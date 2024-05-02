package punch_test

import (
	"context"
	"testing"
	"time"

	"github.com/daxartio/punch"
)

func TestPunch(t *testing.T) {
	t.Parallel()

	handled := make(chan bool)
	defer close(handled)

	p := punch.NewWithConfig(punch.Config[punch.Context]{ //nolint:varnamelen
		ContextCreator: punch.NewContext,
		Handler: func(_ punch.Context) error {
			time.Sleep(time.Millisecond * 100)

			handled <- true

			return nil
		},
	})
	p.Use(func(next punch.HandlerFunc[punch.Context]) punch.HandlerFunc[punch.Context] {
		return func(ctx punch.Context) error {
			t.Log("middleware 1")

			return next(ctx)
		}
	})
	p.Use(func(next punch.HandlerFunc[punch.Context]) punch.HandlerFunc[punch.Context] {
		return func(ctx punch.Context) error {
			t.Log("middleware 2")

			return next(ctx)
		}
	})

	t.Log("punch starting")

	go p.Run() //nolint:errcheck

	<-p.Started()
	t.Log("punch started")

	<-handled
	t.Log("handled")

	t.Log("punch stoping")

	err := p.Shutdown(context.TODO())
	if err != nil {
		t.Fatal(err.Error())
	}

	<-p.Stopped()
	t.Log("punch stopped")
}
