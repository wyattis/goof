package test

import (
	"net/http/httptest"

	"github.com/jmoiron/sqlx"
	"github.com/wyattis/goof/goof"
	"github.com/wyattis/goof/log"
	"github.com/wyattis/goof/sql/driver"
)

type ctrlFactory func(db *sqlx.DB) goof.Controller
type ginTestModule struct {
	goof.BaseModule
	controller ctrlFactory
}

func (m *ginTestModule) Id() string {
	return "test_ctrl_module"
}
func (m *ginTestModule) Controllers(db *sqlx.DB) []goof.Controller {
	return []goof.Controller{m.controller(db)}
}

// TestGinCtrl creates a test server with a single controller using sqlite.
func TestSqliteGinCtrl(fact func(db *sqlx.DB) goof.Controller) (s *httptest.Server, err error) {
	m := &ginTestModule{controller: fact}
	return TestGinModule(m, goof.RootConfig{
		Log: log.Config{
			Level: "debug",
		},
		DB: driver.Config{
			DriverName: driver.TypeSqlite3,
			Database:   ":memory:",
		},
	})
}

func TestSqliteGinModule(m goof.Module) (s *httptest.Server, err error) {
	return TestGinModule(m, goof.RootConfig{
		DB: driver.Config{
			DriverName: driver.TypeSqlite3,
			Database:   ":memory:",
		},
	})
}

// TestGinModule creates a test server with a single module using the provided config.
func TestGinModule(m goof.Module, config goof.RootConfig) (s *httptest.Server, err error) {
	root := &goof.RootModule{Config: config}
	root.Add(m)
	if err = root.Init(); err != nil {
		return
	}
	s = httptest.NewServer(root.Engine())
	return
}
