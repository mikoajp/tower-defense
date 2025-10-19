package ecs

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	gameconfig "tower-defense/internal/game/config"
)

// EntityFactory creates entities based on configuration
type EntityFactory struct {
	config *gameconfig.GameConfig
}

// NewEntityFactory creates a new entity factory
func NewEntityFactory(config *gameconfig.GameConfig) *EntityFactory {
	return &EntityFactory{config: config}
}

// CreateTower creates a new tower entity
func (f *EntityFactory) CreateTower(towerType string, pos Position) (*TowerEntity, error) {
	cfg, err := f.config.GetTowerConfig(towerType)
	if err != nil {
		return nil, err
	}
	
	tower := &TowerEntity{
		BaseEntity: BaseEntity{
			ID:       uuid.New().String(),
			Type:     EntityTypeTower,
			Position: pos,
			Alive:    true,
		},
		TowerType:    towerType,
		Range:        cfg.Range,
		Damage:       cfg.Damage,
		FireRate:     cfg.FireRate,
		SplashRadius: cfg.SplashRadius,
		LastShot:     time.Now().Add(-time.Hour), // Can shoot immediately
	}
	
	return tower, nil
}

// CreateEnemy creates a new enemy entity
func (f *EntityFactory) CreateEnemy(enemyType string, pos Position, wave int) (*EnemyEntity, error) {
	cfg, err := f.config.GetEnemyConfig(enemyType)
	if err != nil {
		return nil, err
	}
	
	// Scale HP based on wave
	hp := f.config.ScaleEnemyHP(cfg.HP, wave)
	
	enemy := &EnemyEntity{
		BaseEntity: BaseEntity{
			ID:       uuid.New().String(),
			Type:     EntityTypeEnemy,
			Position: pos,
			Alive:    true,
		},
		EnemyType:   enemyType,
		HP:          hp,
		MaxHP:       hp,
		Speed:       cfg.Speed,
		PathIndex:   0,
		GoldReward:  cfg.GoldReward,
		ScoreReward: cfg.ScoreReward,
	}
	
	return enemy, nil
}

// CreateProjectile creates a new projectile entity
func (f *EntityFactory) CreateProjectile(projType string, pos Position, targetID string, damage int, splashRadius float64) (*ProjectileEntity, error) {
	cfg, err := f.config.GetProjectileConfig(projType)
	if err != nil {
		return nil, err
	}
	
	projectile := &ProjectileEntity{
		BaseEntity: BaseEntity{
			ID:       uuid.New().String(),
			Type:     EntityTypeProjectile,
			Position: pos,
			Alive:    true,
		},
		ProjectileType: projType,
		Target:         targetID,
		Speed:          cfg.Speed,
		Damage:         damage,
		SplashRadius:   splashRadius,
	}
	
	return projectile, nil
}

// CreateEnemiesForWave creates all enemies for a given wave
func (f *EntityFactory) CreateEnemiesForWave(wave int, startPos Position) ([]*EnemyEntity, error) {
	// Calculate total number of enemies for this wave
	totalEnemies := f.config.CalculateEnemiesForWave(wave)
	composition := f.config.GetWaveComposition(wave)
	enemies := []*EnemyEntity{}
	
	// Calculate total weight from composition percentages
	totalWeight := composition.Basic + composition.Fast + composition.Tank + composition.Boss
	if totalWeight == 0 {
		totalWeight = 100 // Default if not specified
		composition.Basic = 100
	}
	
	// Calculate actual count for each enemy type based on percentages
	basicCount := (totalEnemies * composition.Basic) / totalWeight
	fastCount := (totalEnemies * composition.Fast) / totalWeight
	tankCount := (totalEnemies * composition.Tank) / totalWeight
	bossCount := (totalEnemies * composition.Boss) / totalWeight
	
	// Ensure at least totalEnemies are created (handle rounding)
	currentTotal := basicCount + fastCount + tankCount + bossCount
	if currentTotal < totalEnemies {
		// Add remaining to the most common type
		if composition.Basic > 0 {
			basicCount += totalEnemies - currentTotal
		} else if composition.Fast > 0 {
			fastCount += totalEnemies - currentTotal
		} else if composition.Tank > 0 {
			tankCount += totalEnemies - currentTotal
		} else {
			bossCount += totalEnemies - currentTotal
		}
	}
	
	// Helper to create N enemies of a type
	createN := func(enemyType string, count int) error {
		for i := 0; i < count; i++ {
			enemy, err := f.CreateEnemy(enemyType, startPos, wave)
			if err != nil {
				return err
			}
			enemies = append(enemies, enemy)
		}
		return nil
	}
	
	// Create enemies based on calculated counts
	if basicCount > 0 {
		if err := createN("basic", basicCount); err != nil {
			return nil, err
		}
	}
	if fastCount > 0 {
		if err := createN("fast", fastCount); err != nil {
			return nil, err
		}
	}
	if tankCount > 0 {
		if err := createN("tank", tankCount); err != nil {
			return nil, err
		}
	}
	if bossCount > 0 {
		if err := createN("boss", bossCount); err != nil {
			return nil, err
		}
	}
	
	if len(enemies) == 0 {
		return nil, fmt.Errorf("no enemies created for wave %d", wave)
	}
	
	return enemies, nil
}
