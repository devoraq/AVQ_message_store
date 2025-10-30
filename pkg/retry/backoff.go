// Package retry предоставляет инструменты для реализации экспоненциального
// backoff при повторных попытках выполнения операций.
//
// Пакет используется для управления задержками между повторными запросами
// (например, при работе с Kafka, сетевыми вызовами, БД и т.д.),
// поддерживая экспоненциальный рост интервала и опциональный джиттер.ё
package retry

import (
	"context"
	"time"

	"github.com/devoraq/AVQ_message_store/internal/infrastructure/config"
)

// Backoff представляет собой структуру, управляющую задержками между повторными
// попытками выполнения операции. Использует параметры из RetryConfig.
type Backoff struct {
	cfg   *config.Config
	delay time.Duration
}

// NewBackoff создаёт новый экземпляр Backoff с начальными параметрами.
// Если в конфигурации заданы нулевые значения, они будут заменены на безопасные по умолчанию.
//
// Пример:
//
//	b := retry.NewBackoff(&config.RetryConfig{Attempts: 5, Initial: 500 * time.Millisecond})
func NewBackoff(cfg *config.Config) *Backoff {
	sanitize(cfg)
	return &Backoff{
		cfg:   cfg,
		delay: cfg.Initial,
	}
}

// Sleep приостанавливает выполнение на рассчитанный интервал задержки.
// Метод увеличивает текущую задержку экспоненциально и учитывает джиттер (если включён).
// Поддерживает отмену по контексту.
func (b *Backoff) Sleep(ctx context.Context) {
	b.delay = nextDelay(b.delay, b.cfg)

	timer := time.NewTimer(b.delay)
	select {
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
	case <-timer.C:
	}
}

// Reset сбрасывает накопленный интервал задержки к начальному значению,
// определённому в конфигурации.
func (b *Backoff) Reset() {
	b.delay = b.cfg.Initial
}
