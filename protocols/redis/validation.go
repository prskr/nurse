package redis

import (
	"errors"

	"github.com/go-redis/redis/v8"

	"github.com/baez90/nurse/grammar"
	"github.com/baez90/nurse/validation"
)

var (
	_ CmdValidator = (*GenericCmdValidator)(nil)

	registry = validation.NewRegistry[redis.Cmder]()
)

func init() {
	registry.Register("equals", func() validation.FromCall[redis.Cmder] {
		return new(GenericCmdValidator)
	})
}

type CmdValidator interface {
	Validate(cmder redis.Cmder) error
}

func GenericCommandValidatorFor[T validation.Value](want T) (*GenericCmdValidator, error) {
	comparator, err := validation.JSONValueComparatorFor(want)
	if err != nil {
		return nil, err
	}
	return &GenericCmdValidator{
		comparator: comparator,
	}, nil
}

type GenericCmdValidator struct {
	comparator validation.ValueComparator
}

func (g *GenericCmdValidator) UnmarshalCall(c grammar.Call) error {
	if err := grammar.ValidateParameterCount(c.Params, 1); err != nil {
		return err
	}

	var err error
	switch c.Params[0].Type() {
	case grammar.ParamTypeInt:
		if g.comparator, err = validation.JSONValueComparatorFor(*c.Params[0].Int); err != nil {
			return err
		}
	case grammar.ParamTypeFloat:
		if g.comparator, err = validation.JSONValueComparatorFor(*c.Params[0].Float); err != nil {
			return err
		}
	case grammar.ParamTypeString:
		if g.comparator, err = validation.JSONValueComparatorFor(*c.Params[0].String); err != nil {
			return err
		}
	case grammar.ParamTypeUnknown:
		fallthrough
	default:
		return errors.New("param type is unknown")
	}

	return nil
}

func (g *GenericCmdValidator) Validate(cmder redis.Cmder) error {
	if err := cmder.Err(); err != nil {
		return err
	}

	if in, ok := cmder.(*redis.StringCmd); ok {
		res, err := in.Result()
		if err != nil {
			return err
		}
		return g.comparator.Equals(res)
	}

	return nil
}
