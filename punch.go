package punch

import (
	"context"
	"sync/atomic"
)

type HandlerFunc func(ctx context.Context) error

type MiddlewareFunc func(next HandlerFunc) HandlerFunc

type Lifecycle interface {
	Before(ctx context.Context) error
	After(ctx context.Context) error
}

var _ Lifecycle = (*dummyLifecycle)(nil)

type dummyLifecycle struct{}

func (d *dummyLifecycle) After(_ context.Context) error {
	return nil
}

func (d *dummyLifecycle) Before(_ context.Context) error {
	return nil
}

type Config struct {
	Handler   HandlerFunc
	Lifecycle Lifecycle
}

type Punch struct {
	handler                HandlerFunc
	handlerWithMiddlewares HandlerFunc
	lifecycle              Lifecycle
	started                chan struct{}
	stop                   chan struct{}
	stopped                chan struct{}
	stopFlag               *atomic.Value
}

func New() *Punch {
	return NewWithConfig(Config{}) //nolint:exhaustruct
}

func NewWithConfig(config Config) *Punch {
	if config.Handler == nil {
		config.Handler = func(_ context.Context) error { return nil }
	}

	if config.Lifecycle == nil {
		config.Lifecycle = &dummyLifecycle{}
	}

	return &Punch{
		handler:                config.Handler,
		handlerWithMiddlewares: config.Handler,
		lifecycle:              config.Lifecycle,
		started:                make(chan struct{}),
		stop:                   nil,
		stopped:                nil,
		stopFlag:               new(atomic.Value),
	}
}

func (w *Punch) Run() error {
	w.stop = make(chan struct{})
	w.stopped = make(chan struct{})

	defer close(w.stop)
	defer close(w.stopped)

	close(w.started)

	for {
		if w.stopFlag.Load() != nil {
			return nil
		}

		ctx, cancel := context.WithCancel(context.Background())

		_ = w.lifecycle.Before(ctx)

		go func() {
			defer cancel()

			_ = w.handlerWithMiddlewares(ctx)
		}()

		select {
		case <-ctx.Done():
		case <-w.stop:
			cancel()
			<-ctx.Done()

			return nil
		}

		_ = w.lifecycle.After(ctx)
	}
}

func (w *Punch) Shutdown(_ context.Context) error {
	w.stopFlag.Store(true)
	w.stop <- struct{}{}
	<-w.stopped

	return nil
}

func (w *Punch) Started() <-chan struct{} {
	return w.started
}

func (w *Punch) Stopped() <-chan struct{} {
	return w.stopped
}

func (w *Punch) Use(middleware MiddlewareFunc) {
	w.handlerWithMiddlewares = middleware(w.handlerWithMiddlewares)
}
