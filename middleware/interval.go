package middleware

import (
	"context"
	"time"

	"github.com/daxartio/punch"
)

type IntervalConfig struct {
	Interval func() time.Duration
}

func Interval() punch.MiddlewareFunc {
	return IntervalWithConfig(IntervalConfig{
		Interval: nil,
	})
}

func IntervalWithConfig(config IntervalConfig) punch.MiddlewareFunc {
	return func(next punch.HandlerFunc) punch.HandlerFunc {
		if config.Interval == nil {
			return next
		}

		return func(ctx context.Context) error {
			time.Sleep(config.Interval())

			return next(ctx)
		}
	}
}
