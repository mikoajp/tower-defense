package game

import (
	"encoding/json"
	"sync"
	"time"

	"tower-defense/internal/game/config"
	"tower-defense/internal/game/ecs"
	"tower-defense/internal/game/systems"
	"tower-defense/internal/logging"
)

// GameState represents the current state of a game
type GameState struct {
	Wave     int  `json:"wave"`
	Gold     int  `json:"gold"`
	Lives    int  `json:"lives"`
	Score    int  `json:"score"`
	GameOver bool `json:"gameOver"`
}

// Game represents a single game instance using ECS architecture
type Game struct {
	mu              sync.RWMutex
	id              string
	config          *config.GameConfig
	world           *ecs.World
	factory         *ecs.EntityFactory
	systemManager   *systems.SystemManager
	state           GameState
	running         bool
	ticker          *time.Ticker
	lastUpdate      time.Time
	
	// Systems
	movementSystem  *systems.MovementSystem
	combatSystem    *systems.CombatSystem
	projectileSystem *systems.ProjectileSystem
	waveSystem      *systems.WaveSystem
	rewardSystem    *systems.RewardSystem
	lifecycleSystem *systems.LifecycleSystem
	
	// Callbacks
	onTick          func(TickStats)
}

// TickStats contains statistics about the current tick
type TickStats struct {
	Enemies     int
	Projectiles int
	Towers      int
	Dt          float64
}

// NewGame creates a new game instance
func NewGame(id string, cfg *config.GameConfig) *Game {
	world := ecs.NewWorld()
	factory := ecs.NewEntityFactory(cfg)
	systemManager := systems.NewSystemManager()
	
	// Get start position from config
	startPos := ecs.Position{X: 0, Y: 250}
	if len(cfg.Map.Path) > 0 {
		startPos = ecs.Position{
			X: cfg.Map.Path[0].X,
			Y: cfg.Map.Path[0].Y,
		}
	}
	
	game := &Game{
		id:            id,
		config:        cfg,
		world:         world,
		factory:       factory,
		systemManager: systemManager,
		state: GameState{
			Wave:     0,
			Gold:     cfg.Game.StartingGold,
			Lives:    cfg.Game.StartingLives,
			Score:    0,
			GameOver: false,
		},
		lastUpdate: time.Now(),
	}
	
	// Initialize systems
	game.movementSystem = systems.NewMovementSystem(cfg)
	game.combatSystem = systems.NewCombatSystem(cfg, factory)
	game.projectileSystem = systems.NewProjectileSystem()
	game.waveSystem = systems.NewWaveSystem(cfg, factory, startPos)
	
	game.rewardSystem = systems.NewRewardSystem(func(gold, score int) {
		// Note: This callback is called from Update() which already holds the lock
		// So we don't lock again to avoid deadlock
		game.state.Gold += gold
		game.state.Score += score
	})
	
	game.lifecycleSystem = systems.NewLifecycleSystem(len(cfg.Map.Path), func(lives int) {
		// Note: This callback is called from Update() which already holds the lock
		// So we don't lock again to avoid deadlock
		game.state.Lives -= lives
		if game.state.Lives <= 0 {
			game.state.GameOver = true
		}
	})
	
	// Register systems in order
	systemManager.AddSystem(game.waveSystem)
	systemManager.AddSystem(game.movementSystem)
	systemManager.AddSystem(game.combatSystem)
	systemManager.AddSystem(game.projectileSystem)
	systemManager.AddSystem(game.rewardSystem)
	systemManager.AddSystem(game.lifecycleSystem)
	
	return game
}

// Start starts the game loop
func (g *Game) Start() {
	g.mu.Lock()
	if g.running {
		g.mu.Unlock()
		return
	}
	g.running = true
	g.lastUpdate = time.Now()
	g.mu.Unlock()
	
	tickRate := time.Duration(g.config.Game.TickRateMs) * time.Millisecond
	g.ticker = time.NewTicker(tickRate)
	
	go func() {
		for range g.ticker.C {
			g.mu.RLock()
			running := g.running
			g.mu.RUnlock()
			
			if !running {
				return
			}
			g.Update()
		}
	}()
	
	logging.Infow("game_started", "game_id", g.id)
}

// Stop stops the game loop
func (g *Game) Stop() {
	g.mu.Lock()
	defer g.mu.Unlock()
	
	if !g.running {
		return
	}
	
	g.running = false
	if g.ticker != nil {
		g.ticker.Stop()
	}
	
	logging.Infow("game_stopped", "game_id", g.id)
}

// Update processes one game tick
func (g *Game) Update() {
	g.mu.Lock()
	defer g.mu.Unlock()
	
	if g.state.GameOver {
		return
	}
	
	now := time.Now()
	dt := now.Sub(g.lastUpdate).Seconds()
	
	// Clamp dt to prevent large jumps
	if dt < 0 {
		dt = 0
	}
	if dt > 0.05 {
		dt = 0.05
	}
	
	g.lastUpdate = now
	
	// Update wave number from wave system
	g.state.Wave = g.waveSystem.GetCurrentWave()
	
	// Run all systems
	g.systemManager.Update(g.world, dt)
	
	// Send tick stats
	if g.onTick != nil {
		stats := TickStats{
			Enemies:     len(g.world.GetEnemies()),
			Projectiles: len(g.world.GetProjectiles()),
			Towers:      len(g.world.GetTowers()),
			Dt:          dt,
		}
		g.onTick(stats)
	}
}

// AddTower attempts to place a tower at the given position
func (g *Game) AddTower(towerType string, x, y float64) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	
	// Get tower config
	towerCfg, err := g.config.GetTowerConfig(towerType)
	if err != nil {
		return err
	}
	
	// Check if player has enough gold
	if g.state.Gold < towerCfg.Cost {
		return ErrNotEnoughGold
	}
	
	// Check tower placement rules
	pos := ecs.Position{X: x, Y: y}
	if !g.isValidPlacement(pos) {
		return ErrInvalidPlacement
	}
	
	// Create and place tower
	tower, err := g.factory.CreateTower(towerType, pos)
	if err != nil {
		return err
	}
	
	g.world.AddEntity(tower)
	g.state.Gold -= towerCfg.Cost
	
	logging.Infow("tower_placed", 
		"game_id", g.id, 
		"tower_type", towerType,
		"x", x, "y", y, 
		"gold_remaining", g.state.Gold)
	
	return nil
}

// isValidPlacement checks if a tower can be placed at the given position
func (g *Game) isValidPlacement(pos ecs.Position) bool {
	// Check tower count limit
	if g.world.TowerCount() >= g.config.Placement.MaxTowers {
		return false
	}
	
	// Check distance from path
	path := g.movementSystem.GetPath()
	minDistFromPath := g.config.Placement.MinDistanceFromPath
	
	for i := 0; i < len(path)-1; i++ {
		p1 := path[i]
		p2 := path[i+1]
		
		dist := distanceToSegment(pos, p1, p2)
		if dist < minDistFromPath {
			return false
		}
	}
	
	// Check distance from other towers
	towers := g.world.GetTowers()
	minSpacing := g.config.Placement.MinTowerSpacing
	
	for _, tower := range towers {
		dx := pos.X - tower.Position.X
		dy := pos.Y - tower.Position.Y
		dist := dx*dx + dy*dy
		
		if dist < minSpacing*minSpacing {
			return false
		}
	}
	
	return true
}

// GetState returns the current game state (thread-safe)
func (g *Game) GetState() GameStateSnapshot {
	g.mu.RLock()
	defer g.mu.RUnlock()
	
	return GameStateSnapshot{
		Towers:      g.convertTowers(),
		Enemies:     g.convertEnemies(),
		Projectiles: g.convertProjectiles(),
		Wave:        g.state.Wave,
		Gold:        g.state.Gold,
		Lives:       g.state.Lives,
		Score:       g.state.Score,
		GameOver:    g.state.GameOver,
	}
}

// MarshalState returns the game state as JSON
func (g *Game) MarshalState() ([]byte, error) {
	state := g.GetState()
	return json.Marshal(state)
}

// Reset resets the game to initial state
func (g *Game) Reset() {
	g.mu.Lock()
	defer g.mu.Unlock()
	
	// Clear world
	g.world.Clear()
	
	// Reset state
	g.state = GameState{
		Wave:     0,
		Gold:     g.config.Game.StartingGold,
		Lives:    g.config.Game.StartingLives,
		Score:    0,
		GameOver: false,
	}
	
	// Reset wave system
	g.waveSystem.Reset()
	
	logging.Infow("game_reset", "game_id", g.id)
}

// SetOnTick sets the tick callback
func (g *Game) SetOnTick(f func(TickStats)) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.onTick = f
}

// GetID returns the game ID
func (g *Game) GetID() string {
	return g.id
}

// SaveState saves the current game state and returns the serialized data
func (g *Game) SaveState() ([]byte, error) {
	return g.MarshalState()
}

// LoadFromState loads game state from serialized data
func (g *Game) LoadFromState(data []byte) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	
	var snapshot GameStateSnapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return err
	}
	
	// Clear current world
	g.world.Clear()
	
	// Restore basic state
	g.state.Wave = snapshot.Wave
	g.state.Gold = snapshot.Gold
	g.state.Lives = snapshot.Lives
	g.state.Score = snapshot.Score
	g.state.GameOver = snapshot.GameOver
	
	// Restore towers
	for _, towerDTO := range snapshot.Towers {
		tower := &ecs.TowerEntity{
			BaseEntity: ecs.BaseEntity{
				ID:       towerDTO.ID,
				Type:     ecs.EntityTypeTower,
				Position: ecs.Position{X: towerDTO.Position.X, Y: towerDTO.Position.Y},
				Alive:    true,
			},
			TowerType:    towerDTO.Type,
			Range:        towerDTO.Range,
			Damage:       towerDTO.Damage,
			FireRate:     towerDTO.FireRate,
			SplashRadius: towerDTO.SplashRadius,
			LastShot:     time.Now(),
		}
		g.world.AddEntity(tower)
	}
	
	// Restore enemies
	for _, enemyDTO := range snapshot.Enemies {
		enemy := &ecs.EnemyEntity{
			BaseEntity: ecs.BaseEntity{
				ID:       enemyDTO.ID,
				Type:     ecs.EntityTypeEnemy,
				Position: ecs.Position{X: enemyDTO.Position.X, Y: enemyDTO.Position.Y},
				Alive:    true,
			},
			EnemyType: enemyDTO.Type,
			HP:        enemyDTO.HP,
			MaxHP:     enemyDTO.MaxHP,
			Speed:     enemyDTO.Speed,
			PathIndex: enemyDTO.PathIndex,
		}
		g.world.AddEntity(enemy)
	}
	
	// Restore projectiles
	for _, projDTO := range snapshot.Projectiles {
		projectile := &ecs.ProjectileEntity{
			BaseEntity: ecs.BaseEntity{
				ID:       projDTO.ID,
				Type:     ecs.EntityTypeProjectile,
				Position: ecs.Position{X: projDTO.Position.X, Y: projDTO.Position.Y},
				Alive:    true,
			},
			ProjectileType: projDTO.Type,
			Target:         projDTO.Target,
			Speed:          projDTO.Speed,
			Damage:         projDTO.Damage,
			SplashRadius:   projDTO.SplashRadius,
		}
		g.world.AddEntity(projectile)
	}
	
	// Update wave system
	g.waveSystem.SetCurrentWave(snapshot.Wave)
	
	logging.Infow("game_loaded", "game_id", g.id, "wave", snapshot.Wave, "gold", snapshot.Gold)
	
	return nil
}

// Helper function to calculate distance from point to line segment
func distanceToSegment(p, a, b ecs.Position) float64 {
	dx := b.X - a.X
	dy := b.Y - a.Y
	
	if dx == 0 && dy == 0 {
		// Segment is a point
		px := p.X - a.X
		py := p.Y - a.Y
		return px*px + py*py
	}
	
	t := ((p.X-a.X)*dx + (p.Y-a.Y)*dy) / (dx*dx + dy*dy)
	
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}
	
	nearestX := a.X + t*dx
	nearestY := a.Y + t*dy
	
	px := p.X - nearestX
	py := p.Y - nearestY
	
	return px*px + py*py
}
