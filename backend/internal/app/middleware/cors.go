package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS returns a gin middleware that sets CORS headers according to allowed origins.
// If list contains "*", all origins are allowed. Otherwise only exact matches are allowed.
func CORS(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// origin allow check
		allowedAll := false
		for _, ao := range allowedOrigins { if ao == "*" { allowedAll = true; break } }
		if allowedAll {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			for _, ao := range allowedOrigins {
				if origin == ao {
					c.Writer.Header().Set("Access-Control-Allow-Origin", ao)
					break
				}
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
