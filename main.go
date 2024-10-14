package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/fmiskovic/eth-auth/api"
	"github.com/fmiskovic/eth-auth/logging"
	"github.com/fmiskovic/eth-auth/server"
)

func main() {
	logger := logging.Logger()

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		logger.Warn("environment variable JWT_SECRET is not set, falling back to default")
		secret = "jwt-default-secret"
	}

	router := api.New(secret)
	srv := server.New()

	// Start the server in a goroutine.
	go func() {
		if err := srv.Start("localhost", 8080, router); !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server failed to run", "error", err)
			os.Exit(1)
		}
		logger.Info("stopped serving new connections.")
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()
	<-ctx.Done()

	// Timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Shutdown the server
	if err := srv.Stop(ctx); err != nil {
		logger.Error("server shutdown error", "error", err)
	}
	logger.Info("graceful shutdown completed.")

}
