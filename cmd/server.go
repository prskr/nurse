package cmd

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"

	"github.com/urfave/cli/v2"

	"code.icb4dc0.de/prskr/nurse/api"
)

type server struct {
	*app
}

func (a *server) RunServer(ctx *cli.Context) error {
	logger := slog.Default()

	mux, err := api.PrepareMux(a.nurseInstance, a.registry, a.lookup)
	if err != nil {
		logger.Error("Failed to prepare server mux", slog.String("err", err.Error()))
	}

	srv := http.Server{
		Addr:              ctx.String(httpAddressFlag),
		ReadHeaderTimeout: ctx.Duration(httpReadHeaderTimeout),
		BaseContext: func(listener net.Listener) context.Context {
			return ctx.Context
		},
		Handler: mux,
	}

	if err := srv.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		logger.Error("Failed to serve HTTP", slog.String("err", err.Error()))
	}

	return nil
}
