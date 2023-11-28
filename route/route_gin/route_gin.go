package route_gin

import (
	"github.com/gin-gonic/gin"
	"github.com/wyattis/goof/route"
)

func ginHandler(handler route.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := handler(c)
		if err != nil {
			c.AbortWithError(500, err)
		}
	}
}

// Mount applies the given routes to the gin router provided.
func Mount(router gin.IRouter, routes ...route.IRoute) {
	for _, rInt := range routes {
		route := rInt.Route()
		if route.IsGroup {
			g := router.Group(route.Pattern)
			Mount(g, route.Children...)
		} else {
			handlers := make([]gin.HandlerFunc, 0, len(route.Uses)+1)
			for _, use := range route.Uses {
				handlers = append(handlers, ginHandler(use.Handler()))
			}
			handlers = append(handlers, ginHandler(route.Handler))
			if route.IsAny() {
				// TODO: this isn't really necessary since separate routes are created for each method anyway
				router.Any(route.Pattern, handlers...)
			} else {
				router.Match(route.Methods.Items(), route.Pattern, handlers...)
			}
		}
	}
}
