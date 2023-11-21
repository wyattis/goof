package test

import (
	"net/http/httptest"

	"github.com/wyattis/goof/goof"
	"github.com/wyattis/goof/log"
	"github.com/wyattis/goof/sql/driver"
)

func TestSqliteGinModule(modules ...goof.Module) (s *httptest.Server, err error) {
	return TestGinModule(goof.RootConfig{
		Production: false,
		Log: log.Config{
			Level: "debug",
		},
		DB: driver.Config{
			DriverName: driver.TypeSqlite3,
			Database:   ":memory:",
		},
	}, modules...)
}

// TestGinModule creates a test server with a single module using the provided config.
func TestGinModule(config goof.RootConfig, modules ...goof.Module) (s *httptest.Server, err error) {
	root := &goof.RootModule{Config: config}
	root.Add(modules...)
	if err = root.Init(); err != nil {
		return
	}
	s = httptest.NewServer(root.Engine())
	return
}
