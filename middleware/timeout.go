package middleware

import (
	"context"
	"time"

	"github.com/daxartio/punch"
)

type TimeoutConfig struct {
	Timeout time.Duration
}

func Timeout[T punch.Context]() punch.MiddlewareFunc[T] {
	return TimeoutWithConfig[T](TimeoutConfig{
		Timeout: 0,
	})
}

func TimeoutWithConfig[T punch.Context](config TimeoutConfig) punch.MiddlewareFunc[T] {
	return func(next punch.HandlerFunc[T]) punch.HandlerFunc[T] {
		if config.Timeout == 0 {
			return next
		}

		return func(ctx T) error {
			newCtx, cancel := context.WithTimeout(ctx.Context(), config.Timeout)
			defer cancel()

			ctx.SetContext(newCtx)

			return next(ctx)
		}
	}
}
