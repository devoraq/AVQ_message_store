// Package main запускает HTTP-сервис хранения сообщений.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/devoraq/AVQ_message_store/internal/app"
	"github.com/devoraq/AVQ_message_store/internal/infrastructure/config"
	"github.com/devoraq/AVQ_message_store/internal/infrastructure/logger"
)

const pathConfig = "./config/config.yaml"

func main() {
	rootCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger := initLogger()
	cfg := config.LoadConfig(pathConfig)

	appCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a := mustInitApp(appCtx, cfg, logger)
	a.StartAsync(appCtx)

	<-rootCtx.Done()
	logger.Info("termination signal received, shutting down...")

	cancel()

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancelShutdown()

	if err := a.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", slog.String("error", err.Error()))
	} else {
		logger.Info("service stopped gracefully")
	}
}

func initLogger() *slog.Logger {
	logHandler := logger.NewPrettyHandler(os.Stdout, logger.PrettyHandlerOptions{
		Opts: slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)
	return logger
}

func mustInitApp(ctx context.Context, cfg *config.Config, log *slog.Logger) *app.App {
	a, err := app.New(ctx, cfg, log)
	if err != nil {
		log.Error("failed to initialize app", slog.String("error", err.Error()))
		os.Exit(1)
	}
	return a
}
