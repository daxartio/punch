package middleware

import (
	"errors"
	"time"

	"github.com/daxartio/punch"
)

var ErrIntervalCancelled = errors.New("interval cancelled")

type IntervalConfig struct {
	Interval func() time.Duration
}

func Interval[T punch.Context]() punch.MiddlewareFunc[T] {
	return IntervalWithConfig[T](IntervalConfig{
		Interval: nil,
	})
}

func IntervalWithConfig[T punch.Context](config IntervalConfig) punch.MiddlewareFunc[T] {
	return func(next punch.HandlerFunc[T]) punch.HandlerFunc[T] {
		if config.Interval == nil {
			return next
		}

		return func(ctx T) error {
			timer := time.NewTimer(config.Interval())

			select {
			case <-timer.C:
			case <-ctx.Context().Done():
				timer.Stop()

				return ErrIntervalCancelled
			}

			return next(ctx)
		}
	}
}
