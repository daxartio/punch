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

			return nil
		},
	})
	p.Use(middleware.CronWithConfig(middleware.CronConfig{
		Interval: time.Second,
	}))

	_ = p.Run()
}
