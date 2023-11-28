package route

import (
	"fmt"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	Any("/").Handle(func(c *gin.Context) (err error) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
		return
	})
}

type testHasRoleMiddleware struct {
	role string
}

func (m *testHasRoleMiddleware) Name() string {
	return fmt.Sprintf("hasRole(%s)", m.role)
}

func (m *testHasRoleMiddleware) Handler() HandlerFunc {
	return func(c *gin.Context) (err error) {
		c.Next()
		return
	}
}

func TestPrintRoutes(t *testing.T) {
	PrintRoutes(os.Stdout, Group("/api",
		Get("/").Use(&testHasRoleMiddleware{role: "owner"}).Name("return pong"),
		Get("/hey").Use(Or(
			&testHasRoleMiddleware{role: "admin"},
			&testHasRoleMiddleware{role: "owner"},
		)).Name("return hey"),
	), Group("/api/v2").Use(Or(
		&testHasRoleMiddleware{role: "admin"},
		&testHasRoleMiddleware{role: "owner"},
	)).Routes(
		Get("/").Name("return pong"),
		Get("/hey").Name("return hey"),
	))
}
