package redis

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"

	"code.icb4dc0.de/prskr/nurse/check"
	"code.icb4dc0.de/prskr/nurse/config"
	"code.icb4dc0.de/prskr/nurse/grammar"
	"code.icb4dc0.de/prskr/nurse/internal/retry"
	"code.icb4dc0.de/prskr/nurse/validation"
)

var _ check.SystemChecker = (*SetCheck)(nil)

type SetCheck struct {
	redis.UniversalClient
	validators validation.Validator[redis.Cmder]
	Key, Value string
}

func (s *SetCheck) Execute(ctx check.Context) error {
	logger := slog.Default().With(
		slog.String("check", "redis.SET"),
		slog.String("key", s.Key),
		slog.String("value", s.Value),
	)

	return retry.Retry(ctx, ctx.AttemptCount(), ctx.AttemptTimeout(), func(ctx context.Context, attempt int) error {
		logger.Debug("Execute check", slog.Int("attempt", attempt))

		cmd := s.Set(ctx, s.Key, s.Value, -1)

		if err := cmd.Err(); err != nil {
			return err
		}

		return s.validators.Validate(cmd)
	})
}

func (s *SetCheck) UnmarshalCheck(c grammar.Check, lookup config.ServerLookup) error {
	const serverKeyAndValueArgsNumber = 3
	inst := c.Initiator
	if err := grammar.ValidateParameterCount(inst.Params, serverKeyAndValueArgsNumber); err != nil {
		return err
	}

	var err error
	if s.UniversalClient, err = clientFromParam(inst.Params[0], lookup); err != nil {
		return err
	}

	if s.Key, err = inst.Params[1].AsString(); err != nil {
		return err
	}

	if s.Value, err = inst.Params[2].AsString(); err != nil {
		return err
	}

	if s.validators, err = registry.ValidatorsForFilters(c.Validators); err != nil {
		return err
	}

	return nil
}
