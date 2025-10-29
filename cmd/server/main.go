package main

import (
	"context"
	"fmt"
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
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger := initLogger()

	cfg := config.LoadConfig(pathConfig)
	fmt.Println(cfg)

	app, err := app.New(ctx, cfg, logger)
	if err != nil {
		logger.Error("app error: ", slog.String("err", err.Error()))
		os.Exit(1)
	}
	app.StartAsync()

	<-ctx.Done()
	logger.Info("signal received, shutting down...")

	app.Shutdown(ctx)
}

func initLogger() *slog.Logger {
	logHandler := logger.NewPrettyHandler(os.Stdout, logger.PrettyHandlerOptions{})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)
	return logger
}
