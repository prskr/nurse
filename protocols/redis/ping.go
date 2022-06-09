package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"

	"github.com/baez90/nurse/check"
	"github.com/baez90/nurse/config"
	"github.com/baez90/nurse/grammar"
	"github.com/baez90/nurse/validation"
)

var _ check.SystemChecker = (*PingCheck)(nil)

type PingCheck struct {
	redis.UniversalClient
	validators validation.Validator[redis.Cmder]
	Message    string
}

func (p PingCheck) Execute(ctx context.Context) error {
	if p.Message == "" {
		return p.Ping(ctx).Err()
	}
	if resp, err := p.Do(ctx, "PING", p.Message).Text(); err != nil {
		return err
	} else if resp != p.Message {
		return fmt.Errorf("expected value %s got %s", p.Message, resp)
	}

	return nil
}

func (p *PingCheck) UnmarshalCheck(c grammar.Check, lookup config.ServerLookup) error {
	const (
		serverOnlyArgCount       = 1
		serverAndMessageArgCount = 2
	)

	val, _ := GenericCommandValidatorFor("PONG")

	validators := validation.Chain[redis.Cmder]{}
	validators = append(validators, val)

	p.validators = validators

	init := c.Initiator
	switch len(init.Params) {
	case 0:
		return grammar.ErrMissingServer
	case serverAndMessageArgCount:
		if msg, err := init.Params[1].AsString(); err != nil {
			return err
		} else {
			val, _ := GenericCommandValidatorFor(msg)
			p.validators = validation.Chain[redis.Cmder]{val}
		}
		fallthrough
	case serverOnlyArgCount:
		if cli, err := clientFromParam(init.Params[0], lookup); err != nil {
			return err
		} else {
			p.UniversalClient = cli
		}
		return nil
	default:
		return grammar.ErrAmbiguousParamCount
	}
}
