package main

import (
	"context"
	"fmt"
	"time"

	"github.com/daxartio/punch"
	"github.com/daxartio/punch/middleware"
)

func main() {
	p := punch.NewWithConfig(punch.Config{ //nolint
		Handler: func(_ context.Context) error {
			fmt.Println("tick") //nolint
			panic("test")
		},
	})
	p.Use(middleware.IntervalWithConfig(middleware.IntervalConfig{
		Interval: func() time.Duration { return time.Second },
	}))
	p.Use(middleware.Recover())
	p.Use(middleware.ErrorWithConfig(middleware.ErrorConfig{
		Handler: func(err error, _ context.Context) error {
			fmt.Println(err) //nolint

			return nil
		},
	}))

	_ = p.Run()
}
