package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"code.1533b4dc0.de/prskr/nurse/api"
	"code.1533b4dc0.de/prskr/nurse/check"
	"code.1533b4dc0.de/prskr/nurse/config"
	"code.1533b4dc0.de/prskr/nurse/protocols/redis"
)

var (
	logLevel = zapcore.InfoLevel
	cfgFile  string
	cfg      config.Nurse
)

func main() {
	if err := prepareFlags(); err != nil {
		log.Fatalf("Failed to parse flags: %v", err)
	}
	setupLogging()

	logger := zap.L()

	nurseInstance, err := config.New(
		config.WithValuesFrom(cfg),
		config.WithConfigFile(cfgFile),
		config.WithServersFromEnv(),
		config.WithEndpointsFromEnv(),
	)
	if err != nil {
		logger.Fatal("Failed to load config from environment", zap.Error(err))
	}

	logger.Debug("Loaded config", zap.Any("config", nurseInstance))

	chkRegistry := check.NewRegistry()
	if err = chkRegistry.Register(redis.Module()); err != nil {
		logger.Fatal("Failed to register Redis module", zap.Error(err))
	}

	srvLookup, err := nurseInstance.ServerLookup()
	if err != nil {
		logger.Fatal("Failed to prepare server lookup", zap.Error(err))
	}

	mux, err := api.PrepareMux(nurseInstance, chkRegistry, srvLookup)
	if err != nil {
		logger.Fatal("Failed to prepare server mux", zap.Error(err))
	}

	if err := http.ListenAndServe(":8080", mux); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return
		}

		logger.Fatal("Failed to serve HTTP", zap.Error(err))
	}
}

func setupLogging() {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(logLevel)
	logger, err := cfg.Build()
	if err != nil {
		log.Fatalf("Failed to setup logging: %v", err)
	}

	zap.ReplaceGlobals(logger)
}

func prepareFlags() error {
	set := config.ConfigureFlags(&cfg)

	set.Var(&logLevel, "log-level", "Log level to use")

	set.StringVar(
		&cfgFile,
		"config",
		config.LookupEnvOr[string]("NURSE_CONFIG", "", config.Identity[string]),
		"Config file to load, if not set $HOME/.nurse.yaml, /etc/nurse/config.yaml and ./nurse.yaml are tried  - optional",
	)

	return set.Parse(os.Args)
}
