package app

import (
	"api-gateway/internal/config"
	"time"

	"github.com/er-davo/retry"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

func zapLogger(log *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)

			stop := time.Now()
			latency := stop.Sub(start)

			req := c.Request()
			res := c.Response()

			log.Info("request handled",
				zap.String("method", req.Method),
				zap.String("path", req.URL.Path),
				zap.Int("status", res.Status),
				zap.Duration("latency", latency),
				zap.String("remote_ip", c.RealIP()),
				zap.String("user_agent", req.UserAgent()),
			)

			return err
		}
	}
}

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
