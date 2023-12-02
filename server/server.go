package server

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/valli0x/apod/model"
	"gorm.io/gorm"
)

type Server struct {
	port     string
	metaStor *gorm.DB
}

func NewServer(port string, metaStor *gorm.DB) *Server {
	return &Server{
		port:     port,
		metaStor: metaStor,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	mux.Handle("/v1/apod", s.apod()) // пока только получение метаданных

	ln, err := net.Listen("tcp", s.port)
	if err != nil {
		return err
	}
	server := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		IdleTimeout:       5 * time.Minute,
	}
	server.Serve(ln)

	return nil
}

func (s *Server) apod() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}
		queryVals := r.URL.Query()

		var list bool
		var err error

		listStr := queryVals.Get("list")
		if listStr != "" {
			list, err = strconv.ParseBool(listStr)
			if err != nil {
				respondError(w, http.StatusBadRequest, nil)
				return
			}
		}

		dateStr := queryVals.Get("date")

		switch list {
		case false:
			apod := &model.APOD{}
			result := s.metaStor.First(apod, "date = ?", dateStr)
			if result.Error != nil || apod == nil {
				respondError(w, http.StatusInternalServerError, nil)
				return
			}
			respondOk(w, apod)
			return
		case true:
			var apods []model.APOD
			result := s.metaStor.First(&apods)
			if result.Error != nil || apods == nil {
				respondError(w, http.StatusInternalServerError, nil)
				return
			}
			respondOk(w, apods)
			return
		}

	})
}

func respondError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	type ErrorResponse struct {
		Errors []string `json:"errors"`
	}
	resp := &ErrorResponse{Errors: make([]string, 0, 1)}
	if err != nil {
		resp.Errors = append(resp.Errors, err.Error())
	}

	enc := json.NewEncoder(w)
	enc.Encode(resp)
}

func respondOk(w http.ResponseWriter, body interface{}) {
	w.Header().Set("Content-Type", "application/json")

	if body == nil {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusOK)
		enc := json.NewEncoder(w)
		enc.Encode(body)
	}
}
