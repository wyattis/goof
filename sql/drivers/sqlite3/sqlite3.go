//go:build all || sqlite3 || !nosqlite3
// +build all sqlite3 !nosqlite3

package sqlite3

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/wyattis/goof/sql/driver"
)

func init() {
	fmt.Println("init sqlite3")
	driver.Connectors[driver.TypeSqlite3] = func(config driver.Config) (*sql.DB, error) {
		fmt.Println("connecting to sqlite", config.Database)
		return sql.Open(string(driver.TypeSqlite3), config.Database)
	}
}
