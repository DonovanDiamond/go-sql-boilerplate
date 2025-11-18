package boilerplate

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

func Connect(driver Driver, connString string) (*sqlx.DB, error) {
	return sqlx.Connect(string(driver), connString)
}
