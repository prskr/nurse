package config

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

//nolint:ireturn // required for interface implementation
func WithServersFromEnv() Option {
	return OptionFunc(func(n Nurse) (Nurse, error) {
		envServers, err := ServersFromEnv()
		if err != nil {
			return Nurse{}, err
		}

		if n.Servers == nil || len(n.Servers) == 0 {
			n.Servers = envServers
			return n, nil
		}

		for name, srv := range envServers {
			if _, ok := n.Servers[name]; !ok {
				n.Servers[name] = srv
			}
		}

		return n, nil
	})
}

func WithEndpointsFromEnv() Option {
	return OptionFunc(func(n Nurse) (Nurse, error) {
		envEndpoints, err := EndpointsFromEnv()
		if err != nil {
			return Nurse{}, err
		}

		if n.Endpoints == nil || len(n.Endpoints) == 0 {
			n.Endpoints = envEndpoints
			return n, nil
		}

		for route, spec := range envEndpoints {
			if _, ok := n.Endpoints[route]; !ok {
				n.Endpoints[route] = spec
			}
		}

		return n, nil
	})
}

func WithValuesFrom(other Nurse) Option {
	return OptionFunc(func(n Nurse) (Nurse, error) {
		return n.Merge(other), nil
	})
}

//nolint:ireturn // required to implement interface
func WithConfigFile(configFilePath string) Option {
	logger := zap.L()
	return OptionFunc(func(n Nurse) (Nurse, error) {
		var out Nurse
		if configFilePath != "" {
			logger.Debug("Attempt to load config file")
			if err := out.ReadFromFile(configFilePath); err == nil {
				return out, nil
			} else {
				logger.Warn(
					"Failed to load config file",
					zap.String("config_file_path", configFilePath),
					zap.Error(err),
				)
			}
		}

		if workingDir, err := os.Getwd(); err == nil {
			configFilePath = filepath.Join(workingDir, "nurse.yaml")
			logger.Debug("Attempt to load config file from current working directory")
			if err = out.ReadFromFile(configFilePath); err == nil {
				return out, nil
			} else {
				logger.Warn(
					"Failed to load config file",
					zap.String("config_file_path", configFilePath),
					zap.Error(err),
				)
			}
		}

		if home, err := os.UserHomeDir(); err == nil {
			configFilePath = filepath.Join(home, ".nurse.yaml")
			logger.Debug("Attempt to load config file from user home directory")
			if err = out.ReadFromFile(configFilePath); err == nil {
				return out, nil
			} else {
				logger.Warn(
					"Failed to load config file",
					zap.String("config_file_path", configFilePath),
					zap.Error(err),
				)
			}
		}

		configFilePath = filepath.Join("", "etc", "nurse", "config.yaml")
		logger.Debug("Attempt to load config file from global config directory")
		if err := out.ReadFromFile(configFilePath); err == nil {
			return out, nil
		} else {
			logger.Warn(
				"Failed to load config file",
				zap.String("config_file_path", configFilePath),
				zap.Error(err),
			)
		}

		return Nurse{}, nil
	})
}
