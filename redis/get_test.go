package redis_test

import (
	"context"
	"fmt"
	"testing"

	redisCli "github.com/go-redis/redis/v8"
	"github.com/maxatome/go-testdeep/td"

	"github.com/baez90/nurse/config"
	"github.com/baez90/nurse/grammar"
	"github.com/baez90/nurse/redis"
)

func TestGetCheck_Execute(t *testing.T) {
	t.Parallel()
	srv := PrepareRedisContainer(t)
	if err := config.DefaultLookup.Register(t.Name(), *srv); err != nil {
		t.Fatalf("DefaultLookup.Register() err = %v", err)
	}

	cli, err := redis.ClientForServer(srv)
	if err != nil {
		t.Fatalf("redis.ClientForServer() err = %v", err)
	}

	tests := []struct {
		name    string
		check   string
		setup   func(tb testing.TB, cli redisCli.UniversalClient)
		wantErr bool
	}{
		{
			name:  "Get value",
			check: fmt.Sprintf(`redis.GET("%s", "some_key")`, t.Name()),
			setup: func(tb testing.TB, cli redisCli.UniversalClient) {
				tb.Helper()
				td.CmpNoError(tb, cli.Set(context.Background(), "some_key", "some_value", 0).Err())
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.setup != nil {
				tt.setup(t, cli)
			}

			get := new(redis.GetCheck)

			parser, err := grammar.NewParser[grammar.Check]()
			td.CmpNoError(t, err, "grammar.NewParser()")
			check, err := parser.Parse(tt.check)
			td.CmpNoError(t, err, "parser.Parse()")

			td.CmpNoError(t, get.UnmarshalCheck(*check), "get.UnmarshalCheck()")
			td.CmpNoError(t, get.Execute(context.Background()))

			if tt.wantErr {
				td.CmpError(t, get.Execute(context.Background()))
			} else {
				td.CmpNoError(t, get.Execute(context.Background()))
			}
		})
	}
}
