package systems

import (
	"math/rand"
	"time"

	"tower-defense/internal/game/config"
	"tower-defense/internal/game/ecs"
	"tower-defense/internal/logging"
)

// WaveSystem handles wave spawning and enemy creation
type WaveSystem struct {
	config          *config.GameConfig
	factory         *ecs.EntityFactory
	startPos        ecs.Position
	currentWave     int
	remainingInWave int
	nextEnemySpawn  time.Time
	lastWaveTime    time.Time
	waveInterval    time.Duration
	rng             *rand.Rand
}

// NewWaveSystem creates a new wave system
func NewWaveSystem(cfg *config.GameConfig, factory *ecs.EntityFactory, startPos ecs.Position) *WaveSystem {
	return &WaveSystem{
		config:       cfg,
		factory:      factory,
		startPos:     startPos,
		currentWave:  0,
		waveInterval: 10 * time.Second,
		lastWaveTime: time.Now(),
		rng:          rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Update processes wave spawning
func (s *WaveSystem) Update(world *ecs.World, dt float64) {
	now := time.Now()

	// Check if it's time to spawn a new wave
	if s.remainingInWave == 0 && now.Sub(s.lastWaveTime) > s.waveInterval {
		s.spawnWave(world)
		s.lastWaveTime = now
	}

	// Spawn enemies from current wave
	for s.remainingInWave > 0 && now.After(s.nextEnemySpawn) {
		s.spawnNextEnemy(world)
		s.nextEnemySpawn = now.Add(s.nextSpawnDelay())
	}
}

// spawnWave starts a new wave
func (s *WaveSystem) spawnWave(world *ecs.World) {
	s.currentWave++
	
	// Create enemies for this wave
	enemies, err := s.factory.CreateEnemiesForWave(s.currentWave, s.startPos)
	if err != nil {
		logging.Errorw("wave_spawn_error", "wave", s.currentWave, "error", err)
		return
	}

	s.remainingInWave = len(enemies)
	logging.Infow("wave_started", "wave", s.currentWave, "enemy_count", len(enemies))

	// Spawn first enemy immediately
	if len(enemies) > 0 {
		world.AddEntity(enemies[0])
		s.remainingInWave--
		// Store remaining enemies for later spawning
		for i := 1; i < len(enemies); i++ {
			// We'll spawn these over time
			s.remainingInWave++
		}
		s.nextEnemySpawn = time.Now().Add(s.nextSpawnDelay())
	}
}

// spawnNextEnemy spawns the next enemy in the current wave
func (s *WaveSystem) spawnNextEnemy(world *ecs.World) {
	if s.remainingInWave <= 0 {
		return
	}

	// Determine enemy type based on wave composition
	composition := s.config.GetWaveComposition(s.currentWave)
	enemyType := s.selectEnemyType(composition)

	enemy, err := s.factory.CreateEnemy(enemyType, s.startPos, s.currentWave)
	if err != nil {
		logging.Errorw("enemy_spawn_error", "type", enemyType, "error", err)
		return
	}

	world.AddEntity(enemy)
	s.remainingInWave--
}

// selectEnemyType selects a random enemy type based on wave composition
func (s *WaveSystem) selectEnemyType(comp config.WaveComposition) string {
	total := comp.Basic + comp.Fast + comp.Tank + comp.Boss
	if total == 0 {
		return "basic"
	}

	roll := s.rng.Intn(total)
	
	if roll < comp.Basic {
		return "basic"
	}
	roll -= comp.Basic
	
	if roll < comp.Fast {
		return "fast"
	}
	roll -= comp.Fast
	
	if roll < comp.Tank {
		return "tank"
	}
	
	return "boss"
}

// nextSpawnDelay returns a random delay for next enemy spawn
func (s *WaveSystem) nextSpawnDelay() time.Duration {
	baseDelay := 120
	variance := 181
	delay := baseDelay + s.rng.Intn(variance)
	return time.Duration(delay) * time.Millisecond
}

// GetCurrentWave returns the current wave number
func (s *WaveSystem) GetCurrentWave() int {
	return s.currentWave
}

// SetCurrentWave sets the current wave number (for loading saved games)
func (s *WaveSystem) SetCurrentWave(wave int) {
	s.currentWave = wave
}

// Reset resets the wave system
func (s *WaveSystem) Reset() {
	s.currentWave = 0
	s.remainingInWave = 0
	s.lastWaveTime = time.Now()
}
