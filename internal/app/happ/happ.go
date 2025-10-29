package happ

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/devoraq/AVQ_message_store/internal/infrastructure/config"
)

type HApp struct {
	log    *slog.Logger
	server *http.Server
	cfg    *config.HTTPConfig
}

func NewHApp(cfg *config.HTTPConfig, log *slog.Logger, mux *http.ServeMux) *HApp {
	server := &http.Server{
		Addr:              cfg.Addr,
		Handler:           mux,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}

	return &HApp{
		log:    log,
		server: server,
		cfg:    cfg,
	}
}

func (ha *HApp) Start() error {
	const op = "HApp.Start"
	log := ha.log.With("op", op)

	log.Info(
		"HTTP server is starting",
		slog.String("address", ha.cfg.Addr),
	)

	if err := ha.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (ha *HApp) MustStart() {
	if err := ha.Start(); err != nil {
		ha.log.Error("HTTP server failed", "err", err)
		panic(err)
	}
}

func (ha *HApp) Shutdown(ctx context.Context) error {
	if err := ha.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
