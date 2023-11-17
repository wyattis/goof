package goof

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func init() {
	var _ Routable = CRUD(&sqlx.DB{}, struct{}{}, nil)
	var _ Routable = Json("", func(_ *gin.Context, _ struct{}) (r struct{}, s int, e error) { return })
	var _ Routable = ToJson("", func(_ *gin.Context) (r struct{}, s int, e error) { return })
	var _ Routable = FromJson("", func(_ *gin.Context, _ struct{}) (s int, e error) { return })
}
