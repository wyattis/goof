package goof

import "github.com/gin-gonic/gin"

type PipelineHandler[Payload any, Response any] func(*gin.Context, Payload) (Response, int, error)
type ResponseHandler[Response any] func(*gin.Context) (Response, int, error)
type RequestHandler[Payload any] func(*gin.Context, Payload) (int, error)

func ToJson[Response any](cb ResponseHandler[Response]) gin.HandlerFunc {
	return func(c *gin.Context) {
		response, status, err := cb(c)
		if err != nil {
			if status == 0 {
				status = 500
			}
			c.AbortWithError(status, err)
			return
		}
		c.JSON(status, response)
	}
}

func JsonToJson[Payload any, Response any](cb PipelineHandler[Payload, Response]) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload Payload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.AbortWithError(400, err)
			return
		}
		// TODO: add validation logic
		response, status, err := cb(c, payload)
		if err != nil {
			if status == 0 {
				status = 500
			}
			c.AbortWithError(status, err)
			return
		}
		c.JSON(status, response)
	}
}
