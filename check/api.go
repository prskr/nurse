package check

import (
	"context"
	"errors"

	"code.1533b4dc0.de/prskr/nurse/config"
	"code.1533b4dc0.de/prskr/nurse/grammar"
)

var (
	ErrNoSuchCheck      = errors.New("no such check")
	ErrConflictingCheck = errors.New("check with same name already registered")
	ErrNoSuchValidator  = errors.New("no such validator")
)

type (
	Unmarshaler interface {
		UnmarshalCheck(c grammar.Check, lookup config.ServerLookup) error
	}

	SystemChecker interface {
		Unmarshaler
		Execute(ctx context.Context) error
	}

	CallUnmarshaler interface {
		UnmarshalCall(c grammar.Call) error
	}

	CheckerLookup interface {
		Lookup(c grammar.Check, srvLookup config.ServerLookup) (SystemChecker, error)
	}

	ModuleLookup interface {
		Lookup(modName string) (CheckerLookup, error)
	}
)
