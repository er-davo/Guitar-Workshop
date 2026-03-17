package app

import (
	"audio-sep-task-service/internal/config"

	"github.com/er-davo/retry"
)

func newRetrier(cfg config.Retry) retry.Retrier {
	opts := []retry.RetryOption{
		retry.WithMaxAttempts(cfg.MaxAttempts),
	}

	if cfg.Backoff == "exponential" {
		opts = append(opts, retry.WithBackoff(retry.ExponentialBackoff{
			Base:   cfg.Base,
			Factor: cfg.Factor,
			Max:    cfg.Max,
			Jitter: cfg.Jitter,
		}))
	}

	return retry.New(opts...)
}
