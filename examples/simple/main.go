package main

import (
	"fmt"
	"time"

	"github.com/daxartio/punch"
	"github.com/daxartio/punch/middleware"
)

func main() {
	p := punch.New() //nolint

	p.SetHandler(func(_ *punch.Ctx) error {
		fmt.Println("tick") //nolint

		return nil
	})

	p.Use(middleware.IntervalWithConfig[*punch.Ctx](middleware.IntervalConfig{
		Interval: func() time.Duration { return time.Second },
	}))

	_ = p.Run()
}
