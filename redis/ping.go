package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"

	"github.com/baez90/nurse/check"
	"github.com/baez90/nurse/grammar"
)

var (
	_ check.SystemChecker      = (*PingCheck)(nil)
	_ grammar.CheckUnmarshaler = (*PingCheck)(nil)
)

type PingCheck struct {
	redis.UniversalClient
	validators ValidationChain
	Message    string
}

func (p PingCheck) Execute(ctx context.Context) error {
	if p.Message == "" {
		return p.Ping(ctx).Err()
	}
	if resp, err := p.Do(ctx, "PING", p.Message).Text(); err != nil {
		return err
	} else if resp != p.Message {
		return fmt.Errorf("expected value %s got %s", p.Message, resp)
	}

	return nil
}

func (p *PingCheck) UnmarshalCheck(c grammar.Check) error {
	const (
		serverOnlyArgCount       = 1
		serverAndMessageArgCount = 2
	)

	p.validators = append(p.validators, StringCmdValidator("PONG"))

	init := c.Initiator
	switch len(init.Params) {
	case 0:
		return grammar.ErrMissingServer
	case serverAndMessageArgCount:
		if msg, err := init.Params[1].AsString(); err != nil {
			return err
		} else {
			p.validators = ValidationChain{StringCmdValidator(msg)}
		}
		fallthrough
	case serverOnlyArgCount:
		if cli, err := clientFromParam(init.Params[0]); err != nil {
			return err
		} else {
			p.UniversalClient = cli
		}
		return nil
	default:
		return grammar.ErrAmbiguousParamCount
	}
}
