package redis

import (
	"context"

	"github.com/go-redis/redis/v8"

	"code.icb4dc0.de/prskr/nurse/check"
	"code.icb4dc0.de/prskr/nurse/config"
	"code.icb4dc0.de/prskr/nurse/grammar"
	"code.icb4dc0.de/prskr/nurse/validation"
)

var _ check.SystemChecker = (*GetCheck)(nil)

type GetCheck struct {
	redis.UniversalClient
	validators validation.Validator[redis.Cmder]
	Key        string
}

func (g *GetCheck) Execute(ctx check.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			attemptCtx, cancel := ctx.AttemptContext()
			err := g.executeAttempt(attemptCtx)
			cancel()
			if err == nil {
				return nil
			}
		}
	}
}

func (g *GetCheck) executeAttempt(ctx context.Context) error {
	cmd := g.Get(ctx, g.Key)

	if err := cmd.Err(); err != nil {
		return err
	}

	return g.validators.Validate(cmd)
}

func (g *GetCheck) UnmarshalCheck(c grammar.Check, lookup config.ServerLookup) error {
	const serverAndKeyArgsNumber = 2
	inst := c.Initiator
	if err := grammar.ValidateParameterCount(inst.Params, serverAndKeyArgsNumber); err != nil {
		return err
	}

	var err error
	if g.UniversalClient, err = clientFromParam(inst.Params[0], lookup); err != nil {
		return err
	}

	if g.Key, err = inst.Params[1].AsString(); err != nil {
		return err
	}

	if g.validators, err = registry.ValidatorsForFilters(c.Validators); err != nil {
		return err
	}

	return nil
}
