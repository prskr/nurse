package redis

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"

	"code.icb4dc0.de/prskr/nurse/internal/retry"

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
	logger := slog.Default().With(
		slog.String("check", "redis.GET"),
		slog.String("key", g.Key),
	)

	return retry.Retry(ctx, ctx.AttemptCount(), ctx.AttemptTimeout(), func(ctx context.Context, attempt int) error {
		logger.Debug("Execute check", slog.Int("attempt", attempt))

		cmd := g.Get(ctx, g.Key)

		if err := cmd.Err(); err != nil {
			return err
		}

		return g.validators.Validate(cmd)
	})
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
