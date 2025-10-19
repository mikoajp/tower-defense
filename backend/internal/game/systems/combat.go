package systems

import (
	"math"

	"tower-defense/internal/game/config"
	"tower-defense/internal/game/ecs"
)

// CombatSystem handles tower shooting and target acquisition
type CombatSystem struct {
	config  *config.GameConfig
	factory *ecs.EntityFactory
}

// NewCombatSystem creates a new combat system
func NewCombatSystem(cfg *config.GameConfig, factory *ecs.EntityFactory) *CombatSystem {
	return &CombatSystem{
		config:  cfg,
		factory: factory,
	}
}

// Update processes tower shooting logic
func (s *CombatSystem) Update(world *ecs.World, dt float64) {
	towers := world.GetTowers()
	enemies := world.GetEnemies()

	for _, tower := range towers {
		if !tower.Alive {
			continue
		}

		// Check if tower can shoot
		if !tower.CanShoot() {
			continue
		}

		// Find closest enemy in range
		var closestEnemy *ecs.EnemyEntity
		minDist := tower.Range

		for _, enemy := range enemies {
			if !enemy.Alive {
				continue
			}

			dx := enemy.Position.X - tower.Position.X
			dy := enemy.Position.Y - tower.Position.Y
			dist := math.Sqrt(dx*dx + dy*dy)

			if dist < minDist {
				minDist = dist
				closestEnemy = enemy
			}
		}

		// Shoot at closest enemy
		if closestEnemy != nil {
			// Determine projectile type based on tower type
			projType := "basic"
			if tower.TowerType == "sniper" {
				projType = "sniper"
			} else if tower.TowerType == "splash" {
				projType = "splash"
			}

			projectile, err := s.factory.CreateProjectile(
				projType,
				tower.Position,
				closestEnemy.ID,
				tower.Damage,
				tower.SplashRadius,
			)
			if err == nil {
				world.AddEntity(projectile)
				tower.Shoot()
			}
		}
	}
}
