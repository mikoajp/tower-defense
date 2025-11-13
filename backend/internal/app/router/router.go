package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"tower-defense/internal/app/middleware"
	"tower-defense/internal/infrastructure/metrics"
)

// NewRouter wires up the HTTP routes.
func NewRouter(wsHandler gin.HandlerFunc, addTower gin.HandlerFunc, getState gin.HandlerFunc, reset gin.HandlerFunc, saveGame gin.HandlerFunc, loadGame gin.HandlerFunc, createGame gin.HandlerFunc, listGames gin.HandlerFunc, listMaps gin.HandlerFunc, changeMap gin.HandlerFunc, allowedOrigins []string) *gin.Engine {
	r := gin.New()
	// logging + recovery
	r.Use(middleware.RequestLogger(), gin.Recovery())

	// CORS
	r.Use(middleware.CORS(allowedOrigins))

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
	metrics.MountMetrics(r)

	return r
}
