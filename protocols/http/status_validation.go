package http

import (
	"fmt"
	"net/http"

	"code.icb4dc0.de/prskr/nurse/grammar"
	"code.icb4dc0.de/prskr/nurse/validation"
)

var _ validation.FromCall[*http.Response] = (*StatusValidator)(nil)

type StatusValidator struct {
	Want int
}

func (s *StatusValidator) Validate(resp *http.Response) error {
	if resp.StatusCode != s.Want {
		return fmt.Errorf("want HTTP status %d but got %d", s.Want, resp.StatusCode)
	}

	return nil
}

func (s *StatusValidator) UnmarshalCall(c grammar.Call) error {
	if err := grammar.ValidateParameterCount(c.Params, 1); err != nil {
		return err
	}

	var err error
	s.Want, err = c.Params[0].AsInt()

	return err
}
