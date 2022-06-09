package http

import (
	"errors"
	"io"
	"net/http"

	"github.com/valyala/bytebufferpool"

	"github.com/baez90/nurse/grammar"
	"github.com/baez90/nurse/validation"
)

var _ validation.FromCall[*http.Response] = (*JSONPathValidator)(nil)

type JSONPathValidator struct {
	validator *validation.JSONPathValidator
}

func (j *JSONPathValidator) UnmarshalCall(c grammar.Call) (err error) {
	if err = grammar.ValidateParameterCount(c.Params, 2); err != nil {
		return err
	}

	var jsonPath string

	if jsonPath, err = c.Params[0].AsString(); err != nil {
		return err
	}

	switch c.Params[1].Type() {
	case grammar.ParamTypeInt:
		j.validator, err = validation.JSONPathValidatorFor(jsonPath, *c.Params[1].Int)
	case grammar.ParamTypeFloat:
		j.validator, err = validation.JSONPathValidatorFor(jsonPath, *c.Params[1].Float)
	case grammar.ParamTypeString:
		j.validator, err = validation.JSONPathValidatorFor(jsonPath, *c.Params[1].String)
	case grammar.ParamTypeUnknown:
		fallthrough
	default:
		return errors.New("param type unknown")
	}

	return nil
}

func (j *JSONPathValidator) Validate(resp *http.Response) error {
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	readBytes, err := io.Copy(buf, resp.Body)
	if err != nil {
		return err
	}

	return j.validator.Equals(buf.B[:readBytes])
}
