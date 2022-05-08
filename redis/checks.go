package redis

import (
	"github.com/baez90/nurse/check"
)

func Module() *check.Module {
	m, _ := check.NewModule(
		check.WithCheck("ping", check.FactoryFunc(func() check.SystemChecker {
			return new(PingCheck)
		})),
		check.WithCheck("get", check.FactoryFunc(func() check.SystemChecker {
			return new(GetCheck)
		})),
	)

	return m
}
