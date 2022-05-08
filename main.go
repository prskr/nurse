package main

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/baez90/nurse/config"
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

	envCfg, err := config.New(
		config.WithValuesFrom(cfg),
		config.WithConfigFile(cfgFile),
		config.WithServersFromEnv(),
		config.WithEndpointsFromEnv(),
	)

	if err != nil {
		logger.Error("Failed to load config from environment", zap.Error(err))
		os.Exit(1)
	}

	logger.Debug("Loaded config", zap.Any("config", envCfg))
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
