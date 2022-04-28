package redis_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/maxatome/go-testdeep/td"

	"github.com/baez90/nurse/config"
	"github.com/baez90/nurse/grammar"
	"github.com/baez90/nurse/redis"
)

func TestPingCheck_Execute(t *testing.T) {
	t.Parallel()
	srv := PrepareRedisContainer(t)
	if err := config.DefaultLookup.Register(t.Name(), *srv); err != nil {
		t.Fatalf("DefaultLookup.Register() err = %v", err)
	}

	tests := []struct {
		name    string
		check   string
		wantErr bool
	}{
		{
			name:    "PING check",
			check:   fmt.Sprintf(`redis.PING("%s")`, t.Name()),
			wantErr: false,
		},
		// redis.PING("my-redis")
		// redis.PING("my-redis", "Hello, Redis!")
		// redis.GET("my-redis", "some-key") -> String("Hello")
		// redis.PING("my-redis"); redis.GET("my-redis", "some-key") -> String("Hello")
		{
			name:    "PING check - with custom message",
			check:   fmt.Sprintf(`redis.PING("%s", "Hello, Redis!")`, t.Name()),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ping := new(redis.PingCheck)

			parser, err := grammar.NewParser[grammar.Check]()
			td.CmpNoError(t, err, "grammar.NewParser()")
			check, err := parser.Parse(tt.check)
			td.CmpNoError(t, err, "parser.Parse()")

			if err := ping.UnmarshalCheck(*check); err != nil {
				t.Fatalf("UnmarshalCheck() err = %v", err)
			}
			if err := ping.Execute(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
