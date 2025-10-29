package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/devoraq/AVQ_message_store/internal/infrastructure/config"
	"github.com/devoraq/AVQ_message_store/pkg/retry"
)

// Component описывает компонент, жизненным циклом которого управляет контейнер.
type Component interface {
	Name() string
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

// Container хранит набор компонентов и управляет их жизненным циклом.
type Container struct {
	comps       []Component
	log         *slog.Logger
	retryConfig *config.RetryConfig
}

// NewContainer создаёт пустой контейнер без зарегистрированных компонентов.
func NewContainer(log *slog.Logger, retryConfig *config.RetryConfig) *Container {
	return &Container{
		log:         log,
		retryConfig: retryConfig,
	}
}

// Add добавляет один или несколько компонентов в контейнер.
func (c *Container) Add(comp ...Component) { c.comps = append(c.comps, comp...) }

// StartAll последовательно запускает компоненты, накапливая ошибки запуска.
func (c *Container) StartAll(ctx context.Context) error {
	var errs []error
	for _, comp := range c.comps {
		component := comp

		err := retry.Do(ctx, c.log, c.retryConfig, func(ctx context.Context) error {
			return component.Start(ctx)
		})

		if err != nil {
			errs = append(errs, fmt.Errorf("%s start failed: %w", component.Name(), err))
		}
	}
	return errors.Join(errs...)
}

// StopAll останавливает компоненты в обратном порядке, обеспечивая корректное завершение зависимостей.
func (c *Container) StopAll(ctx context.Context) error {
	var errs []error
	for i := len(c.comps) - 1; i >= 0; i-- {
		if err := c.comps[i].Stop(ctx); err != nil {
			errs = append(errs, fmt.Errorf("%s stop failed: %w", c.comps[i].Name(), err))
		}
	}
	return errors.Join(errs...)
}
