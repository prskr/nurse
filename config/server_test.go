package config_test

import (
	"testing"

	"github.com/maxatome/go-testdeep/td"

	"code.icb4dc0.de/prskr/nurse/config"
	"code.icb4dc0.de/prskr/nurse/internal/values"
)

func TestParseFromURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		url     string
		want    any
		wantErr bool
	}{
		{
			name: "Single Redis server",
			url:  "redis://localhost:6379/0",
			want: &config.Server{
				Type:  config.ServerTypeRedis,
				Hosts: []string{"localhost:6379"},
				Path:  []string{"0"},
				Args:  make(map[string]any),
			},
		},
		{
			name: "Single Redis server with args",
			url:  "redis://localhost:6379/0?MaxRetries=3",
			want: &config.Server{
				Type:  config.ServerTypeRedis,
				Hosts: []string{"localhost:6379"},
				Path:  []string{"0"},
				Args: map[string]any{
					"MaxRetries": 3.0,
				},
			},
		},
		{
			name: "Single Redis server with user",
			url:  "redis://user@localhost:6379/0",
			want: &config.Server{
				Type:  config.ServerTypeRedis,
				Hosts: []string{"localhost:6379"},
				Path:  []string{"0"},
				Args:  make(map[string]any),
				Credentials: &config.Credentials{
					Username: "user",
				},
			},
		},
		{
			name: "Single Redis server with credentials",
			url:  "redis://user:p4$$w0rd@localhost:6379/0",
			want: &config.Server{
				Type:  config.ServerTypeRedis,
				Hosts: []string{"localhost:6379"},
				Path:  []string{"0"},
				Args:  make(map[string]any),
				Credentials: &config.Credentials{
					Username: "user",
					Password: values.StringP("p4$$w0rd"),
				},
			},
		},
		{
			name: "Multiple Redis servers",
			url:  "redis://{localhost:6379,localhost:6380}/0",
			want: &config.Server{
				Type: config.ServerTypeRedis,
				Hosts: []string{
					"localhost:6379",
					"localhost:6380",
				},
				Path: []string{"0"},
				Args: make(map[string]any),
			},
		},
		{
			name: "Multiple Redis servers with args",
			url:  "redis://{localhost:6379,localhost:6380}/0?MaxRetries=3",
			want: &config.Server{
				Type: config.ServerTypeRedis,
				Hosts: []string{
					"localhost:6379",
					"localhost:6380",
				},
				Path: []string{"0"},
				Args: map[string]any{
					"MaxRetries": 3.0,
				},
			},
		},
		{
			name: "Multiple Redis servers with user",
			url:  "redis://user@{localhost:6379,localhost:6380}/0",
			want: &config.Server{
				Type: config.ServerTypeRedis,
				Hosts: []string{
					"localhost:6379",
					"localhost:6380",
				},
				Credentials: &config.Credentials{
					Username: "user",
				},
				Path: []string{"0"},
				Args: make(map[string]any),
			},
		},
		{
			name: "Multiple Redis servers with credentials",
			url:  "redis://user:p4$$w0rd@{localhost:6379,localhost:6380}/0",
			want: &config.Server{
				Type: config.ServerTypeRedis,
				Hosts: []string{
					"localhost:6379",
					"localhost:6380",
				},
				Credentials: &config.Credentials{
					Username: "user",
					Password: values.StringP("p4$$w0rd"),
				},
				Path: []string{"0"},
				Args: make(map[string]any),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := new(config.Server)

			if err := got.UnmarshalURL(tt.url); (err != nil) != tt.wantErr {
				t.Errorf("ParseFromURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			td.Cmp(t, got, tt.want)
		})
	}
}
