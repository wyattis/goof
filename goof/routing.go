package goof

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wyattis/goof/route"
)

type PageQuery struct {
	Page    int    `form:"page" binding:"min=0"`
	Size    int    `form:"size" binding:"min=5,max=100"`
	OrderBy string `form:"orderBy"`
	Desc    bool   `form:"desc"`
}

type PayloadInterceptor[T any] func(*gin.Context, T) (T, error)

type PipelineHandler[Payload any, Response any] func(*gin.Context, Payload) (Response, int, error)
type ResponseHandler[Response any] func(*gin.Context) (Response, int, error)
type RequestHandler[Payload any] func(*gin.Context, Payload) (int, error)
type StatusHandler func(*gin.Context) (int, error)

// Takes in a JSON payload and returns a JSON response of the types provided
func Json[Req any, Res any](pattern string, handler PipelineHandler[Req, Res]) *route.RouteBuilder {
	return route.NewRouteBuilder(pattern).Handle(func(c *gin.Context) (err error) {
		var payload Req
		if err = c.BindJSON(&payload); err != nil {
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
		return
	})
}

// Takes in any type of request and responds with a JSON response of the type provided
func ToJson[Res any](pattern string, handler ResponseHandler[Res]) *route.RouteBuilder {
	return route.R(pattern).Handle(func(c *gin.Context) (err error) {
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
		return
	})
}

// Takes in a JSON payload and returns any response
func FromJson[Req any](pattern string, handler RequestHandler[Req]) *route.RouteBuilder {
	return route.R(pattern).Handle(func(c *gin.Context) (err error) {
		var payload Req
		if err = c.BindJSON(&payload); err != nil {
			c.Error(fmt.Errorf("invalid request: %+v", payload))
			return
		}

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
		return
	})
}

// Takes in any type of request and returns a status code
func Status(pattern string, handler StatusHandler) *route.RouteBuilder {
	return route.R(pattern).Handle(func(c *gin.Context) (err error) {
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
		return
	})
}
