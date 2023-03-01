package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"code.icb4dc0.de/prskr/nurse/check"
	"code.icb4dc0.de/prskr/nurse/config"
	"code.icb4dc0.de/prskr/nurse/grammar"
	"code.icb4dc0.de/prskr/nurse/validation"
)

var _ check.SystemChecker = (*PingCheck)(nil)

type PingCheck struct {
	redis.UniversalClient
	validators validation.Validator[redis.Cmder]
	Message    string
}

func (p PingCheck) Execute(ctx check.Context) error {
	if p.Message == "" {
		return p.Ping(ctx).Err()
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			attemptCtx, cancel := ctx.AttemptContext()
			err := p.executeAttempt(attemptCtx)
			cancel()
			if err == nil {
				return nil
			}
		}
	}
}

func (p PingCheck) executeAttempt(ctx context.Context) error {
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
