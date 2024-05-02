package main

import (
	"context"
	"fmt"
	"time"

	"github.com/daxartio/punch"
	"github.com/daxartio/punch/middleware"
)

type Context interface {
	punch.Context
	MyMethod()
}

var _ Context = (*MyContext)(nil)

type MyContext struct {
	ctx context.Context //nolint
}

func (m *MyContext) Context() context.Context {
	return m.ctx
}

func (m *MyContext) SetContext(ctx context.Context) {
	m.ctx = ctx
}

func (m *MyContext) MyMethod() {
	fmt.Println("my context") //nolint
}

func MyMiddleware() punch.MiddlewareFunc[Context] {
	return func(next punch.HandlerFunc[Context]) punch.HandlerFunc[Context] {
		return func(ctx Context) error {
			ctx.MyMethod()

			return next(ctx)
		}
	}
}

func main() {
	p := punch.NewWithConfig( //nolint
		punch.Config[Context]{
			Handler: func(_ Context) error {
				fmt.Println("tick") //nolint

				return nil
			},
			ContextCreator: func(ctx context.Context) Context {
				return &MyContext{ctx: ctx}
			},
		},
	)

	p.Use(MyMiddleware())
	p.Use(middleware.IntervalWithConfig[Context](middleware.IntervalConfig{
		Interval: func() time.Duration { return time.Second },
	}))

	_ = p.Run()
}
