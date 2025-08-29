package app

import (
	"time"

	"github.com/labstack/echo"
	"go.uber.org/zap"
)

func ZapLogger(log *zap.Logger) echo.MiddlewareFunc {
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
