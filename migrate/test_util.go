package migrate

import (
	"database/sql"
	"goof/schema"
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
