package http

import (
	"net/http"

	"code.icb4dc0.de/prskr/nurse/check"
)

func Module() *check.Module {
	m, err := check.NewModule(
		"http",
		check.WithCheck("get", check.FactoryFunc(func() check.SystemChecker {
			return &GenericCheck{Method: http.MethodGet}
		})),
		check.WithCheck("post", check.FactoryFunc(func() check.SystemChecker {
			return &GenericCheck{Method: http.MethodPost}
		})),
		check.WithCheck("put", check.FactoryFunc(func() check.SystemChecker {
			return &GenericCheck{Method: http.MethodPut}
		})),
		check.WithCheck("delete", check.FactoryFunc(func() check.SystemChecker {
			return &GenericCheck{Method: http.MethodDelete}
		})),
	)
	if err != nil {
		panic(err)
	}

	return m
}
