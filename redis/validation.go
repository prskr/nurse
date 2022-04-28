package redis

import (
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var (
	_ CmdValidator = (ValidationChain)(nil)
	_ CmdValidator = StringCmdValidator("")
)

type (
	CmdValidator interface {
		Validate(cmder redis.Cmder) error
	}

	ValidationChain []CmdValidator
)

func (v ValidationChain) Validate(cmder redis.Cmder) error {
	for i := range v {
		if err := v[i].Validate(cmder); err != nil {
			return err
		}
	}

	return nil
}

type StringCmdValidator string

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
