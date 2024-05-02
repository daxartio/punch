package middleware

import (
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

type RecoverHandler[T punch.Context] func(error, T, []byte)

type RecoverConfig[T punch.Context] struct {
	StackSize      int
	EnableStackAll bool
	Handler        RecoverHandler[T]
}

func Recover[T punch.Context]() punch.MiddlewareFunc[T] {
	return RecoverWithConfig(RecoverConfig[T]{
		StackSize:      DefaultRecoverStackSize,
		EnableStackAll: DefaultRecoverEnableStackAll,
		Handler:        nil,
	})
}

func RecoverWithConfig[T punch.Context](config RecoverConfig[T]) punch.MiddlewareFunc[T] {
	return func(next punch.HandlerFunc[T]) punch.HandlerFunc[T] {
		return func(ctx T) error {
			var errPanic error

			err := func() error {
				defer handlePanic(ctx, config, &errPanic)

				return next(ctx)
			}()

			if errPanic != nil {
				return errPanic
			}

			return err
		}
	}
}

func handlePanic[T punch.Context](ctx T, config RecoverConfig[T], errPanic *error) {
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
