package middleware

import (
	"context"
	"time"

	"github.com/daxartio/punch"
)

var DefaultTimeoutConfig = TimeoutConfig{
	Timeout: 0,
}

type TimeoutConfig struct {
	Timeout time.Duration
}

func Timeout() punch.MiddlewareFunc {
	return TimeoutWithConfig(DefaultTimeoutConfig)
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
