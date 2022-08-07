package api

import (
	"net/http"
	"time"

	"code.1533b4dc0.de/prskr/nurse/check"
)

var _ http.Handler = (*CheckHandler)(nil)

type CheckHandler struct {
	Timeout  time.Duration
	Attempts uint
	Check    check.SystemChecker
}

func (c CheckHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx, cancel := check.AttemptsContext(request.Context(), c.Attempts, c.Timeout)
	defer cancel()

	if err := c.Check.Execute(ctx); err != nil {
		writer.WriteHeader(http.StatusServiceUnavailable)
		_, _ = writer.Write([]byte(err.Error()))
		return
	}

	writer.WriteHeader(http.StatusOK)
	return
}
