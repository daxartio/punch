# Punch

[![Go Documentation](https://godocs.io/github.com/daxartio/punch?status.svg)](https://godocs.io/github.com/daxartio/punch)

## Usage

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/daxartio/punch"
	"github.com/daxartio/punch/middleware"
)

func main() {
	p := punch.NewWithConfig(punch.Config{
		Handler: func(_ context.Context) error {
			fmt.Println("tick")

			return nil
		},
	})
	p.Use(middleware.CronWithConfig(middleware.CronConfig{
		Interval: time.Second,
	}))

	_ = p.Run()
}
```
