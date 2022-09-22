package sql_test

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/td"

	"code.icb4dc0.de/prskr/nurse/check"
	"code.icb4dc0.de/prskr/nurse/config"
	"code.icb4dc0.de/prskr/nurse/grammar"
	sqlchk "code.icb4dc0.de/prskr/nurse/protocols/sql"
)

func TestChecks_Execute(t *testing.T) {
	t.Parallel()

	td.DefaultContextConfig.FailureIsFatal = true
	sqlModule := sqlchk.Module()

	dbTypes := []config.ServerType{
		config.ServerTypePostgres,
		config.ServerTypeMysql,
	}

	tests := []struct {
		name    string
		check   string
		setup   func(tb testing.TB, db *sql.DB)
		wantErr bool
	}{
		{
			name:  "Simple SELECT 1",
			check: `sql.SELECT("%s", "SELECT 1;")`,
		},
		{
			name:  "Simple SELECT 1 with column name",
			check: `sql.SELECT("%s", "SELECT 1 as Idx;")`,
		},
		{
			name:  "Simple SELECT 1 - check for rows",
			check: `sql.SELECT("%s", "SELECT 1;") => Rows(1)`,
		},
	}
	for _, tt := range tests {
		tt := tt
		for _, st := range dbTypes {
			st := st
			t.Run(fmt.Sprintf("%s: %s", strings.ToUpper(st.Scheme()), tt.name), func(t *testing.T) {
				t.Parallel()

				var (
					srv     *config.Server
					srvName string
				)

				switch st {
				case config.ServerTypeMysql:
					srvName, srv = PrepareMariaDBContainer(t)
				case config.ServerTypePostgres:
					srvName, srv = PreparePostgresContainer(t)
				case config.ServerTypeRedis:
					fallthrough
				default:
					t.Fatalf("unexpected server type: %s", st.Scheme())
				}

				register := config.NewServerRegister()

				td.CmpNoError(t, register.Register(srvName, *srv))

				db, err := sqlchk.DBForServer(srv)
				td.CmpNoError(t, err, "sql.DBForServer()")

				if tt.setup != nil {
					tt.setup(t, db)
				}

				rawCheck := strings.Clone(tt.check)

				if strings.Contains(rawCheck, "%s") {
					rawCheck = fmt.Sprintf(rawCheck, srvName)
				}

				parser, err := grammar.NewParser[grammar.Check]()
				td.CmpNoError(t, err, "grammar.NewParser()")
				parsedCheck, err := parser.Parse(rawCheck)
				td.CmpNoError(t, err, "parser.Parse()")

				chk, err := sqlModule.Lookup(*parsedCheck, register)
				td.CmpNoError(t, err, "redis.LookupCheck()")

				ctx, cancel := check.AttemptsContext(context.Background(), 100, 500*time.Millisecond)
				t.Cleanup(cancel)

				if tt.wantErr {
					td.CmpError(t, chk.Execute(ctx))
				} else {
					td.CmpNoError(t, chk.Execute(ctx))
				}
			})
		}
	}
}
