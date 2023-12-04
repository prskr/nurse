package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"code.icb4dc0.de/prskr/nurse/check"
	"code.icb4dc0.de/prskr/nurse/config"
	"code.icb4dc0.de/prskr/nurse/protocols/http"
	"code.icb4dc0.de/prskr/nurse/protocols/redis"
	"code.icb4dc0.de/prskr/nurse/protocols/sql"
	"github.com/urfave/cli/v2"
)

const (
	defaultCheckTimeout = 500 * time.Millisecond
	defaultAttemptCount = 20
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
				Name:    "config",
				Usage:   "Config file to load, if not set `$HOME/.nurse.yaml`, `/etc/nurse/config.yaml` and `./nurse.yaml` are tried  - optional",
				Aliases: []string{"c"},
				EnvVars: []string{"NURSE_CONFIG"},
			},
			&cli.DurationFlag{
				Name:    "check-timeout",
				Usage:   "Timeout when running checks",
				Value:   defaultCheckTimeout,
				EnvVars: []string{"NURSE_CHECK_TIMEOUT"},
			},
			&cli.UintFlag{
				Name:    "check-attempts",
				Usage:   "Number of attempts for a check",
				Value:   defaultAttemptCount,
				EnvVars: []string{"NURSE_CHECK_ATTEMPTS"},
			},
			&cli.StringFlag{
				Name:  "log-level",
				Usage: "Log level to use",
				Value: "info",
			},
			&cli.StringSliceFlag{
				Name:    "servers",
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
						Usage:   "",
						Aliases: []string{"ep"},
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
	if err = a.configureLogging(ctx.String("log-level")); err != nil {
		return err
	}

	a.nurseInstance, err = config.New(
		config.WithCheckAttempts(ctx.Uint("check-attempts")),
		config.WithCheckDuration(ctx.Duration("check-timeout")),
		config.WithConfigFile(ctx.String("config")),
		config.WithServersFromEnv(),
		config.WithServersFromArgs(ctx.StringSlice("servers")),
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
