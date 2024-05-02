# Punch

[![Go Documentation](https://godocs.io/github.com/daxartio/punch?status.svg)](https://godocs.io/github.com/daxartio/punch)

## Usage

```go
package main

import (
	"fmt"
	"time"

	"github.com/daxartio/punch"
	"github.com/daxartio/punch/middleware"
)

func main() {
	p := punch.New() //nolint

	p.SetHandler(func(_ punch.Context) error {
		fmt.Println("tick") //nolint

		return nil
	})

	p.Use(middleware.IntervalWithConfig[punch.Context](middleware.IntervalConfig{
		Interval: func() time.Duration { return time.Second },
	}))

	_ = p.Run()
}

```
