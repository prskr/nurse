package api

import (
	"log/slog"
	"net/http"

	"code.icb4dc0.de/prskr/nurse/check"
	"code.icb4dc0.de/prskr/nurse/config"
)

func PrepareMux(instance *config.Nurse, modLookup check.ModuleLookup, srvLookup config.ServerLookup) (http.Handler, error) {
	mux := http.NewServeMux()

	for route, spec := range instance.Endpoints {
		slog.Info("Configuring route", slog.String("route", route.String()))
		chk, err := check.CheckForScript(spec.Checks, modLookup, srvLookup)
		if err != nil {
			return nil, err
		}

		mux.Handle(route.String(), CheckHandler{
			Timeout:  spec.Timeout(instance.CheckTimeout),
			Attempts: spec.Attempts(instance.CheckAttempts),
			Check:    chk,
		})
	}

	return mux, nil
}
