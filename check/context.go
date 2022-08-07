package check

import (
	"context"
	"time"
)

var _ Context = (*checkContext)(nil)

func AttemptsContext(parent context.Context, numberOfAttempts uint, attemptTimeout time.Duration) (*checkContext, context.CancelFunc) {
	finalTimeout := time.Duration(numberOfAttempts) * attemptTimeout
	base, cancel := context.WithTimeout(parent, finalTimeout)

	return &checkContext{
		Context:        base,
		attemptTimeout: attemptTimeout,
	}, cancel
}

type checkContext struct {
	attemptTimeout time.Duration
	context.Context
}

func (c *checkContext) WithParent(ctx context.Context) Context {
	return &checkContext{
		Context:        ctx,
		attemptTimeout: c.attemptTimeout,
	}
}

func (c *checkContext) AttemptContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(c, c.attemptTimeout)
}
