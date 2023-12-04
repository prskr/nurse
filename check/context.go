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
		Context:          base,
		attemptTimeout:   attemptTimeout,
		numberOfAttempts: numberOfAttempts,
	}, cancel
}

type checkContext struct {
	attemptTimeout   time.Duration
	numberOfAttempts uint
	context.Context
}

func (c *checkContext) AttemptCount() uint {
	return c.numberOfAttempts
}

func (c *checkContext) AttemptTimeout() time.Duration {
	return c.attemptTimeout
}

func (c *checkContext) WithParent(ctx context.Context) Context {
	return &checkContext{
		Context:        ctx,
		attemptTimeout: c.attemptTimeout,
	}
}
