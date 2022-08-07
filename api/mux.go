package api

import (
	"net/http"

	"go.uber.org/zap"

	"code.1533b4dc0.de/prskr/nurse/check"
	"code.1533b4dc0.de/prskr/nurse/config"
)

func PrepareMux(instance *config.Nurse, modLookup check.ModuleLookup, srvLookup config.ServerLookup) (http.Handler, error) {
	mux := http.NewServeMux()

	logger := zap.L()

	for route, spec := range instance.Endpoints {
		logger.Info("Configuring route", zap.String("route", route.String()))
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
