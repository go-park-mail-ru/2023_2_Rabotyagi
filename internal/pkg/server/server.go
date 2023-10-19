package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/pkg/config"
)

const (
	basicTimeout = 10 * time.Second
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(config *config.Config, handler http.Handler) error {
	s.httpServer = &http.Server{ //nolint:exhaustruct
		Addr:           ":" + config.PortServer,
		Handler:        handler,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
		ReadTimeout:    basicTimeout,
		WriteTimeout:   basicTimeout,
	}

	log.Printf("Start server:%s", config.PortServer)

	return s.httpServer.ListenAndServe() //nolint:wrapcheck
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx) //nolint:wrapcheck
}
