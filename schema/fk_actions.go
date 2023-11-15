package schema

import "github.com/wyattis/goof/sql/driver"

type FkAction interface {
	Action(driverType driver.Type) string
}

type NO_ACTION struct{}

func (n NO_ACTION) Action(driverType driver.Type) string {
	switch driverType {
	case driver.TypeSqlite3:
		return "NO ACTION"
	default:
		panic("unknown driver type")
	}
}

type RESTRICT struct{}

func (n RESTRICT) Action(driverType driver.Type) string {
	switch driverType {
	case driver.TypeSqlite3:
		return "RESTRICT"
	default:
		panic("unknown driver type")
	}
}

type SET_NULL struct{}

func (n SET_NULL) Action(driverType driver.Type) string {
	switch driverType {
	case driver.TypeSqlite3:
		return "SET NULL"
	default:
		panic("unknown driver type")
	}
}

type SET_DEFAULT struct{}

func (n SET_DEFAULT) Action(driverType driver.Type) string {
	switch driverType {
	case driver.TypeSqlite3:
		return "SET DEFAULT"
	default:
		panic("unknown driver type")
	}
}

type CASCADE struct{}

func (n CASCADE) Action(driverType driver.Type) string {
	switch driverType {
	case driver.TypeSqlite3:
		return "CASCADE"
	default:
		panic("unknown driver type")
	}
}
