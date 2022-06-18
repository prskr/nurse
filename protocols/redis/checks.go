package redis

import (
	"github.com/baez90/nurse/check"
)

func Module() *check.Module {
	m, err := check.NewModule(
		"redis",
		check.WithCheck("ping", check.FactoryFunc(func() check.SystemChecker {
			return new(PingCheck)
		})),
		check.WithCheck("get", check.FactoryFunc(func() check.SystemChecker {
			return new(GetCheck)
		})),
	)

	if err != nil {
		panic(err)
	}

	return m
}
