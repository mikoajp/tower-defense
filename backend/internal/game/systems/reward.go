package systems

import (
	"tower-defense/internal/game/ecs"
	"tower-defense/internal/logging"
)

// RewardSystem handles giving gold and score when enemies die
type RewardSystem struct {
	onReward func(gold, score int)
}

// NewRewardSystem creates a new reward system
func NewRewardSystem(onReward func(gold, score int)) *RewardSystem {
	return &RewardSystem{
		onReward: onReward,
	}
}

// Update processes dead enemies and grants rewards
func (s *RewardSystem) Update(world *ecs.World, dt float64) {
	enemies := world.GetEnemies()

	for _, enemy := range enemies {
		// Check if enemy just died (HP <= 0 but still marked alive)
		if enemy.HP <= 0 && enemy.Alive {
			// Grant rewards
			if s.onReward != nil {
				s.onReward(enemy.GoldReward, enemy.ScoreReward)
				logging.Debugw("enemy_killed", 
					"enemy_id", enemy.ID, 
					"gold", enemy.GoldReward, 
					"score", enemy.ScoreReward)
			}
			// Mark as dead after granting rewards
			enemy.Alive = false
		}
	}
}
