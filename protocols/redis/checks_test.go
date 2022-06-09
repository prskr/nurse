package redis_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	redisCli "github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/maxatome/go-testdeep/td"

	"github.com/baez90/nurse/config"
	"github.com/baez90/nurse/grammar"
	"github.com/baez90/nurse/protocols/redis"
)

func TestChecks_Execute(t *testing.T) {
	t.Parallel()

	redisModule := redis.Module()

	tests := []struct {
		name    string
		check   string
		setup   func(tb testing.TB, cli redisCli.UniversalClient)
		wantErr bool
	}{
		{
			name:  "Get value",
			check: `redis.GET("%s", "some_key")`,
			setup: func(tb testing.TB, cli redisCli.UniversalClient) {
				tb.Helper()
				td.CmpNoError(tb, cli.Set(context.Background(), "some_key", "some_value", 0).Err())
			},
			wantErr: false,
		},
		{
			name:  "Get value - validate value",
			check: `redis.GET("%s", "some_key") => Equals("some_value")`,
			setup: func(tb testing.TB, cli redisCli.UniversalClient) {
				tb.Helper()
				td.CmpNoError(tb, cli.Set(context.Background(), "some_key", "some_value", 0).Err())
			},
			wantErr: false,
		},
		{
			name:    "PING check",
			check:   `redis.PING("%s")`,
			wantErr: false,
		},
		{
			name:    "PING check - with custom message",
			check:   `redis.PING("%s", "Hello, Redis!")`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			srv := PrepareRedisContainer(t)
			serverName := uuid.NewString()

			register := config.NewServerRegister()

			if err := register.Register(serverName, *srv); err != nil {
				t.Fatalf("DefaultLookup.Register() err = %v", err)
			}

			cli, err := redis.ClientForServer(srv)
			if err != nil {
				t.Fatalf("redis.ClientForServer() err = %v", err)
			}

			if tt.setup != nil {
				tt.setup(t, cli)
			}

			if strings.Contains(tt.check, "%s") {
				tt.check = fmt.Sprintf(tt.check, serverName)
			}

			parser, err := grammar.NewParser[grammar.Check]()
			td.CmpNoError(t, err, "grammar.NewParser()")
			parsedCheck, err := parser.Parse(tt.check)
			td.CmpNoError(t, err, "parser.Parse()")

			chk, err := redisModule.Lookup(*parsedCheck, register)
			td.CmpNoError(t, err, "redis.LookupCheck()")

			if tt.wantErr {
				td.CmpError(t, chk.Execute(context.Background()))
			} else {
				td.CmpNoError(t, chk.Execute(context.Background()))
			}
		})
	}
}
