package middleware

import (
	"context"
	"time"

	"github.com/daxartio/punch"
)

type TimeoutConfig struct {
	Timeout time.Duration
}

func Timeout() punch.MiddlewareFunc {
	return TimeoutWithConfig(TimeoutConfig{
		Timeout: 0,
	})
}

func TimeoutWithConfig(config TimeoutConfig) punch.MiddlewareFunc {
	return func(next punch.HandlerFunc) punch.HandlerFunc {
		if config.Timeout == 0 {
			return next
		}

		return func(ctx context.Context) error {
			newCtx, cancel := context.WithTimeout(ctx, config.Timeout)
			defer cancel()

			return next(newCtx)
		}
	}
}
