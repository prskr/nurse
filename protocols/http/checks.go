package http

import (
	"net/http"

	"github.com/baez90/nurse/check"
)

func Module() *check.Module {
	m, _ := check.NewModule(
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

	return m
}
