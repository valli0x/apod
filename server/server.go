package server

import (
	"net"
	"net/http"
	"time"

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
