package sql

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"code.icb4dc0.de/prskr/nurse/check"
	"code.icb4dc0.de/prskr/nurse/config"
	"code.icb4dc0.de/prskr/nurse/grammar"
	"code.icb4dc0.de/prskr/nurse/internal/retry"
	"code.icb4dc0.de/prskr/nurse/validation"
)

var _ check.SystemChecker = (*SelectCheck)(nil)

type SelectCheck struct {
	*sql.DB
	validators validation.Validator[*sql.Rows]
	Query      string
}

func (s *SelectCheck) UnmarshalCheck(c grammar.Check, lookup config.ServerLookup) error {
	const serverAndKeyArgsNumber = 2
	inst := c.Initiator
	if err := grammar.ValidateParameterCount(inst.Params, serverAndKeyArgsNumber); err != nil {
		return err
	}

	var err error
	if s.DB, err = dbFromParam(inst.Params[0], lookup); err != nil {
		return err
	}

	if s.Query, err = inst.Params[1].AsString(); err != nil {
		return err
	}

	if s.validators, err = registry.ValidatorsForFilters(c.Validators); err != nil {
		return err
	}

	return nil
}

func (s *SelectCheck) Execute(ctx check.Context) error {
	logger := slog.Default().With(
		slog.String("check", "sql.SELECT"),
		slog.String("query", s.Query),
	)

	return retry.Retry(ctx, ctx.AttemptCount(), ctx.AttemptTimeout(), func(ctx context.Context, attempt int) error {
		logger.Debug("Execute check", slog.Int("attempt", attempt))

		logger.Debug("")

		rows, err := s.QueryContext(ctx, s.Query)
		if err != nil {
			return err
		}

		defer func() {
			err = errors.Join(rows.Close(), rows.Err())
		}()

		return s.validators.Validate(rows)
	})
}
