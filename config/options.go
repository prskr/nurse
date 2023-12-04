package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//nolint:ireturn // required for interface implementation
func WithCheckDuration(d time.Duration) Option {
	return OptionFunc(func(n *Nurse) error {
		n.CheckTimeout = d

		return nil
	})
}

func WithCheckAttempts(attempts uint) Option {
	return OptionFunc(func(n *Nurse) error {
		n.CheckAttempts = attempts
		return nil
	})
}

//nolint:ireturn // required for interface implementation
func WithServersFromArgs(servers []string) Option {
	return OptionFunc(func(n *Nurse) error {
		if len(servers) == 0 {
			return nil
		}

		for _, rawSrv := range servers {
			name, rawURL, found := strings.Cut(rawSrv, "=")
			if !found {
				return fmt.Errorf("couldn't parse %s as server, expected format: <server-name>=<url>", rawSrv)
			}

			var srv Server
			if err := srv.UnmarshalURL(rawURL); err != nil {
				return err
			}

			name = strings.ToLower(strings.TrimSpace(name))

			n.Servers[name] = srv
		}

		return nil
	})
}

//nolint:ireturn // required for interface implementation
func WithServersFromEnv() Option {
	return OptionFunc(func(n *Nurse) error {
		envServers, err := ServersFromEnv()
		if err != nil {
			return err
		}

		if n.Servers == nil || len(n.Servers) == 0 {
			n.Servers = envServers
			return nil
		}

		for name, srv := range envServers {
			if _, ok := n.Servers[name]; !ok {
				n.Servers[name] = srv
			}
		}

		return nil
	})
}

func WithEndpointsFromEnv() Option {
	return OptionFunc(func(n *Nurse) error {
		envEndpoints, err := EndpointsFromEnv()
		if err != nil {
			return err
		}

		if n.Endpoints == nil || len(n.Endpoints) == 0 {
			n.Endpoints = envEndpoints
			return nil
		}

		for route, spec := range envEndpoints {
			if _, ok := n.Endpoints[route]; !ok {
				n.Endpoints[route] = spec
			}
		}

		return nil
	})
}

//nolint:ireturn // required to implement interface
func WithConfigFile(configFilePath string) Option {
	return OptionFunc(func(n *Nurse) error {
		var out Nurse
		if configFilePath != "" {
			slog.Debug("Attempt to load config file")
			if err := out.ReadFromFile(configFilePath); err == nil {
				return nil
			} else {
				slog.Warn(
					"Failed to load config file",
					slog.String("config_file_path", configFilePath),
					slog.String("err", err.Error()),
				)
			}
		}

		if workingDir, err := os.Getwd(); err == nil {
			configFilePath = filepath.Join(workingDir, "nurse.yaml")
			slog.Debug("Attempt to load config file from current working directory")
			if err = out.ReadFromFile(configFilePath); err == nil {
				return nil
			} else {
				slog.Warn(
					"Failed to load config file",
					slog.String("config_file_path", configFilePath),
					slog.String("err", err.Error()),
				)
			}
		}

		if home, err := os.UserHomeDir(); err == nil {
			configFilePath = filepath.Join(home, ".nurse.yaml")
			slog.Debug("Attempt to load config file from user home directory")
			if err = out.ReadFromFile(configFilePath); err == nil {
				return nil
			} else {
				slog.Warn(
					"Failed to load config file",
					slog.String("config_file_path", configFilePath),
					slog.String("err", err.Error()),
				)
			}
		}

		configFilePath = filepath.Join("", "etc", "nurse", "config.yaml")
		slog.Debug("Attempt to load config file from global config directory")
		if err := out.ReadFromFile(configFilePath); err == nil {
			return nil
		} else {
			slog.Warn(
				"Failed to load config file",
				slog.String("config_file_path", configFilePath),
				slog.String("err", err.Error()),
			)
		}

		return nil
	})
}
