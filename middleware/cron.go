package middleware

import (
	"context"
	"time"

	"github.com/daxartio/punch"
)

var DefaultCronConfig = CronConfig{
	Interval: 0,
}

type CronConfig struct {
	Interval time.Duration
}

func Cron() punch.MiddlewareFunc {
	return CronWithConfig(DefaultCronConfig)
}

func CronWithConfig(config CronConfig) punch.MiddlewareFunc {
	return func(next punch.HandlerFunc) punch.HandlerFunc {
		if config.Interval == 0 {
			return next
		}

		return func(ctx context.Context) error {
			time.Sleep(config.Interval)

			return next(ctx)
		}
	}
}
