package cmd

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"code.icb4dc0.de/prskr/nurse/api"
	"github.com/urfave/cli/v2"
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
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 100 * time.Millisecond,
	}

	if err := srv.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		logger.Error("Failed to serve HTTP", slog.String("err", err.Error()))
	}

	return nil
}
