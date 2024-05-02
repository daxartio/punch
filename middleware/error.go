package middleware

import (
	"github.com/daxartio/punch"
)

type ErrorHandler[T punch.Context] func(error, T) error

type ErrorConfig[T punch.Context] struct {
	Handler ErrorHandler[T]
}

func Error[T punch.Context]() punch.MiddlewareFunc[T] {
	return ErrorWithConfig[T](ErrorConfig[T]{
		Handler: nil,
	})
}

func ErrorWithConfig[T punch.Context](config ErrorConfig[T]) punch.MiddlewareFunc[T] {
	return func(next punch.HandlerFunc[T]) punch.HandlerFunc[T] {
		if config.Handler == nil {
			return next
		}

		return func(ctx T) error {
			if err := next(ctx); err != nil {
				return config.Handler(err, ctx)
			}

			return nil
		}
	}
}
