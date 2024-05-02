package punch

import "context"

type Context interface {
	Context() context.Context
	SetContext(ctx context.Context)
}

var _ Context = (*baseContext)(nil)

func NewContext(ctx context.Context) Context { //nolint:ireturn
	return &baseContext{ctx: ctx}
}

type baseContext struct {
	ctx context.Context //nolint
}

func (c *baseContext) SetContext(ctx context.Context) {
	c.ctx = ctx
}

func (c *baseContext) Context() context.Context {
	return c.ctx
}
