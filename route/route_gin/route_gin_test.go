package route_gin

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/wyattis/goof/route"
)

func TestMount(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	Mount(r, route.Any("/user").Handle(func(c *gin.Context) (err error) {
		c.String(200, "hello")
		return
	}))
	res := r.Routes()[0]
	if res.Path != "/user" {
		t.Errorf("expected /user, got %s", res.Path)
	}
}
