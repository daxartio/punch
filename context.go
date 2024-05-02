package punch

import "context"

type Context interface {
	Context() context.Context
	SetContext(ctx context.Context)
}

var _ Context = (*Ctx)(nil)

func NewCtx(ctx context.Context) *Ctx {
	return &Ctx{ctx: ctx}
}

type Ctx struct {
	ctx context.Context //nolint
}

func (c *Ctx) SetContext(ctx context.Context) {
	c.ctx = ctx
}

func (c *Ctx) Context() context.Context {
	return c.ctx
}
