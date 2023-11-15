package migrate

import (
	"database/sql"

	"github.com/wyattis/goof/sql/driver"
)

var sqliteDbFile = ":memory:"

func setupSqlite() (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", sqliteDbFile)
	if err != nil {
		return
	}
	err = initializeSchema(db, driver.TypeSqlite3, "test")
	return
}
