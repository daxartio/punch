package main

import (
	"context"
	"fmt"
	"time"

	"github.com/daxartio/punch"
	"github.com/daxartio/punch/middleware"
)

type Config = punch.Config[*punch.Ctx]

var (
	IntervalWithConfig = middleware.IntervalWithConfig[*punch.Ctx]
	Recover            = middleware.Recover[*punch.Ctx]
	ErrorWithConfig    = middleware.ErrorWithConfig[*punch.Ctx]
)

func main() {
	p := punch.NewWithConfig(Config{ //nolint
		CreateContext: punch.NewCtx,
		Handler: func(_ *punch.Ctx) error {
			fmt.Println("tick") //nolint
			panic("test")
		},
	})

	p.Use(IntervalWithConfig(middleware.IntervalConfig{
		Interval: func() time.Duration { return time.Second },
	}))
	p.Use(Recover())
	p.Use(ErrorWithConfig(middleware.ErrorConfig{
		Handler: func(err error, _ context.Context) error {
			fmt.Println(err) //nolint

			return nil
		},
	}))

	_ = p.Run()
}
