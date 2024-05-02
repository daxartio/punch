package punch

import (
	"context"
	"sync/atomic"
)

type HandlerFunc[T Context] func(ctx T) error

type MiddlewareFunc[T Context] func(next HandlerFunc[T]) HandlerFunc[T]

type Config[T Context] struct {
	Handler       HandlerFunc[T]
	CreateContext func(ctx context.Context) T
}

type Punch[T Context] struct {
	handler                HandlerFunc[T]
	handlerWithMiddlewares HandlerFunc[T]
	createContext          func(ctx context.Context) T
	started                chan struct{}
	stop                   chan struct{}
	stopped                chan struct{}
	stopFlag               *atomic.Value
}

func New() *Punch[*Ctx] {
	return NewWithConfig(Config[*Ctx]{}) //nolint:exhaustruct
}

func NewWithConfig[T Context](config Config[T]) *Punch[T] {
	return &Punch[T]{
		handler:                config.Handler,
		handlerWithMiddlewares: config.Handler,
		createContext:          config.CreateContext,
		started:                make(chan struct{}),
		stop:                   nil,
		stopped:                nil,
		stopFlag:               new(atomic.Value),
	}
}

func (w *Punch[T]) SetHandler(handler HandlerFunc[T]) {
	w.handler = handler
	w.handlerWithMiddlewares = handler
}

func (w *Punch[T]) Run() error {
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
		pCtx := w.createContext(ctx)

		go func() {
			defer cancel()

			_ = w.handlerWithMiddlewares(pCtx)
		}()

		select {
		case <-ctx.Done():
		case <-w.stop:
			cancel()
			<-ctx.Done()

			return nil
		}
	}
}

func (w *Punch[T]) Shutdown(_ context.Context) error {
	w.stopFlag.Store(true)
	w.stop <- struct{}{}

	return nil
}

func (w *Punch[T]) Started() <-chan struct{} {
	return w.started
}

func (w *Punch[T]) Stopped() <-chan struct{} {
	return w.stopped
}

func (w *Punch[T]) Use(middleware MiddlewareFunc[T]) {
	w.handlerWithMiddlewares = middleware(w.handlerWithMiddlewares)
}
