package middleware

import (
	"context"
	"errors"
	"fmt"
	"runtime"

	"github.com/daxartio/punch"
)

const (
	DefaultRecoverStackSize      = 4 << 10 // 4kb
	DefaultRecoverEnableStackAll = true
)

var ErrRecoverPanic = errors.New("panic")

var DefaultRecoverConfig = RecoverConfig{
	StackSize:      DefaultRecoverStackSize,
	EnableStackAll: DefaultRecoverEnableStackAll,
	Handler:        nil,
}

type RecoverHandler func(error, context.Context, []byte)

type RecoverConfig struct {
	StackSize      int
	EnableStackAll bool
	Handler        RecoverHandler
}

func Recover[T punch.Context]() punch.MiddlewareFunc[T] {
	return RecoverWithConfig[T](DefaultRecoverConfig)
}

func RecoverWithConfig[T punch.Context](config RecoverConfig) punch.MiddlewareFunc[T] {
	return func(next punch.HandlerFunc[T]) punch.HandlerFunc[T] {
		return func(ctx T) error {
			var errPanic error

			err := func() error {
				defer handlePanic(ctx.Context(), config, &errPanic)

				return next(ctx)
			}()

			if errPanic != nil {
				return errPanic
			}

			return err
		}
	}
}

func handlePanic(ctx context.Context, config RecoverConfig, errPanic *error) {
	if r := recover(); r != nil {
		err, ok := r.(error)
		if !ok {
			err = fmt.Errorf("%w: %r", ErrRecoverPanic, r)
		}

		var (
			stack  []byte
			length int
		)

		stack = make([]byte, config.StackSize)
		length = runtime.Stack(stack, !config.EnableStackAll)
		stack = stack[:length]

		*errPanic = err

		if config.Handler != nil {
			config.Handler(err, ctx, stack)
		}
	}
}
