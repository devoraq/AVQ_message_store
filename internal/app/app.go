// Package app содержит координатор жизненного цикла компонентов приложения.
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
	"github.com/devoraq/AVQ_message_store/internal/infrastructure/eventbus/kafka"
	"github.com/devoraq/AVQ_message_store/internal/infrastructure/nosql/mongodb"
)

// App управляет жизненным циклом компонентов приложения.
type App struct {
	log  *slog.Logger
	happ *happ.HApp

	container *Container

	wg           sync.WaitGroup
	shutdownOnce sync.Once
}

// New создает контейнер приложения и подготавливает инфраструктурные компоненты.
func New(_ context.Context, cfg *config.Config, log *slog.Logger) (*App, error) {
	app := &App{
		log:       log,
		container: NewContainer(log, cfg),
	}

	mongo := mustInitMongo(cfg, log)
	kafka := mustInitKafka(cfg, log)

	app.container.Add(mongo, kafka)

	if cfg.IsHTTPEnabled {
		app.happ = buildHTTP(cfg.HTTPConfig, log)
	}

	return app, nil
}

// StartAsync запускает фоновые компоненты, включая HTTP-сервер.
func (a *App) StartAsync(ctx context.Context) {
	const op = "App.StartAsync"
	log := a.log.With("op", op)

	log.Info(
		"Application started",
		slog.Int("PID", os.Getpid()),
	)

	if err := a.container.StartAll(ctx); err != nil {
		log.Error("failed to start components", slog.Any("error", err))
		os.Exit(1)
	}

	// if a.GRPC != nil {
	// 	go a.GRPC.MustStart()
	// }
	if a.happ != nil {
		go a.happ.MustStart()
	}
}

// Shutdown корректно останавливает запущенные компоненты.
func (a *App) Shutdown(ctx context.Context) error {
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
		if a.happ != nil {
			if err := a.happ.Shutdown(ctx); err != nil {
				errs = append(errs, err)
			}
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

func mustInitMongo(cfg *config.Config, log *slog.Logger) *mongodb.MongoDB {
	client, err := mongodb.New(&mongodb.MongoDeps{
		Cfg:    cfg.MongoConfig,
		Logger: log,
	})
	if err != nil {
		log.Error("failed to initialize MongoDB client", slog.String("error", err.Error()))
		os.Exit(1)
	}

	return client
}

func mustInitKafka(cfg *config.Config, log *slog.Logger) *kafka.Kafka {
	return kafka.NewKafka(&kafka.KafkaDeps{Cfg: cfg, Log: log})
}
