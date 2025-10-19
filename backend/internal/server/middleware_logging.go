package server

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tower-defense/internal/logging"
)

// RequestLogger adds request-id and logs basic request/response info
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" { reqID = uuid.NewString() }
		c.Writer.Header().Set("X-Request-ID", reqID)
		start := time.Now()
		c.Next()
		dur := time.Since(start)
		logging.Infow("http_request",
			"req_id", reqID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration_ms", dur.Milliseconds(),
			"client_ip", c.ClientIP(),
		)
	}
}
