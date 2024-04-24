package middleware

import (
	"context"

	"github.com/daxartio/punch"
)

type ErrorHandler = func(error, context.Context) error

var DefaultErrorConfig = ErrorConfig{
	Handler: nil,
}

type ErrorConfig struct {
	Handler ErrorHandler
}

func Error() punch.MiddlewareFunc {
	return ErrorWithConfig(DefaultErrorConfig)
}

func ErrorWithConfig(config ErrorConfig) punch.MiddlewareFunc {
	return func(next punch.HandlerFunc) punch.HandlerFunc {
		if config.Handler == nil {
			return next
		}

		return func(ctx context.Context) error {
			if err := next(ctx); err != nil {
				return config.Handler(err, ctx)
			}

			return nil
		}
	}
}
