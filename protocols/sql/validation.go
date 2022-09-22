package sql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"

	"code.icb4dc0.de/prskr/nurse/validation"
)

var registry = validation.NewRegistry[*sql.Rows]()

func init() {
	registry.Register("rows", func() validation.FromCall[*sql.Rows] {
		return new(RowsValidator)
	})
}
