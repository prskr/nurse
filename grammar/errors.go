package grammar

import "errors"

var (
	ErrMissingServer       = errors.New("initiator is missing a server")
	ErrTypeMismatch        = errors.New("param has a different type")
	ErrAmbiguousParamCount = errors.New("the supplied number of arguments does not match the expected one")
)
