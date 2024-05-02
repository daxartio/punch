package middleware

import (
	"context"

	"github.com/daxartio/punch"
)

type ErrorHandler = func(error, context.Context) error

type ErrorConfig struct {
	Handler ErrorHandler
}

func Error[T punch.Context]() punch.MiddlewareFunc[T] {
	return ErrorWithConfig[T](ErrorConfig{
		Handler: nil,
	})
}

func ErrorWithConfig[T punch.Context](config ErrorConfig) punch.MiddlewareFunc[T] {
	return func(next punch.HandlerFunc[T]) punch.HandlerFunc[T] {
		if config.Handler == nil {
			return next
		}

		return func(ctx T) error {
			if err := next(ctx); err != nil {
				return config.Handler(err, ctx.Context())
			}

			return nil
		}
	}
}
