package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tower-defense/internal/config"
	"tower-defense/internal/game"
	gameconfig "tower-defense/internal/game/config"
	"tower-defense/internal/logging"
	"tower-defense/internal/server"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func main() {
	cfg := config.FromEnv()
	// init structured logging
	_ = logging.Init(cfg.LogLevel)
	defer logging.Sync()

	// Load game configuration
	gameCfg, err := gameconfig.Load()
	if err != nil {
		logging.Errorw("failed_to_load_game_config", "error", err)
		panic(err)
	}

	// Initialize game manager (supports multi-room)
	gameManager := game.NewManager(gameCfg)
	defer gameManager.Shutdown()

	// Get or create default game
	defaultGame := gameManager.GetOrCreateDefault()
	defaultGame.Start()

	// Prepare websocket upgrader with origin check
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			for _, ao := range cfg.AllowedOrigins {
				if ao == "*" || origin == ao { return true }
			}
			return false
		},
	}

	// Handlers
	// WebSocket hub setup
	hub := server.NewHub()
	go hub.Run()

	// Broadcaster: encode state once and distribute to clients
	go func() {
		// base interval 100ms, adaptive: skip if previous broadcast is recent
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		var last time.Time
		for range ticker.C {
			if time.Since(last) < 50*time.Millisecond { // simple adaptive throttling
				continue
			}
			b, err := defaultGame.MarshalState()
			if err != nil {
				continue
			}
			hub.Broadcast(b)
			last = time.Now()
		}
	}()

	wsHandler := gin.HandlerFunc(func(c *gin.Context) {
		server.WsConnections.Inc()
		defer server.WsConnections.Dec()
		serverHandler := hub.ServeWS(upgrader)
		serverHandler(c.Writer, c.Request)
		return // no JSON write here
		})

	addTower := func(c *gin.Context) {
		var req struct {
			X         float64 `json:"x"`
			Y         float64 `json:"y"`
			TowerType string  `json:"towerType"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		// Default to basic tower if not specified
		towerType := req.TowerType
		if towerType == "" {
			towerType = "basic"
		}
		
		if err := defaultGame.AddTower(towerType, req.X, req.Y); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"success": true})
		}
	}

	getState := func(c *gin.Context) {
		c.JSON(http.StatusOK, defaultGame.GetState())
	}

	reset := func(c *gin.Context) {
		logging.Infow("game_reset")
		defaultGame.Reset()
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Game reset successfully"})
	}
	
	// Multi-room handlers
	createGame := func(c *gin.Context) {
		game, err := gameManager.CreateGame()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		game.Start()
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"game_id": game.GetID(),
			"message": "Game created",
		})
	}
	
	listGames := func(c *gin.Context) {
		stats := gameManager.GetStats()
		c.JSON(http.StatusOK, stats)
	}
	
	// Save/Load handlers
	saveGame := func(c *gin.Context) {
		data, err := defaultGame.SaveState()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		// For now, just return the data as base64
		// In production, you'd save to repository
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Game saved",
			"size": len(data),
		})
	}
	
	loadGame := func(c *gin.Context) {
		// Accept raw JSON state
		var stateData []byte
		var err error
		
		// Try to read raw body
		stateData, err = c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		if err := defaultGame.LoadFromState(stateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Game loaded"})
	}

	// wire Prometheus metrics via on-tick hook
	defaultGame.SetOnTick(func(st game.TickStats) {
		server.TicksTotal.Inc()
		server.EngineEnemies.Set(float64(st.Enemies))
		server.EngineProjectiles.Set(float64(st.Projectiles))
		server.EngineTowers.Set(float64(st.Towers))
		server.EngineTickSeconds.Observe(st.Dt)
	})

	r := server.NewRouter(wsHandler, addTower, getState, reset, saveGame, loadGame, createGame, listGames, cfg.AllowedOrigins)
	// plug request logger is already in router; nothing else needed here
	// optional debug pprof
	server.MountPprof(r, cfg.EnablePprof)

	httpSrv := &http.Server{
		Addr:    cfg.Port,
		Handler: r,
	}

	// graceful shutdown
	go func() {
		logging.Infow("server_start", "port", cfg.Port)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Errorw("server_error", "error", err)
		}
	}()

	// Wait for termination signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logging.Infow("server_shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpSrv.Shutdown(ctx); err != nil {
		logging.Errorw("server_shutdown_error", "error", err)
	}
	defaultGame.Stop()
	logging.Infow("server_stopped")
}
