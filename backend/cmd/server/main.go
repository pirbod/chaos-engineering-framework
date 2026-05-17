package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pirbod/chaos-engineering-framework/backend/internal/ai"
	"github.com/pirbod/chaos-engineering-framework/backend/internal/api"
	"github.com/pirbod/chaos-engineering-framework/backend/internal/observability"
	"github.com/pirbod/chaos-engineering-framework/backend/internal/store"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	port := firstNonEmpty(os.Getenv("PORT"), "8080")
	version := firstNonEmpty(os.Getenv("VERSION"), "dev")

	metrics := observability.NewRecorder()
	router := api.NewRouter(store.NewInMemoryStore(), ai.NewProviderFromEnv(), metrics, version, logger)
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errs := make(chan error, 1)
	go func() {
		logger.Info("server_starting", "addr", server.Addr, "version", version)
		errs <- server.ListenAndServe()
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errs:
		if err != nil && err != http.ErrServerClosed {
			logger.Error("server_failed", "error", err)
			os.Exit(1)
		}
	case sig := <-signals:
		logger.Info("server_stopping", "signal", sig.String())
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			logger.Error("server_shutdown_failed", "error", err)
			os.Exit(1)
		}
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
