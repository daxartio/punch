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

func Recover() punch.MiddlewareFunc {
	return RecoverWithConfig(DefaultRecoverConfig)
}

func RecoverWithConfig(config RecoverConfig) punch.MiddlewareFunc {
	return func(next punch.HandlerFunc) punch.HandlerFunc {
		return func(ctx context.Context) error {
			defer handlePanic(ctx, config)

			return next(ctx)
		}
	}
}

func handlePanic(ctx context.Context, config RecoverConfig) {
	if r := recover(); r != nil {
		err, ok := r.(error)
		if !ok {
			err = fmt.Errorf("%w: %v", ErrRecoverPanic, r)
		}

		var (
			stack  []byte
			length int
		)

		stack = make([]byte, config.StackSize)
		length = runtime.Stack(stack, !config.EnableStackAll)
		stack = stack[:length]

		if config.Handler != nil {
			config.Handler(err, ctx, stack)
		}
	}
}