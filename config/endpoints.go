package config

import (
	"fmt"
	"path"
	"strings"
	"time"

	"code.1533b4dc0.de/prskr/nurse/grammar"
)

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

func (s EndpointSpec) Timeout(fallback time.Duration) time.Duration {
	if s.CheckTimeout != 0 {
		return s.CheckTimeout
	}

	return fallback
}

func (s *EndpointSpec) Parse(text string) error {
	parser, err := grammar.NewParser[grammar.Script]()
	if err != nil {
		return err
	}

	script, err := parser.Parse(text)
	if err != nil {
		return err
	}

	s.Checks = script.Checks
	return nil
}
