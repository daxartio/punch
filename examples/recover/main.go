package main

import (
	"fmt"
	"time"

	"github.com/daxartio/punch"
	"github.com/daxartio/punch/middleware"
)

type (
	Config      = punch.Config[punch.Context]
	ErrorConfig = middleware.ErrorConfig[punch.Context]
)

var (
	IntervalWithConfig = middleware.IntervalWithConfig[punch.Context]
	Recover            = middleware.Recover[punch.Context]
	ErrorWithConfig    = middleware.ErrorWithConfig[punch.Context]
)

func main() {
	p := punch.NewWithConfig(Config{ //nolint
		ContextCreator: punch.NewContext,
		Handler: func(_ punch.Context) error {
			fmt.Println("tick") //nolint
			panic("test")
		},
	})

	p.Use(IntervalWithConfig(middleware.IntervalConfig{
		Interval: func() time.Duration { return time.Second },
	}))
	p.Use(Recover())
	p.Use(ErrorWithConfig(ErrorConfig{
		Handler: func(err error, _ punch.Context) error {
			fmt.Println(err) //nolint

			return nil
		},
	}))

	_ = p.Run()
}
