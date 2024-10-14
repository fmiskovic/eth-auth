package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fmiskovic/eth-auth/logging"
)

type Server interface {
	Start(host string, port uint, api http.Handler) error
	Stop(ctx context.Context) error
}

type server struct {
	*http.Server
}

func New() Server {
	return &server{}
}

func (s *server) Start(host string, port uint, api http.Handler) error {
	s.Server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: api,
	}

	logging.Logger().Info("Starting server", "host", host, "port", port)
	return s.ListenAndServe()
}

func (s *server) Stop(ctx context.Context) error {
	logging.Logger().InfoContext(ctx, "Shutting down server")
	return s.Shutdown(ctx)
}
