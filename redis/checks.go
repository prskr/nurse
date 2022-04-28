package redis

import (
	"fmt"
	"strings"

	"github.com/baez90/nurse/check"
	"github.com/baez90/nurse/grammar"
)

var knownChecks = map[string]func() check.SystemChecker{
	"ping": func() check.SystemChecker {
		return new(PingCheck)
	},
	"get": func() check.SystemChecker {
		return new(GetCheck)
	},
}

func LookupCheck(c grammar.Check) (check.SystemChecker, error) {
	var (
		provider func() check.SystemChecker
		ok       bool
	)
	if provider, ok = knownChecks[strings.ToLower(c.Initiator.Name)]; !ok {
		return nil, fmt.Errorf("%w: %s", check.ErrNoSuchCheck, c.Initiator.Name)
	}

	chk := provider()
	if err := chk.UnmarshalCheck(c); err != nil {
		return nil, err
	}

	return chk, nil
}
