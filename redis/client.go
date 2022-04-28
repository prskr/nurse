package redis

import (
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/mitchellh/mapstructure"

	"github.com/baez90/nurse/config"
	"github.com/baez90/nurse/grammar"
)

func clientFromParam(p grammar.Param) (redis.UniversalClient, error) {
	if srvName, err := p.AsString(); err != nil {
		return nil, err
	} else if srv, err := config.DefaultLookup.Lookup(srvName); err != nil {
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

	return redis.NewUniversalClient(opts), nil
}
