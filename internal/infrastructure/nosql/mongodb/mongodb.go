// Package mongodb предоставляет вспомогательные типы для работы с клиентом MongoDB.
package mongodb

import (
	"context"
	"errors"
	"log/slog"

	"github.com/devoraq/AVQ_message_store/internal/infrastructure/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// MongoDB инкапсулирует клиент MongoDB и метаданные сервиса.
type MongoDB struct {
	deps *MongoDeps

	name string
	*mongo.Client
}

// MongoDeps содержит зависимости, необходимые для инициализации MongoDB-клиента.
type MongoDeps struct {
	Cfg    *config.MongoConfig
	Logger *slog.Logger
}

// New создает клиент MongoDB по заданной конфигурации.
func New(deps *MongoDeps) (*MongoDB, error) {
	opts := options.ClientOptions{
		Hosts: []string{deps.Cfg.Addr},
		Auth: &options.Credential{
			Username: deps.Cfg.Username,
			Password: deps.Cfg.Password,
		},
		ConnectTimeout: &deps.Cfg.ConnectTimeout,
		MaxPoolSize:    &deps.Cfg.MaxPoolSize,
	}

	client, err := mongo.Connect(&opts)
	if err != nil {
		return nil, errors.Join(ErrConnect, err)
	}

	mongo := MongoDB{
		deps:   deps,
		name:   "mongodb",
		Client: client,
	}

	return &mongo, nil
}

// Name возвращает имя компонента.
func (md *MongoDB) Name() string { return md.name }

// Start проверяет доступность кластера MongoDB.
func (md *MongoDB) Start(ctx context.Context) error {
	log := md.deps.Logger.With(
		slog.String("component", "mongodb"),
		slog.String("addr", md.deps.Cfg.Addr),
	)

	if err := md.Ping(ctx, readpref.Primary()); err != nil {
		return errors.Join(ErrPing, err)
	}

	log.Debug("MongoDB connection established",
		slog.String("database", md.deps.Cfg.DB),
		slog.Uint64("pool_size", md.deps.Cfg.MaxPoolSize),
	)

	return nil
}

// Stop корректно закрывает соединение с MongoDB.
func (md *MongoDB) Stop(ctx context.Context) error {
	log := md.deps.Logger.With(
		slog.String("component", "mongodb"),
		slog.String("addr", md.deps.Cfg.Addr),
	)

	log.Debug("Disconnecting from MongoDB")

	if err := md.Disconnect(ctx); err != nil {
		return errors.Join(ErrDisconnect, err)
	}

	return nil
}

// GetDB возвращает дескриптор указанной базы данных.
func (md *MongoDB) GetDB(name string) *mongo.Database { return md.Database(name) }
