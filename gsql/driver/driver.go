package driver

import "database/sql"

//go:generate go-enum --marshal --flag

// ENUM(postgres, mysql, sqlite3)
type Type string

// ENUM(local, cloud_run)
type Environment string

type Config struct {
	Environment Environment `default:"local"`
	Driver      Type        `default:"sqlite3"`
	SslMode     string      `default:"disable"`
	Host        string      `default:"127.0.0.1"`
	Port        string      `default:"5432"`
	User        string
	Password    string
	Database    string
	SocketDir   string `default:"/cloudsql"`
}

type ConnectionFactory func(config Config) (*sql.DB, error)

var Connectors = map[Type]ConnectionFactory{}

// type SchemaMutatorFactory = func() SchemaMutator

// var Mutators = map[DriverType]SchemaMutatorFactory{}
