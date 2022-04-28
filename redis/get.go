package redis

import (
	"context"

	"github.com/go-redis/redis/v8"

	"github.com/baez90/nurse/check"
	"github.com/baez90/nurse/grammar"
)

var (
	_ check.SystemChecker      = (*GetCheck)(nil)
	_ grammar.CheckUnmarshaler = (*GetCheck)(nil)
)

type GetCheck struct {
	redis.UniversalClient
	validators ValidationChain
	Key        string
}

func (g GetCheck) Execute(ctx context.Context) error {
	cmd := g.Get(ctx, g.Key)

	if err := cmd.Err(); err != nil {
		return err
	}

	return g.validators.Validate(cmd)
}

func (g *GetCheck) UnmarshalCheck(c grammar.Check) error {
	const serverAndKeyArgsNumber = 2
	inst := c.Initiator
	if err := grammar.ValidateParameterCount(inst.Params, serverAndKeyArgsNumber); err != nil {
		return err
	}

	var err error
	if g.UniversalClient, err = clientFromParam(inst.Params[0]); err != nil {
		return err
	}

	if g.Key, err = inst.Params[1].AsString(); err != nil {
		return err
	}

	return nil
}
