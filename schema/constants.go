package schema

import "github.com/wyattis/goof/sql/driver"

type Constant interface {
	Constant(driver driver.Type) string
}

// Constant functions
type NOW struct{}

func (n NOW) Constant(driverType driver.Type) string {
	switch driverType {
	case driver.TypeSqlite3:
		return "CURRENT_TIMESTAMP"
	default:
		panic("unsupported driver type")
	}
}

type CURRENT_TIMESTAMP = NOW
