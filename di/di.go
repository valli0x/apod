package di

import (
	"flag"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/samber/do"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/hashicorp/go-hclog"
	physFile "github.com/hashicorp/vault/sdk/physical/file"
	"github.com/valli0x/apod/config"
	"github.com/valli0x/apod/model"
	"github.com/valli0x/apod/server"
	"github.com/valli0x/apod/worker"
	"github.com/valli0x/apod/worker/fairshare"
)

func NewLogger(i *do.Injector) (*zerolog.Logger, error) {
	cfg := do.MustInvoke[*config.Config](i)
	log := zerolog.New(os.Stdout).With().Timestamp().Str("ServiceName", cfg.Logger.ServiceName).Caller().Logger()

	return &log, nil
}

func NewStor(i *do.Injector) (*gorm.DB, error) {
	cfg := do.MustInvoke[*config.Config](i)
	logger := do.MustInvoke[*zerolog.Logger](i)

	host := cfg.Storage.Host
	port := cfg.Storage.Port
	user := cfg.Storage.User
	pass := cfg.Storage.Password
	dbname := cfg.Storage.DBname
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable", host, port, user, pass)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	closer, err := db.DB()
	if err != nil {
		return nil, err
	}
	defer closer.Close()

	db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbname))
	logger.Info().Msg(fmt.Sprintf("created database %s", dbname))

	mydsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)
	mydb, err := gorm.Open(postgres.Open(mydsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := mydb.AutoMigrate(&model.APOD{}); err != nil {
		return nil, err
	}
	logger.Info().Msg("created database apod table")

	return mydb, nil
}

func NewConfig(*do.Injector) (*config.Config, error) {
	flag.Parse()
	v, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("load config error: %w", err)
	}

	cfg, err := config.ParseConfig(v)
	if err != nil {
		return nil, fmt.Errorf("parse config error: %w", err)
	}

	return cfg, nil
}

func NewJobManager(i *do.Injector) (*fairshare.JobManager, error) {
	return fairshare.NewJobManager("apod manager", 2, hclog.NewInterceptLogger(&hclog.LoggerOptions{
		Output:     os.Stdout,
		Level:      hclog.DefaultLevel,
		JSONFormat: false,
	}), nil), nil
}

func NewWorker(i *do.Injector) (*worker.APODjob, error) {
	cfg := do.MustInvoke[*config.Config](i)
	metastor := do.MustInvoke[*gorm.DB](i)
	logger := do.MustInvoke[*zerolog.Logger](i)

	dataStor, err := physFile.NewFileBackend(map[string]string{
		"path": cfg.Datapath,
	}, nil)
	if err != nil {
		return nil, err
	}

	return worker.NewAPODjob(metastor, dataStor, logger, cfg.Apodkey), nil
}

func NewServer(i *do.Injector) (*server.Server, error) {
	cfg := do.MustInvoke[*config.Config](i)
	metastor := do.MustInvoke[*gorm.DB](i)

	return server.NewServer(cfg.Server.Address, metastor), nil
}
