package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/urfave/cli/v2"

	"code.icb4dc0.de/prskr/nurse/check"
	"code.icb4dc0.de/prskr/nurse/config"
	"code.icb4dc0.de/prskr/nurse/protocols/http"
	"code.icb4dc0.de/prskr/nurse/protocols/redis"
	"code.icb4dc0.de/prskr/nurse/protocols/sql"
)

const (
	defaultCheckTimeout = 500 * time.Millisecond
	defaultAttemptCount = 20
)

const (
	logLevelFlag          = "log.level"
	httpAddressFlag       = "http.address"
	httpReadHeaderTimeout = "http.read-header-timeout"
	maxCheckAttemptsFlag  = "check-attempts"
	checkTimeoutFlag      = "check-timeout"
	serversFlag           = "servers"
	configFlag            = "config"
)

func NewApp() (*cli.App, error) {
	app := &app{
		registry: check.NewRegistry(),
	}

	if err := app.registry.Register(
		redis.Module(),
		http.Module(),
		sql.Module(),
	); err != nil {
		return nil, err
	}

	srv := server{
		app: app,
	}
	exec := executor{
		app: app,
	}

	return &cli.App{
		Name:                 "nurse",
		DefaultCommand:       "server",
		EnableBashCompletion: true,
		Before:               app.init,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    configFlag,
				Usage:   "Config file to load, if not set `$HOME/.nurse.yaml`, `/etc/nurse/config.yaml` and `./nurse.yaml` are tried  - optional",
				Aliases: []string{"c"},
				EnvVars: []string{"NURSE_CONFIG"},
			},
			&cli.DurationFlag{
				Name:    checkTimeoutFlag,
				Usage:   "Timeout when running checks",
				Value:   defaultCheckTimeout,
				EnvVars: []string{"NURSE_CHECK_TIMEOUT"},
			},
			&cli.UintFlag{
				Name:    maxCheckAttemptsFlag,
				Usage:   "Number of attempts for a check",
				Value:   defaultAttemptCount,
				EnvVars: []string{"NURSE_CHECK_ATTEMPTS"},
			},
			&cli.StringFlag{
				Name:  logLevelFlag,
				Usage: "Log level to use",
				Value: "info",
			},
			&cli.StringSliceFlag{
				Name:    serversFlag,
				Usage:   "",
				Aliases: []string{"s"},
			},
		},
		Commands: []*cli.Command{
			{
				Name:     "server",
				Aliases:  []string{"serve"},
				Category: "daemon",
				Action:   srv.RunServer,
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:    "endpoints",
						Usage:   "Endpoints to expose in the HTTP server",
						Aliases: []string{"ep"},
					},
					&cli.StringFlag{
						Name:    httpAddressFlag,
						Usage:   "HTTP server address",
						Value:   ":8080",
						EnvVars: []string{"NURSE_HTTP_ADDRESS"},
					},
					&cli.DurationFlag{
						Name:    httpReadHeaderTimeout,
						Usage:   "Timeout for reading headers in the HTTP server",
						Value:   100 * time.Millisecond,
						EnvVars: []string{"NURSE_HTTP_READ_HEADER_TIMEOUT"},
					},
				},
			},
			{
				Name:      "exec-check",
				Aliases:   []string{"run-check"},
				Category:  "interactive",
				Action:    exec.ExecChecks,
				ArgsUsage: "checks to execute, either in a single argument (within \"\") separated by a ';' or multiple arguments",
			},
		},
	}, nil
}

type app struct {
	nurseInstance *config.Nurse
	registry      *check.Registry
	lookup        *config.ServerRegister
	logging       struct {
		level slog.LevelVar
	}
}

func (a *app) init(ctx *cli.Context) (err error) {
	if err = a.configureLogging(ctx.String(logLevelFlag)); err != nil {
		return err
	}

	a.nurseInstance, err = config.New(
		config.WithCheckAttempts(ctx.Uint(maxCheckAttemptsFlag)),
		config.WithCheckDuration(ctx.Duration(checkTimeoutFlag)),
		config.WithConfigFile(ctx.String(configFlag)),
		config.WithServersFromEnv(),
		config.WithServersFromArgs(ctx.StringSlice(serversFlag)),
		config.WithEndpointsFromEnv(),
	)

	if err != nil {
		return fmt.Errorf("failed to load Nurse config: %w", err)
	}

	if a.lookup, err = a.nurseInstance.ServerLookup(); err != nil {
		return fmt.Errorf("failed to constructor server lookup: %w", err)
	}
	return nil
}

func (a *app) configureLogging(level string) error {
	if err := a.logging.level.UnmarshalText([]byte(level)); err != nil {
		return err
	}

	var addSource bool
	if a.logging.level.Level() == slog.LevelDebug {
		addSource = true
	}

	cfg := slog.HandlerOptions{
		AddSource: addSource,
		Level:     a.logging.level.Level(),
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &cfg)))
	return nil
}
