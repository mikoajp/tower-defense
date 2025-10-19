package server

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

// NewRouter wires up the HTTP routes.
func NewRouter(wsHandler gin.HandlerFunc, addTower gin.HandlerFunc, getState gin.HandlerFunc, reset gin.HandlerFunc, saveGame gin.HandlerFunc, loadGame gin.HandlerFunc, createGame gin.HandlerFunc, listGames gin.HandlerFunc, listMaps gin.HandlerFunc, changeMap gin.HandlerFunc, allowedOrigins []string) *gin.Engine {
	r := gin.New()
	// logging + recovery
	r.Use(RequestLogger(), gin.Recovery())

	// CORS
	r.Use(CORS(allowedOrigins))

	// Versioned API group
	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
		v1.GET("/state", getState)
		v1.POST("/tower", addTower)
		v1.POST("/reset", reset)
		v1.POST("/save", saveGame)
		v1.POST("/load", loadGame)
		v1.POST("/games", createGame)
		v1.GET("/games", listGames)
		v1.GET("/maps", listMaps)
		v1.POST("/map", changeMap)
	}

	// Legacy routes (backward compatibility)
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	r.GET("/state", getState)
	r.POST("/tower", addTower)
	r.POST("/reset", reset)
	r.POST("/save", saveGame)
	r.POST("/load", loadGame)
	r.GET("/maps", listMaps)
	r.POST("/map", changeMap)

	// websocket (keep legacy path)
	r.GET("/ws", wsHandler)

	// metrics mount (optional)
	MountMetrics(r)

	return r
}
