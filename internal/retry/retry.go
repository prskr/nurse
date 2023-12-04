package retry

import (
	"context"
	"errors"
	"time"
)

// Retry executes a function with the given number of attempts and attempt timeouts.
// It returns the last error encountered during the attempts.
// If the context is canceled, it returns the context error (if there is no previous error),
// or the joined error of the last error and the context error (otherwise).
func Retry(ctx context.Context, numberOfAttempts uint, attemptTimeout time.Duration, f func(ctx context.Context, attempt int) error) (lastErr error) {
	baseCtx, baseCancel := context.WithTimeout(ctx, time.Duration(numberOfAttempts)*attemptTimeout)
	defer baseCancel()

	for i := uint(0); i < numberOfAttempts; i++ {
		select {
		case <-ctx.Done():
			if lastErr == nil {
				return ctx.Err()
			}
			return errors.Join(lastErr, ctx.Err())
		default:
			attemptCtx, attemptCancel := context.WithTimeout(baseCtx, attemptTimeout)

			lastErr = f(attemptCtx, int(i))
			if lastErr == nil {
				attemptCancel()
				return nil
			}

			if attemptCtx.Err() == nil {
				<-attemptCtx.Done()
			}

			attemptCancel()
		}
	}

	return lastErr
}
