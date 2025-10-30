// Package retry предоставляет универсальный механизм для повторного выполнения операций
// с использованием настраиваемой экспоненциальной задержки.
package retry

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"time"

	"github.com/devoraq/AVQ_message_store/internal/infrastructure/config"
)

// Do выполняет переданную функцию fn с логикой повторных попыток.
// Функция будет повторяться до успешного выполнения или пока не исчерпаются все попытки.
// Поддерживает отмену через context, экспоненциальный рост задержки и случайный разброс (jitter).
func Do(ctx context.Context, cfg *config.Config, fn func(ctx context.Context) error) (lastErr error) {
	if cfg == nil {
		return errors.New("retry.Do: config cannot be nil")
	}

	// sane defaults, если значения не заданы
	sanitize(cfg)

	delay := cfg.Initial

	for attempt := 1; attempt <= cfg.Attempts; attempt++ {
		if err := fn(ctx); err != nil {
			lastErr = err
		}

		// если это была последняя попытка — выходим
		if attempt == cfg.Attempts {
			break
		}

		// считаем следующую задержку
		delay = nextDelay(delay, cfg)

		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			return ctx.Err()
		case <-timer.C:
		}
	}

	return lastErr
}

// nextDelay возвращает увеличенный интервал с crypto-jitter.
func nextDelay(prev time.Duration, cfg *config.Config) time.Duration {
	backoff := float64(prev) * cfg.Factor
	if backoff > float64(cfg.Max) {
		backoff = float64(cfg.Max)
	}

	if cfg.Jitter {
		jitterMax := int64(backoff / 2)
		if jitterMax > 0 {
			if n, err := rand.Int(rand.Reader, big.NewInt(jitterMax)); err == nil {
				backoff += float64(n.Int64())
			}
		}
	}

	return time.Duration(backoff)
}

func sanitize(cfg *config.Config) {
	if cfg.Attempts <= 0 {
		cfg.Attempts = 3
	}
	if cfg.Initial <= 0 {
		cfg.Initial = time.Second
	}
	if cfg.Max <= 0 {
		cfg.Max = 30 * time.Second
	}
	if cfg.Factor < 1 {
		cfg.Factor = 2.0
	}
}
