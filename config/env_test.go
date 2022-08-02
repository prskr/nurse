package config_test

import (
	"fmt"
	"testing"

	"github.com/maxatome/go-testdeep/td"

	"code.1533b4dc0.de/prskr/nurse/config"
)

//nolint:paralleltest // not possible with env setup
func TestServersFromEnv(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		want    any
		wantErr bool
	}{
		{
			name: "Empty env",
			want: td.NotNil(),
		},
		{
			name: "Single server",
			env: map[string]string{
				fmt.Sprintf("%s_REDIS_1", config.ServerKeyPrefix): "redis://localhost:6379",
			},
			want: td.Map(make(map[string]config.Server), td.MapEntries{
				"redis_1": config.Server{
					Type:  config.ServerTypeRedis,
					Hosts: []string{"localhost:6379"},
					Args:  make(map[string]any),
				},
			}),
			wantErr: false,
		},
		{
			name: "Multiple servers",
			env: map[string]string{
				fmt.Sprintf("%s_REDIS_1", config.ServerKeyPrefix): "redis://localhost:6379",
				fmt.Sprintf("%s_REDIS_2", config.ServerKeyPrefix): "redis://redis:6379",
			},
			want: td.Map(make(map[string]config.Server), td.MapEntries{
				"redis_1": config.Server{
					Type:  config.ServerTypeRedis,
					Hosts: []string{"localhost:6379"},
					Args:  make(map[string]any),
				},
				"redis_2": config.Server{
					Type:  config.ServerTypeRedis,
					Hosts: []string{"redis:6379"},
					Args:  make(map[string]any),
				},
			}),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.env != nil {
				for k, v := range tt.env {
					t.Setenv(k, v)
				}
			}

			got, err := config.ServersFromEnv()
			if (err != nil) != tt.wantErr {
				t.Errorf("ServersFromEnv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			td.Cmp(t, got, tt.want)
		})
	}
}

//nolint:paralleltest // not possible with env setup
func TestEndpointsFromEnv(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		want    any
		wantErr bool
	}{
		{
			name: "Empty env",
			want: td.NotNil(),
		},
		{
			name: "Single endpoint",
			env: map[string]string{
				fmt.Sprintf("%s_readiness", config.EndpointKeyPrefix): `redis.PING("local_redis")`,
			},
			want: td.Map(make(map[config.Route]config.EndpointSpec), td.MapEntries{
				config.Route("readiness"): td.Struct(config.EndpointSpec{}, td.StructFields{
					"Checks": td.Len(1),
				}),
			}),
			wantErr: false,
		},
		{
			name: "Single endpoint - multiple checks",
			env: map[string]string{
				fmt.Sprintf("%s_readiness", config.EndpointKeyPrefix): `redis.PING("local_redis");redis.GET("local_redis", "serving") => String("ok")`,
			},
			want: td.Map(make(map[config.Route]config.EndpointSpec), td.MapEntries{
				config.Route("readiness"): td.Struct(config.EndpointSpec{}, td.StructFields{
					"Checks": td.Len(2),
				}),
			}),
			wantErr: false,
		},
		{
			name: "Single endpoint - sub-route",
			env: map[string]string{
				//nolint:lll // checks might become rather long lines
				fmt.Sprintf("%s_READINESS_REDIS", config.EndpointKeyPrefix): `redis.PING("local_redis");redis.GET("local_redis", "serving") => String("ok")`,
			},
			want: td.Map(make(map[config.Route]config.EndpointSpec), td.MapEntries{
				config.Route("readiness/redis"): td.Struct(config.EndpointSpec{}, td.StructFields{
					"Checks": td.Len(2),
				}),
			}),
			wantErr: false,
		},
		{
			name: "Multiple endpoints",
			env: map[string]string{
				fmt.Sprintf("%s_readiness", config.EndpointKeyPrefix): `redis.PING("local_redis")`,
				fmt.Sprintf("%s_liveness", config.EndpointKeyPrefix):  `redis.GET("local_redis", "serving") => String("ok")`,
			},
			want: td.Map(make(map[config.Route]config.EndpointSpec), td.MapEntries{
				config.Route("readiness"): td.Struct(config.EndpointSpec{}, td.StructFields{
					"Checks": td.Len(1),
				}),
				config.Route("liveness"): td.Struct(config.EndpointSpec{}, td.StructFields{
					"Checks": td.Len(1),
				}),
			}),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.env != nil {
				for k, v := range tt.env {
					t.Setenv(k, v)
				}
			}

			got, err := config.EndpointsFromEnv()
			if (err != nil) != tt.wantErr {
				t.Errorf("EndpointsFromEnv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			td.Cmp(t, got, tt.want)
		})
	}
}
