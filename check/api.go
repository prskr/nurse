package check

import (
	"context"
	"errors"

	"github.com/baez90/nurse/grammar"
)

var ErrNoSuchCheck = errors.New("no such check")

type SystemChecker interface {
	grammar.CheckUnmarshaler
	Execute(ctx context.Context) error
}
