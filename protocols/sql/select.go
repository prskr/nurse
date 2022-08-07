package sql

import (
	"context"
	"database/sql"

	"go.uber.org/multierr"

	"code.1533b4dc0.de/prskr/nurse/check"
	"code.1533b4dc0.de/prskr/nurse/config"
	"code.1533b4dc0.de/prskr/nurse/grammar"
	"code.1533b4dc0.de/prskr/nurse/validation"
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
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			attemptCtx, cancel := ctx.AttemptContext()
			err := s.executeAttempt(attemptCtx)
			cancel()
			if err == nil {
				return nil
			}
		}
	}
}

func (s *SelectCheck) executeAttempt(ctx context.Context) (err error) {
	var rows *sql.Rows
	rows, err = s.QueryContext(ctx, s.Query)
	if err != nil {
		return err
	}

	defer multierr.AppendInvoke(&err, multierr.Close(rows))
	defer multierr.AppendInvoke(&err, multierr.Invoke(rows.Err))

	return s.validators.Validate(rows)
}
