package sql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/baez90/nurse/validation"
)

var registry = validation.NewRegistry[*sql.Rows]()

func init() {
	registry.Register("rows", func() validation.FromCall[*sql.Rows] {
		return new(RowsValidator)
	})
}
