package systems

import (
	"math"

	"tower-defense/internal/game/ecs"
)

// ProjectileSystem handles projectile movement and collision
type ProjectileSystem struct{}

// NewProjectileSystem creates a new projectile system
func NewProjectileSystem() *ProjectileSystem {
	return &ProjectileSystem{}
}

// Update processes projectile movement and hits
func (s *ProjectileSystem) Update(world *ecs.World, dt float64) {
	projectiles := world.GetProjectiles()

	for _, proj := range projectiles {
		if !proj.Alive {
			continue
		}

		// Find target enemy
		target, exists := world.GetEnemy(proj.Target)
		if !exists || !target.Alive {
			// Target is dead or missing, remove projectile
			proj.Alive = false
			continue
		}

		// Calculate distance to target
		dx := target.Position.X - proj.Position.X
		dy := target.Position.Y - proj.Position.Y
		distance := math.Sqrt(dx*dx + dy*dy)

		// Move projectile
		moveDistance := proj.Speed * dt * 60.0 // Normalize to 60 FPS

		if distance <= moveDistance {
			// Hit target
			target.TakeDamage(proj.Damage)
			proj.Alive = false
			
			// Apply splash damage if projectile has splash radius
			if proj.SplashRadius > 0 {
				s.applySplashDamage(world, target.Position, proj.SplashRadius, proj.Damage, target.ID)
			}
		} else {
			// Move towards target
			ratio := moveDistance / distance
			newPos := ecs.Position{
				X: proj.Position.X + dx*ratio,
				Y: proj.Position.Y + dy*ratio,
			}
			proj.SetPosition(newPos)
		}
	}
}

// applySplashDamage applies area damage to enemies near the impact point
func (s *ProjectileSystem) applySplashDamage(world *ecs.World, impactPos ecs.Position, radius float64, damage int, primaryTargetID string) {
	enemies := world.GetEnemies()
	
	// Splash damage is 50% of primary damage
	splashDamage := damage / 2
	if splashDamage < 1 {
		splashDamage = 1
	}
	
	for _, enemy := range enemies {
		if !enemy.Alive || enemy.ID == primaryTargetID {
			continue
		}
		
		// Calculate distance from impact
		dx := enemy.Position.X - impactPos.X
		dy := enemy.Position.Y - impactPos.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		
		// Apply damage if within splash radius
		if dist <= radius {
			enemy.TakeDamage(splashDamage)
		}
	}
}
