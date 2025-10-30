package retry

import (
	"context"
	"time"

	"github.com/devoraq/AVQ_message_store/internal/infrastructure/config"
)

type Backoff struct {
	cfg   *config.Config
	delay time.Duration
}

func NewBackoff(cfg *config.Config) *Backoff {
	sanitize(cfg)
	return &Backoff{
		cfg:   cfg,
		delay: cfg.Initial,
	}
}

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

func (b *Backoff) Reset() {
	b.delay = b.cfg.Initial
}
