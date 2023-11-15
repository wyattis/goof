package sql

import (
	"database/sql"
	"fmt"

	"github.com/wyattis/goof/sql/driver"
	_ "github.com/wyattis/goof/sql/drivers/mysql"
	_ "github.com/wyattis/goof/sql/drivers/postgres"
	_ "github.com/wyattis/goof/sql/drivers/sqlite3"
)

var (
	ErrDriverNotInBuild   = fmt.Errorf("current executable was not built with support for this driver")
	ErrSchemaVersionOlder = fmt.Errorf("schema version is older than the current executable")
	ErrSchemaVersionNewer = fmt.Errorf("schema version is newer than the current executable")
)

func Open(config driver.Config) (db *sql.DB, err error) {
	connector, ok := driver.Connectors[config.DriverName]
	if !ok {
		return nil, ErrDriverNotInBuild
	}
	return connector(config)
}
