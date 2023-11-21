package middleware

import (
	"time"

	"github.com/wyattis/goof/log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func applyRequestEvent(e *zerolog.Event, c *gin.Context) *zerolog.Event {
	return e.Str("method", c.Request.Method).
		Str("path", c.Request.URL.Path).
		Int("status", c.Writer.Status())
}

func Log() gin.HandlerFunc {
	return func(c *gin.Context) {
		if log.Debug().Enabled() {
			startTime := time.Now()
			requestId := uuid.NewString()
			c.Set("requestId", requestId)
			log.Debug().Str("request", requestId).Msg("start")
			c.Next()
			var event *zerolog.Event
			if len(c.Errors) > 0 {
				event = log.Error()
			} else {
				event = log.Debug()
			}
			applyRequestEvent(event, c).
				Str("request", requestId).
				Int64("size", int64(c.Writer.Size())).
				Dur("duration", time.Since(startTime)).
				Strs("errors", c.Errors.Errors()).
				Msg("complete")
		} else if log.Info().Enabled() {
			c.Next()
			var event *zerolog.Event
			if len(c.Errors) > 0 {
				event = log.Error()
			} else {
				event = log.Info()
			}
			applyRequestEvent(event, c).Strs("errors", c.Errors.Errors()).Msg("complete")
		}
	}
}
