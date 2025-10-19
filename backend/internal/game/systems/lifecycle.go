package systems

import (
	"tower-defense/internal/game/ecs"
	"tower-defense/internal/logging"
)

// LifecycleSystem handles entity cleanup and life loss
type LifecycleSystem struct {
	onLifeLost func(lives int)
	pathLength int
}

// NewLifecycleSystem creates a new lifecycle system
func NewLifecycleSystem(pathLength int, onLifeLost func(lives int)) *LifecycleSystem {
	return &LifecycleSystem{
		onLifeLost: onLifeLost,
		pathLength: pathLength,
	}
}

// Update cleans up dead entities and handles enemies reaching the end
func (s *LifecycleSystem) Update(world *ecs.World, dt float64) {
	enemies := world.GetEnemies()

	// Check for enemies that reached the end
	for _, enemy := range enemies {
		if enemy.Alive && enemy.PathIndex >= s.pathLength-1 {
			enemy.Alive = false
			if s.onLifeLost != nil {
				s.onLifeLost(1)
				logging.Warnw("enemy_reached_end", "enemy_id", enemy.ID)
			}
		}
	}

	// Clean up dead entities
	removed := world.CleanupDeadEntities()
	if len(removed) > 0 {
		logging.Debugw("entities_cleaned", "count", len(removed))
	}
}
