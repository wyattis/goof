package route

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

type Middleware interface {
	Name() string
	Handler() HandlerFunc
}

type orMiddleware struct {
	middlewares []Middleware
}

func (m *orMiddleware) Name() string {
	names := make([]string, len(m.middlewares))
	for i, middleware := range m.middlewares {
		names[i] = middleware.Name()
	}
	return fmt.Sprintf("OR(%s)", strings.Join(names, ", "))
}

func (m *orMiddleware) Handler() HandlerFunc {
	handlers := make([]HandlerFunc, len(m.middlewares))
	for i, middleware := range m.middlewares {
		handlers[i] = middleware.Handler()
	}
	return func(c *gin.Context) (err error) {
		potentialErrors := []error{}
		successful := false
		for _, handler := range handlers {
			if err = handler(c); err == nil {
				successful = true
				break
			} else {
				potentialErrors = append(potentialErrors, err)
			}
		}
		if !successful {
			for _, err := range potentialErrors {
				c.Error(err)
			}
			c.Abort()
		}
		return
	}
}

func Or(middlewares ...Middleware) Middleware {
	return &orMiddleware{middlewares}
}
