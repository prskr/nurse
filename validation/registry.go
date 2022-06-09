package validation

import (
	"fmt"
	"strings"
	"sync"

	"github.com/baez90/nurse/check"
	"github.com/baez90/nurse/grammar"
)

type (
	Validator[T any] interface {
		Validate(in T) error
	}

	FromCall[T any] interface {
		Validator[T]
		check.CallUnmarshaler
	}

	Chain[T any] []Validator[T]
)

func (c Chain[T]) Validate(in T) error {
	for i := range c {
		if err := c[i].Validate(in); err != nil {
			return err
		}
	}

	return nil
}

func NewRegistry[R any]() *Registry[R] {
	return &Registry[R]{
		validators: make(map[string]func() FromCall[R]),
	}
}

type Registry[R any] struct {
	lock       sync.Mutex
	validators map[string]func() FromCall[R]
}

func (r *Registry[R]) Register(name string, factory func() FromCall[R]) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.validators[strings.ToLower(name)] = factory
}

func (r *Registry[R]) ValidatorsForFilters(filters *grammar.Filters) (Chain[R], error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if filters == nil || filters.Chain == nil {
		return Chain[R]{}, nil
	}
	chain := make(Chain[R], 0, len(filters.Chain))
	for i := range filters.Chain {
		validationCall := filters.Chain[i]
		if validatorProvider, ok := r.validators[strings.ToLower(validationCall.Name)]; !ok {
			return nil, fmt.Errorf("%w: %s", check.ErrNoSuchValidator, validationCall.Name)
		} else {
			validator := validatorProvider()
			if err := validator.UnmarshalCall(validationCall); err != nil {
				return nil, err
			}
			chain = append(chain, validator)
		}
	}

	return chain, nil
}
