package sql

import "github.com/baez90/nurse/check"

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
