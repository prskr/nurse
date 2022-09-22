package redis

import (
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/mitchellh/mapstructure"

	"code.icb4dc0.de/prskr/nurse/config"
	"code.icb4dc0.de/prskr/nurse/grammar"
)

func clientFromParam(p grammar.Param, srvLookup config.ServerLookup) (redis.UniversalClient, error) {
	if srvName, err := p.AsString(); err != nil {
		return nil, err
	} else if srv, err := srvLookup.Lookup(srvName); err != nil {
		return nil, err
	} else if redisCli, err := ClientForServer(srv); err != nil {
		return nil, err
	} else {
		return redisCli, nil
	}
}

func ClientForServer(srv *config.Server) (redis.UniversalClient, error) {
	opts := &redis.UniversalOptions{
		Addrs: srv.Hosts,
	}

	if pathLen := len(srv.Path); pathLen > 0 {
		if db, err := strconv.Atoi(srv.Path[0]); err == nil {
			opts.DB = db
		}
	}

	if err := mapstructure.Decode(srv.Args, opts); err != nil {
		return nil, err
	}

	if srv.Credentials != nil {
		opts.Username = srv.Credentials.Username
		opts.Password = *srv.Credentials.Password
	}

	return redis.NewUniversalClient(opts), nil
}
