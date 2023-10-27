package schema

type FkAction interface {
	Action(driver DriverType) string
}

type NO_ACTION struct{}

func (n NO_ACTION) Action(driver DriverType) string {
	switch driver {
	case DriverTypeSqlite3:
		return "NO ACTION"
	default:
		panic("unknown driver type")
	}
}

type RESTRICT struct{}

func (n RESTRICT) Action(driver DriverType) string {
	switch driver {
	case DriverTypeSqlite3:
		return "RESTRICT"
	default:
		panic("unknown driver type")
	}
}

type SET_NULL struct{}

func (n SET_NULL) Action(driver DriverType) string {
	switch driver {
	case DriverTypeSqlite3:
		return "SET NULL"
	default:
		panic("unknown driver type")
	}
}

type SET_DEFAULT struct{}

func (n SET_DEFAULT) Action(driver DriverType) string {
	switch driver {
	case DriverTypeSqlite3:
		return "SET DEFAULT"
	default:
		panic("unknown driver type")
	}
}

type CASCADE struct{}

func (n CASCADE) Action(driver DriverType) string {
	switch driver {
	case DriverTypeSqlite3:
		return "CASCADE"
	default:
		panic("unknown driver type")
	}
}
