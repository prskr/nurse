package redis

import (
	"code.icb4dc0.de/prskr/nurse/check"
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
		check.WithCheck("set", check.FactoryFunc(func() check.SystemChecker {
			return new(SetCheck)
		})),
	)
	if err != nil {
		panic(err)
	}

	return m
}
