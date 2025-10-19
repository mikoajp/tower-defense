package systems

import (
	"math"

	"tower-defense/internal/game/config"
	"tower-defense/internal/game/ecs"
)

// MovementSystem handles enemy movement along the path
type MovementSystem struct {
	config *config.GameConfig
	path   []ecs.Position
}

// NewMovementSystem creates a new movement system
func NewMovementSystem(cfg *config.GameConfig) *MovementSystem {
	// Convert config positions to ecs positions
	path := make([]ecs.Position, len(cfg.Map.Path))
	for i, p := range cfg.Map.Path {
		path[i] = ecs.Position{X: p.X, Y: p.Y}
	}
	
	return &MovementSystem{
		config: cfg,
		path:   path,
	}
}

// Update moves all enemies along the path
func (s *MovementSystem) Update(world *ecs.World, dt float64) {
	enemies := world.GetEnemies()
	
	for _, enemy := range enemies {
		if !enemy.Alive {
			continue
		}
		
		// Check if enemy is beyond the path (let LifecycleSystem handle this)
		if enemy.PathIndex >= len(s.path)-1 {
			// Don't set Alive = false here - let LifecycleSystem handle life loss
			continue
		}
		
		target := s.path[enemy.PathIndex+1]
		current := enemy.Position
		
		// Calculate direction
		dx := target.X - current.X
		dy := target.Y - current.Y
		distance := math.Sqrt(dx*dx + dy*dy)
		
		if distance < 1.0 {
			// Reached waypoint, move to next
			enemy.PathIndex++
			// Don't set Alive = false here - let LifecycleSystem handle it
			continue
		}
		
		// Move towards target
		moveDistance := enemy.Speed * dt * 60.0 // Normalize to 60 FPS
		if moveDistance > distance {
			moveDistance = distance
		}
		
		ratio := moveDistance / distance
		newPos := ecs.Position{
			X: current.X + dx*ratio,
			Y: current.Y + dy*ratio,
		}
		
		enemy.SetPosition(newPos)
	}
}

// GetPath returns the path for external use
func (s *MovementSystem) GetPath() []ecs.Position {
	return s.path
}
