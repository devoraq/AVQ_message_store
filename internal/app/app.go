package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"

	"github.com/devoraq/AVQ_message_store/internal/app/happ"
	"github.com/devoraq/AVQ_message_store/internal/infrastructure/config"
)

type App struct {
	log  *slog.Logger
	happ *happ.HApp

	container *Container

	wg           sync.WaitGroup
	shutdownOnce sync.Once
}

func New(ctx context.Context, cfg *config.Config, log *slog.Logger) (*App, error) {
	app := &App{
		log:       log,
		container: NewContainer(log, cfg.RetryConfig),
	}

	if cfg.AppConfig.IsHttpEnabled {
		app.happ = buildHTTP(cfg.HTTPConfig, log)
	}

	return app, nil
}

func (a *App) StartAsync() {
	const op = "App.StartAsync"
	log := a.log.With("op", op)

	log.Info(
		"Application started",
		slog.Int("PID", os.Getpid()),
	)

	// if a.GRPC != nil {
	// 	go a.GRPC.MustStart()
	// }
	if a.happ != nil {
		go a.happ.MustStart()
	}
}

// Shutdown корректно останавливает HTTP-сервер и все зарегистрированные компоненты.
func (a *App) Shutdown(ctx context.Context) error {
	if ctx == nil {
		return errors.New("app shutdown: context is nil")
	}

	select {
	case <-ctx.Done():
		return fmt.Errorf("app shutdown: %w", ctx.Err())
	default:
	}

	var errs []error
	a.shutdownOnce.Do(func() {
		// if a.consumerCancel != nil {
		// 	a.consumerCancel()
		// }
		if err := a.happ.Shutdown(ctx); err != nil {
			errs = append(errs, err)
		}
		if err := a.container.StopAll(ctx); err != nil {
			errs = append(errs, err)
		}
		a.wg.Wait()
	})

	return errors.Join(errs...)
}

func buildHTTP(cfg *config.HTTPConfig, log *slog.Logger) *happ.HApp {
	mux := http.NewServeMux()
	return happ.NewHApp(cfg, log, mux)
}
