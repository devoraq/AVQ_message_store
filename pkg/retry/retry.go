// Package retry предоставляет универсальный механизм для повторного выполнения операций
// с использованием настраиваемой экспоненциальной задержки.
package retry

import (
	"context"
	"crypto/rand"
	"errors"
	"log/slog"
	"math/big"
	"time"

	"github.com/devoraq/AVQ_message_store/internal/infrastructure/config"
)

// Do выполняет переданную функцию fn с логикой повторных попыток.
// Функция будет повторяться до успешного выполнения или пока не исчерпаются все попытки.
// Поддерживает отмену через context, экспоненциальный рост задержки и случайный разброс (jitter).
func Do(ctx context.Context, log *slog.Logger, cfg *config.RetryConfig, fn func(ctx context.Context) error) error {
	if log == nil {
		return errors.New("retry.Do: logger cannot be nil")
	}
	if cfg == nil {
		return errors.New("retry.Do: config cannot be nil")
	}

	var err error
	for i := 0; i < cfg.Attempts; i++ {
		if err = fn(ctx); err == nil {
			return nil
		}

		log.Warn("Operation failed, retrying...",
			slog.String("error", err.Error()),
			slog.Int("attempt", i+1),
			slog.Int("total_attempts", cfg.Attempts),
		)

		if i == cfg.Attempts-1 {
			break
		}

		backoff := float64(cfg.Initial) * pow(cfg.Factor, i)

		if cfg.Jitter {
			maxJitter := int64(backoff / 2)
			if maxJitter > 0 {
				jitter, randErr := rand.Int(rand.Reader, big.NewInt(maxJitter))
				if randErr == nil {
					backoff += float64(jitter.Int64())
				}
			}
		}

		if backoff > float64(cfg.Max) {
			backoff = float64(cfg.Max)
		}

		delay := time.Duration(backoff)

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
	return err
}

// pow - простая реализация возведения в степень для float, чтобы не импортировать math.
func pow(base float64, exp int) float64 {
	result := 1.0
	for i := 0; i < exp; i++ {
		result *= base
	}
	return result
}
