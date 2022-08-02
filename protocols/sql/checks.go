package sql

import "code.1533b4dc0.de/prskr/nurse/check"

func Module() *check.Module {
	m, err := check.NewModule(
		"sql",
		check.WithCheck("select", check.FactoryFunc(func() check.SystemChecker {
			return new(SelectCheck)
		})),
	)

	if err != nil {
		panic(err)
	}

	return m
}
