package config

import (
	"embed"
	"fmt"

	"gopkg.in/yaml.v3"
)

//go:embed balance.yaml
var configFS embed.FS

// GameConfig represents the entire game configuration
type GameConfig struct {
	Game       GameSettings       `yaml:"game"`
	Towers     map[string]TowerConfig     `yaml:"towers"`
	Enemies    map[string]EnemyConfig     `yaml:"enemies"`
	Projectiles map[string]ProjectileConfig `yaml:"projectiles"`
	Waves      WaveConfig         `yaml:"waves"`
	Map        MapConfig          `yaml:"map"`
	Placement  PlacementConfig    `yaml:"placement"`
}

type GameSettings struct {
	StartingGold        int `yaml:"starting_gold"`
	StartingLives       int `yaml:"starting_lives"`
	TickRateMs          int `yaml:"tick_rate_ms"`
	BroadcastIntervalMs int `yaml:"broadcast_interval_ms"`
}

type TowerConfig struct {
	Cost         int     `yaml:"cost"`
	Range        float64 `yaml:"range"`
	Damage       int     `yaml:"damage"`
	FireRate     float64 `yaml:"fire_rate"`
	SplashRadius float64 `yaml:"splash_radius,omitempty"`
}

type EnemyConfig struct {
	HP          int     `yaml:"hp"`
	Speed       float64 `yaml:"speed"`
	GoldReward  int     `yaml:"gold_reward"`
	ScoreReward int     `yaml:"score_reward"`
}

type ProjectileConfig struct {
	Speed float64 `yaml:"speed"`
}

type WaveConfig struct {
	SpawnIntervalTicks       int     `yaml:"spawn_interval_ticks"`
	EnemiesPerWaveBase       int     `yaml:"enemies_per_wave_base"`
	EnemiesPerWaveMultiplier float64 `yaml:"enemies_per_wave_multiplier"`
	HPScalePerWave           float64 `yaml:"hp_scale_per_wave"`
	EarlyWaves               WaveComposition `yaml:"early_waves"`
	MidWaves                 WaveComposition `yaml:"mid_waves"`
	LateWaves                WaveComposition `yaml:"late_waves"`
	BossWaves                WaveComposition `yaml:"boss_waves"`
}

type WaveComposition struct {
	Basic int `yaml:"basic,omitempty"`
	Fast  int `yaml:"fast,omitempty"`
	Tank  int `yaml:"tank,omitempty"`
	Boss  int `yaml:"boss,omitempty"`
}

type MapConfig struct {
	Width         int            `yaml:"width"`
	Height        int            `yaml:"height"`
	Path          []Position     `yaml:"path"`
	PathHalfWidth float64        `yaml:"path_half_width"`
}

type Position struct {
	X float64 `yaml:"x"`
	Y float64 `yaml:"y"`
}

type PlacementConfig struct {
	MinDistanceFromPath float64 `yaml:"min_distance_from_path"`
	MinTowerSpacing     float64 `yaml:"min_tower_spacing"`
	MaxTowers           int     `yaml:"max_towers"`
}

// Global config instance
var Config *GameConfig

// Load reads and parses the balance configuration
func Load() (*GameConfig, error) {
	data, err := configFS.ReadFile("balance.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg GameConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	Config = &cfg
	return &cfg, nil
}

// MustLoad loads config or panics
func MustLoad() *GameConfig {
	cfg, err := Load()
	if err != nil {
		panic(err)
	}
	return cfg
}

// GetTowerConfig returns config for a tower type
func (c *GameConfig) GetTowerConfig(towerType string) (TowerConfig, error) {
	cfg, ok := c.Towers[towerType]
	if !ok {
		return TowerConfig{}, fmt.Errorf("unknown tower type: %s", towerType)
	}
	return cfg, nil
}

// GetEnemyConfig returns config for an enemy type
func (c *GameConfig) GetEnemyConfig(enemyType string) (EnemyConfig, error) {
	cfg, ok := c.Enemies[enemyType]
	if !ok {
		return EnemyConfig{}, fmt.Errorf("unknown enemy type: %s", enemyType)
	}
	return cfg, nil
}

// GetProjectileConfig returns config for a projectile type
func (c *GameConfig) GetProjectileConfig(projType string) (ProjectileConfig, error) {
	cfg, ok := c.Projectiles[projType]
	if !ok {
		return ProjectileConfig{}, fmt.Errorf("unknown projectile type: %s", projType)
	}
	return cfg, nil
}

// GetWaveComposition returns enemy composition for a given wave number
func (c *GameConfig) GetWaveComposition(wave int) WaveComposition {
	if wave%10 == 0 {
		return c.Waves.BossWaves
	} else if wave <= 5 {
		return c.Waves.EarlyWaves
	} else if wave <= 10 {
		return c.Waves.MidWaves
	}
	return c.Waves.LateWaves
}

// CalculateEnemiesForWave returns number of enemies to spawn for a wave
func (c *GameConfig) CalculateEnemiesForWave(wave int) int {
	base := float64(c.Waves.EnemiesPerWaveBase)
	multiplier := c.Waves.EnemiesPerWaveMultiplier
	count := base * (1 + (float64(wave-1) * (multiplier - 1)))
	return int(count)
}

// ScaleEnemyHP scales enemy HP based on wave number
func (c *GameConfig) ScaleEnemyHP(baseHP int, wave int) int {
	if wave <= 1 {
		return baseHP
	}
	scale := c.Waves.HPScalePerWave
	scaled := float64(baseHP) * (1 + (float64(wave-1) * (scale - 1)))
	return int(scaled)
}
