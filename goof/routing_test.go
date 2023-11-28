package goof

import (
	"github.com/gin-gonic/gin"
	"github.com/wyattis/goof/route"
)

func init() {
	var _ route.IRoute = Json("", func(_ *gin.Context, _ struct{}) (r struct{}, s int, e error) { return })
	var _ route.IRoute = ToJson("", func(_ *gin.Context) (r struct{}, s int, e error) { return })
	var _ route.IRoute = FromJson("", func(_ *gin.Context, _ struct{}) (s int, e error) { return })
}
