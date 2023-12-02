package worker

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/rs/zerolog"
	"github.com/valli0x/apod/model"
	"github.com/valli0x/apod/worker/dataparser"
	"gorm.io/gorm"
)

type APODjob struct {
	client   *http.Client
	apikey   string
	metaStor *gorm.DB
	dataStor physical.Backend
	logger   *zerolog.Logger
}

func NewAPODjob(metaStor *gorm.DB, dataStor physical.Backend, logger *zerolog.Logger, apikey string) *APODjob {
	return &APODjob{
		metaStor: metaStor,
		dataStor: dataStor,
		apikey:   apikey,
		logger:   logger,
		client: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

func (a *APODjob) Execute() error {
	a.logger.Info().Msg("get metadata and save it")

	metadata, err := apodmatadata(a.client, a.apikey)
	if err != nil {
		return err
	}
	apod := &model.APOD{}
	if err := dataparser.Decode(metadata, apod); err != nil {
		return err
	}
	apod.IDData, err = uuid.GenerateUUID()
	if err != nil {
		return err
	}
	a.metaStor.Create(apod)

	a.logger.Info().Msg("get jpg and save it")

	jpgdata, err := apodjpg(a.client, apod.URL)
	if err != nil {
		return err
	}
	if err := a.dataStor.Put(context.Background(), &physical.Entry{
		Key:   apod.IDData,
		Value: compress(jpgdata),
	}); err != nil {
		return err
	}

	return nil
}

func (a *APODjob) OnFailure(err error) {
	a.logger.Err(err).Msg("fail") // :))
}

func apodmatadata(client *http.Client, apikey string) (map[string]interface{}, error) {
	url := "https://api.nasa.gov/planetary/apod?api_key="
	req, err := http.NewRequest(http.MethodGet, url+apikey, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data := make(map[string]interface{}, 10)
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data, err
}

func apodjpg(client *http.Client, url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func compress(n []byte) []byte {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	w.Write(n)
	w.Close()
	return buf.Bytes()
}
