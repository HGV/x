package httpx

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Option func(*config)

type config struct {
	logger          *slog.Logger
	shutdownTimeout time.Duration
	signals         []os.Signal
}

func defaultConfig() config {
	return config{
		logger:          slog.Default(),
		shutdownTimeout: 10 * time.Second,
		signals:         []os.Signal{syscall.SIGINT, syscall.SIGTERM},
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(cfg *config) {
		if logger != nil {
			cfg.logger = logger
		}
	}
}

func WithTimeout(d time.Duration) Option {
	return func(cfg *config) {
		if d > 0 {
			cfg.shutdownTimeout = d
		}
	}
}

func WithSignals(signals ...os.Signal) Option {
	return func(cfg *config) {
		if len(signals) > 0 {
			cfg.signals = signals
		}
	}
}

func ListenAndServe(addr string, handler http.Handler, opts ...Option) error {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, cfg.signals...)
	defer signal.Stop(sigCh)

	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case sig := <-sigCh:
		cfg.logger.Info("shutdown signal received", "signal", sig)
	case err := <-errCh:
		cfg.logger.Error("server failed to start", "err", err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.shutdownTimeout)
	defer cancel()

	cfg.logger.Info("shutting down server gracefully")
	if err := srv.Shutdown(ctx); err != nil {
		cfg.logger.Error("server shutdown failed", "err", err)
		return err
	}

	cfg.logger.Info("server shutdown completed gracefully")
	return nil
}
