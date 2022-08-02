package http

import (
	"net/http"

	"code.1533b4dc0.de/prskr/nurse/validation"
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
