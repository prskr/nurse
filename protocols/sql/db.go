package sql

import (
	"database/sql"
	"fmt"

	"code.1533b4dc0.de/prskr/nurse/config"
	"code.1533b4dc0.de/prskr/nurse/grammar"
)

func dbFromParam(p grammar.Param, srvLookup config.ServerLookup) (*sql.DB, error) {
	if srvName, err := p.AsString(); err != nil {
		return nil, err
	} else if srv, err := srvLookup.Lookup(srvName); err != nil {
		return nil, err
	} else if db, err := DBForServer(srv); err != nil {
		return nil, err
	} else {
		return db, nil
	}
}

func DBForServer(srv *config.Server) (*sql.DB, error) {
	switch srv.Type {
	case config.ServerTypePostgres:
		connectionStrings := srv.ConnectionStrings()
		if lenConnStrings := len(connectionStrings); lenConnStrings != 1 {
			return nil, fmt.Errorf("ambigious number of connection strings: %d", lenConnStrings)
		}

		return sql.Open(srv.Type.Driver(), connectionStrings[0])
	case config.ServerTypeMysql:
		dsns := srv.DSNs()
		if lenConnStrings := len(dsns); lenConnStrings != 1 {
			return nil, fmt.Errorf("ambigious number of DSNs: %d", lenConnStrings)
		}

		return sql.Open(srv.Type.Driver(), dsns[0])
	case config.ServerTypeRedis:
		fallthrough
	default:
		return nil, fmt.Errorf("unmatched server type for SQL DB: %s", srv.Type.Scheme())
	}

}
