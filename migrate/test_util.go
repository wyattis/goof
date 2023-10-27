package migrate

import (
	"database/sql"

	"github.com/wyattis/goof/schema"
)

var sqliteDbFile = ":memory:"

func setupSqlite() (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", sqliteDbFile)
	if err != nil {
		return
	}
	err = initializeSchema(db, schema.DriverTypeSqlite3, "test")
	return
}
