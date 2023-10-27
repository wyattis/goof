package schema

type DriverType string

const (
	DriverTypeSqlite3  DriverType = "sqlite3"
	DriverTypeMysql    DriverType = "mysql"
	DriverTypePostgres DriverType = "postgres"
)

type Constant interface {
	Constant(driver DriverType) string
}

// Constant functions
type NOW struct{}

func (n NOW) Constant(driver DriverType) string {
	switch driver {
	case DriverTypeSqlite3:
		return "CURRENT_TIMESTAMP"
	default:
		panic("unknown driver type")
	}
}
