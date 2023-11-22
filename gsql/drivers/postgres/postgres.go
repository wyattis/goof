//go:build all || postgres || !nopostgres
// +build all postgres !nopostgres

package postgres

import (
	"database/sql"
	"fmt"

	"github.com/wyattis/goof/gsql/driver"
	"github.com/wyattis/z/zhtml/ztemplate"
)

func init() {
	driver.Connectors[driver.TypePostgres] = func(config driver.Config) (*sql.DB, error) {
		return connectPostgres(config)
	}
}

func connectPostgres(config driver.Config) (db *sql.DB, err error) {
	if config.Database == "" {
		return nil, fmt.Errorf("database name must be specified to connect to postgres")
	} else if config.User == "" {
		return nil, fmt.Errorf("user must be specified to connect to postgres")
	}

	var source string
	if config.Environment == "cloud run" {
		source, err = ztemplate.ExecString("host={{.SocketDir}}/{{.Host}} port={{.Port}} user={{.User}} password={{.Password}} dbname={{.Database}} sslmode={{.SslMode}}", config)
	} else {
		source, err = ztemplate.ExecString("host={{.Host}} port={{.Port}} user={{.User}} password={{.Password}} dbname={{.Database}} sslmode={{.SslMode}}", config)
	}
	if err != nil {
		return
	}
	return sql.Open(string(driver.TypePostgres), source)
}
