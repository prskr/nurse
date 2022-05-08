package redis

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-redis/redis/v8"

	"github.com/baez90/nurse/check"
	"github.com/baez90/nurse/grammar"
)

var (
	ErrNoSuchValidator = errors.New("no such validator")

	_ CmdValidator = (ValidationChain)(nil)
	_ CmdValidator = (*StringCmdValidator)(nil)

	knownValidators = map[string]func() unmarshallableCmdValidator{
		"string": func() unmarshallableCmdValidator {
			return new(StringCmdValidator)
		},
	}
)

type (
	CmdValidator interface {
		Validate(cmder redis.Cmder) error
	}
	unmarshallableCmdValidator interface {
		CmdValidator
		check.CallUnmarshaler
	}
)

func ValidatorsForFilters(filters *grammar.Filters) (ValidationChain, error) {
	if filters == nil || filters.Chain == nil {
		return ValidationChain{}, nil
	}
	chain := make(ValidationChain, 0, len(filters.Chain))
	for i := range filters.Chain {
		validationCall := filters.Chain[i]
		if validatorProvider, ok := knownValidators[strings.ToLower(validationCall.Name)]; !ok {
			return nil, fmt.Errorf("%w: %s", ErrNoSuchValidator, validationCall.Name)
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

type ValidationChain []CmdValidator

func (v ValidationChain) UnmarshalCall(grammar.Call) error {
	return errors.New("cannot unmarshal chain")
}

func (v ValidationChain) Validate(cmder redis.Cmder) error {
	for i := range v {
		if err := v[i].Validate(cmder); err != nil {
			return err
		}
	}

	return nil
}

type StringCmdValidator string

func (s *StringCmdValidator) UnmarshalCall(c grammar.Call) error {
	if err := grammar.ValidateParameterCount(c.Params, 1); err != nil {
		return err
	}

	want, err := c.Params[0].AsString()
	if err != nil {
		return err
	}

	*s = StringCmdValidator(want)
	return nil
}

func (s StringCmdValidator) Validate(cmder redis.Cmder) error {
	if err := cmder.Err(); err != nil {
		return err
	}

	if stringCmd, ok := cmder.(*redis.StringCmd); !ok {
		return errors.New("not a string result")
	} else if got, err := stringCmd.Result(); err != nil {
		return err
	} else if want := string(s); got != want {
		return fmt.Errorf("want %s but got %s", want, got)
	}

	return nil
}
