package http

import (
	"net/http"

	"github.com/baez90/nurse/validation"
)

var registry = validation.NewRegistry[*http.Response]()

func init() {
	registry.Register("jsonpath", func() validation.FromCall[*http.Response] {
		return new(JSONPathValidator)
	})

	registry.Register("status", func() validation.FromCall[*http.Response] {
		return new(StatusValidator)
	})
}
