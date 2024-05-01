package middleware

import (
	"context"
	"errors"
	"time"

	"github.com/daxartio/punch"
)

var ErrIntervalCancelled = errors.New("interval cancelled")

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
			timer := time.NewTimer(config.Interval())

			select {
			case <-timer.C:
			case <-ctx.Done():
				timer.Stop()

				return ErrIntervalCancelled
			}

			return next(ctx)
		}
	}
}
