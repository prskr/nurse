package http

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"

	"code.icb4dc0.de/prskr/nurse/check"
	"code.icb4dc0.de/prskr/nurse/config"
	"code.icb4dc0.de/prskr/nurse/grammar"
	"code.icb4dc0.de/prskr/nurse/validation"
)

type ClientInjectable interface {
	SetClient(client *http.Client)
}

var (
	_ check.SystemChecker = (*GenericCheck)(nil)
	_ ClientInjectable    = (*GenericCheck)(nil)
)

type GenericCheck struct {
	*http.Client
	validators validation.Validator[*http.Response]
	Method     string
	Body       []byte
	URL        string
}

func (g *GenericCheck) SetClient(client *http.Client) {
	if client == nil {
		return
	}

	g.Client = client
}

func (g *GenericCheck) Execute(ctx check.Context) error {
	slog.Default().Debug("Execute check",
		slog.String("check", "http"),
		slog.String("method", g.Method),
		slog.String("url", g.URL),
	)

	var body io.Reader
	if len(g.Body) > 0 {
		body = bytes.NewReader(g.Body)
	}
	req, err := http.NewRequestWithContext(ctx, g.Method, g.URL, body)
	if err != nil {
		return err
	}

	resp, err := g.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	return g.validators.Validate(resp)
}

func (g *GenericCheck) UnmarshalCheck(c grammar.Check, _ config.ServerLookup) error {
	const urlArgsNumber = 1

	inst := c.Initiator

	if err := grammar.ValidateParameterCount(inst.Params, urlArgsNumber); err != nil {
		return err
	}

	if g.Client == nil {
		g.Client = http.DefaultClient
	}

	var err error
	if g.URL, err = inst.Params[0].AsString(); err != nil {
		return err
	}

	if g.validators, err = registry.ValidatorsForFilters(c.Validators); err != nil {
		return err
	}

	return nil
}
