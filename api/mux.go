package api

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/baez90/nurse/check"
	"github.com/baez90/nurse/config"
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
			Timeout: spec.Timeout(instance.CheckTimeout),
			Check:   chk,
		})
	}

	return mux, nil
}
