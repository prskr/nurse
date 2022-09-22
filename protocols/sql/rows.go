package sql

import (
	"database/sql"
	"fmt"

	"code.icb4dc0.de/prskr/nurse/grammar"
	"code.icb4dc0.de/prskr/nurse/validation"
)

var _ validation.FromCall[*sql.Rows] = (*RowsValidator)(nil)

type RowsValidator struct {
	Want int
}

func (r *RowsValidator) Validate(in *sql.Rows) error {
	readRows := 0
	for in.Next() {
		readRows++
	}

	if err := in.Err(); err != nil {
		return err
	}

	if readRows != r.Want {
		return fmt.Errorf("expected %d rows but got %d", r.Want, readRows)
	}

	return nil
}

func (r *RowsValidator) UnmarshalCall(c grammar.Call) error {
	if err := grammar.ValidateParameterCount(c.Params, 1); err != nil {
		return err
	}

	var err error
	r.Want, err = c.Params[0].AsInt()

	return err
}
