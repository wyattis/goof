package goof

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

type PipelineHandler[Payload any, Response any] func(*gin.Context, Payload) (Response, int, error)
type ResponseHandler[Response any] func(*gin.Context) (Response, int, error)
type RequestHandler[Payload any] func(*gin.Context, Payload) (int, error)
type StatusHandler func(*gin.Context) (int, error)

type IRoute interface {
	Method() string
	Pattern() string
	Uses() []gin.HandlerFunc
	Handler() gin.HandlerFunc

	RequestType() any
	ResponseType() any
}

type Routable interface {
	Routes() []IRoute
}

type routeBuilder[Req any, Res any] struct {
	route *route[Req, Res]
}

type route[Req any, Res any] struct {
	method  string
	pattern string
	use     []gin.HandlerFunc
	handler gin.HandlerFunc
}

func (b *route[Req, Res]) Method() string {
	return b.method
}

func (b *route[Req, Res]) Pattern() string {
	return b.pattern
}

func (b *route[Req, Res]) Uses() []gin.HandlerFunc {
	return b.use
}

func (b *route[Req, Res]) Handler() gin.HandlerFunc {
	return b.handler
}

func (b *route[Req, Res]) RequestType() any {
	var req Req
	return req
}

func (b *route[Req, Res]) ResponseType() any {
	var res Res
	return res
}

func (b *routeBuilder[Req, Res]) Routes() []IRoute {
	return []IRoute{b.route}
}

func (b *routeBuilder[Req, Res]) Get() *routeBuilder[Req, Res] {
	b.route.method = http.MethodGet
	return b
}

func (b *routeBuilder[Req, Res]) Post() *routeBuilder[Req, Res] {
	b.route.method = http.MethodPost
	return b
}

func (b *routeBuilder[Req, Res]) Put() *routeBuilder[Req, Res] {
	b.route.method = http.MethodPut
	return b
}

func (b *routeBuilder[Req, Res]) Delete() *routeBuilder[Req, Res] {
	b.route.method = http.MethodDelete
	return b
}

func (b *routeBuilder[Req, Res]) Patch() *routeBuilder[Req, Res] {
	b.route.method = http.MethodPatch
	return b
}

func (b *routeBuilder[Req, Res]) Options() *routeBuilder[Req, Res] {
	b.route.method = http.MethodOptions
	return b
}

func (b *routeBuilder[Req, Res]) Head() *routeBuilder[Req, Res] {
	b.route.method = http.MethodHead
	return b
}

func (b *routeBuilder[Req, Res]) Use(handlers ...gin.HandlerFunc) *routeBuilder[Req, Res] {
	b.route.use = append(b.route.use, handlers...)
	return b
}

// Takes in a JSON payload and returns a JSON response of the types provided
func Json[Req any, Res any](pattern string, handler PipelineHandler[Req, Res]) *routeBuilder[Req, Res] {
	return &routeBuilder[Req, Res]{
		route: &route[Req, Res]{
			pattern: pattern,
			handler: func(c *gin.Context) {
				var payload Req
				if err := c.BindJSON(&payload); err != nil {
					c.Error(fmt.Errorf("invalid request: %+v", payload))
					return
				}

				response, status, err := handler(c, payload)
				if err != nil {
					if status == 0 {
						status = http.StatusInternalServerError
					}
					c.Error(fmt.Errorf("handler failure: %+v", payload))
					c.AbortWithError(status, err)
					return
				}
				if status == 0 {
					status = http.StatusOK
				}
				c.JSON(status, response)
			},
		},
	}
}

// Takes in any type of request and responds with a JSON response of the type provided
func ToJson[Res any](pattern string, handler ResponseHandler[Res]) *routeBuilder[any, Res] {
	return &routeBuilder[any, Res]{
		route: &route[any, Res]{
			pattern: pattern,
			handler: func(c *gin.Context) {
				response, status, err := handler(c)
				if err != nil {
					if status == 0 {
						status = http.StatusInternalServerError
					}
					c.Error(fmt.Errorf("handler failure"))
					c.AbortWithError(status, err)
					return
				}
				if status == 0 {
					status = http.StatusOK
				}
				c.JSON(status, response)
			},
		},
	}
}

// Takes in a JSON payload and returns any response
func FromJson[Req any](pattern string, handler RequestHandler[Req]) *routeBuilder[Req, any] {
	return &routeBuilder[Req, any]{
		route: &route[Req, any]{
			pattern: pattern,
			handler: func(c *gin.Context) {
				var payload Req
				if err := c.BindJSON(&payload); err != nil {
					c.Error(fmt.Errorf("invalid request: %+v", payload))
					return
				}

				// Struct level validation
				// errs := validate.Struct(payload)
				// if len(errs) > 0 {
				// 	c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid request: %+v", payload))
				// 	for _, err := range errs {
				// 		c.Error(err)
				// 	}
				// 	return
				// }

				status, err := handler(c, payload)
				if err != nil {
					if status == 0 {
						status = http.StatusInternalServerError
					}
					c.Error(fmt.Errorf("handler failure: %+v", payload))
					c.AbortWithError(status, err)
					return
				}
				if status == 0 {
					status = http.StatusOK
				}
				c.Status(status)
			},
		},
	}
}

// Takes in any type of request and returns any response
func Route(pattern string, handler gin.HandlerFunc) *routeBuilder[struct{}, struct{}] {
	return &routeBuilder[struct{}, struct{}]{
		route: &route[struct{}, struct{}]{
			pattern: pattern,
			handler: handler,
		},
	}
}

// Takes in any type of request and returns a status code
func Status(pattern string, handler StatusHandler) *routeBuilder[struct{}, struct{}] {
	return &routeBuilder[struct{}, struct{}]{
		route: &route[struct{}, struct{}]{
			pattern: pattern,
			handler: func(c *gin.Context) {
				status, err := handler(c)
				if err != nil {
					if status == 0 {
						status = http.StatusInternalServerError
					}
					c.AbortWithError(status, err)
					return
				}
				if status == 0 {
					status = http.StatusOK
				}
				c.Status(status)
			},
		},
	}
}

// Register several gin routes at once
func RouteGin(router gin.IRouter, routes ...Routable) {
	for _, r := range routes {
		routes := r.Routes()
		for _, r := range routes {
			router.Handle(r.Method(), r.Pattern(), r.Handler()).Use(r.Uses()...)
		}
	}
}

// TODO: In theory it's possible to generate a TypeScript API, right?
func TSRouter(routes ...Routable) {
	for _, route := range routes {
		routes := route.Routes()
		for _, r := range routes {
			println(r.Method(), r.Pattern(), reflect.TypeOf(route))
		}
	}
}
