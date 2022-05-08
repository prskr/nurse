package config

import (
	"encoding"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/baez90/nurse/grammar"
)

var _ encoding.TextUnmarshaler = (*EndpointSpec)(nil)

type Route string

func (r Route) String() string {
	val := string(r)
	val = strings.Trim(val, "/")

	return path.Clean(fmt.Sprintf("/%s", val))
}

type EndpointSpec struct {
	CheckTimeout time.Duration
	Checks       []grammar.Check
}

func (e *EndpointSpec) UnmarshalText(text []byte) error {
	parser, err := grammar.NewParser[grammar.Script]()
	if err != nil {
		return err
	}

	script, err := parser.Parse(string(text))
	if err != nil {
		return err
	}

	e.Checks = script.Checks
	return nil
}
