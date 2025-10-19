# Game Package - ECS Architecture

This package implements a clean Entity-Component-System (ECS) architecture for the tower defense game.

## Structure

```
game/
├── game.go              # Main game instance
├── manager.go           # Multi-room game manager
├── adapter.go           # Backward compatibility
├── state.go             # State DTOs
├── ecs/                 # Entity-Component-System
│   ├── entity.go        # Entity interfaces & types
│   ├── world.go         # Entity container & queries
│   └── factory.go       # Entity creation (data-driven)
├── systems/             # Game logic systems
│   ├── system.go        # System interface
│   ├── combat.go        # Tower shooting
│   ├── projectile.go    # Projectile movement
│   ├── movement.go      # Enemy movement
│   ├── wave.go          # Wave spawning
│   ├── reward.go        # Gold/score rewards
│   └── lifecycle.go     # Entity cleanup
├── config/              # Configuration
│   ├── loader.go        # YAML config loader
│   └── balance.yaml     # Game balance values
└── repository/          # Persistence
    ├── repository.go    # Repository interface
    ├── memory.go        # In-memory implementation
    └── file.go          # File-based implementation
```

## Quick Start

### Creating a Game

```go
import (
    "tower-defense/internal/game"
    "tower-defense/internal/game/config"
)

// Load configuration
cfg, _ := config.Load()

// Create game manager
manager := game.NewManager(cfg)

// Create a game
gameInstance := manager.GetOrCreateDefault()
gameInstance.Start()

// Place a tower
err := gameInstance.AddTower("basic", 100, 100)

// Get state
state := gameInstance.GetState()

// Stop game
gameInstance.Stop()
```

### Using Multiple Game Rooms

```go
// Create multiple games
game1, _ := manager.CreateGame()
game2, _ := manager.CreateGame()

game1.Start()
game2.Start()

// Each game runs independently
game1.AddTower("basic", 100, 100)
game2.AddTower("sniper", 200, 200)

// List all games
gameIDs := manager.ListGames()

// Remove a game
manager.RemoveGame(game1.GetID())
```

### Persistence

```go
import "tower-defense/internal/game/repository"

// Create repository
repo := repository.NewMemoryRepository()

// Save game state
stateJSON, _ := gameInstance.MarshalState()
saveID, _ := repo.Save(gameInstance.GetID(), stateJSON)

// Load game state
save, _ := repo.Load(saveID)
// Use save.Data to restore state
```

## Core Concepts

### Entities

Entities are game objects (towers, enemies, projectiles). They implement the `Entity` interface:

```go
type Entity interface {
    GetID() string
    GetType() EntityType
    GetPosition() Position
    SetPosition(Position)
    Update(dt float64)
    IsAlive() bool
}
```

**Entity Types:**
- `TowerEntity` - Defense towers that shoot at enemies
- `EnemyEntity` - Enemies that follow the path
- `ProjectileEntity` - Projectiles shot by towers

### World

The World manages all entities and provides efficient queries:

```go
world := ecs.NewWorld()

// Add entities
world.AddEntity(tower)
world.AddEntity(enemy)

// Query by type
towers := world.GetTowers()
enemies := world.GetEnemies()
projectiles := world.GetProjectiles()

// Get specific entity
enemy, exists := world.GetEnemy(enemyID)

// Cleanup dead entities
removed := world.CleanupDeadEntities()
```

### Systems

Systems contain game logic and operate on entities:

```go
type System interface {
    Update(world *ecs.World, dt float64)
}
```

**Available Systems:**

1. **WaveSystem** - Spawns waves of enemies
   - Reads composition from config
   - Spawns enemies with delay
   - Scales difficulty per wave

2. **MovementSystem** - Moves enemies along path
   - Uses path from config
   - Handles waypoint progression
   - Marks enemies that reached end

3. **CombatSystem** - Tower shooting logic
   - Finds targets in range
   - Respects fire rate cooldown
   - Creates projectiles

4. **ProjectileSystem** - Projectile behavior
   - Moves projectiles toward targets
   - Applies damage on hit
   - Removes dead projectiles

5. **RewardSystem** - Grants rewards
   - Gives gold when enemies die
   - Grants score points
   - Callbacks for state updates

6. **LifecycleSystem** - Entity cleanup
   - Removes dead entities
   - Handles life loss
   - Checks game over condition

### Entity Factory

The factory creates entities from configuration:

```go
factory := ecs.NewEntityFactory(config)

// Create tower
tower, _ := factory.CreateTower("basic", position)
tower, _ := factory.CreateTower("sniper", position)

// Create enemy
enemy, _ := factory.CreateEnemy("basic", position, wave)
enemy, _ := factory.CreateEnemy("tank", position, wave)

// Create projectile
proj, _ := factory.CreateProjectile("basic", position, targetID, damage)

// Create full wave
enemies, _ := factory.CreateEnemiesForWave(waveNumber, startPos)
```

### Configuration

All game balance comes from `balance.yaml`:

```yaml
game:
  starting_gold: 100
  starting_lives: 20

towers:
  basic:
    cost: 50
    range: 100.0
    damage: 10
    fire_rate: 1.0

enemies:
  basic:
    hp: 50
    speed: 1.0
    gold_reward: 10
    score_reward: 10

waves:
  early_waves:
    basic: 100
  mid_waves:
    basic: 70
    fast: 30
```

## System Update Order

Systems run in this order each tick:

1. WaveSystem - Spawn new enemies
2. MovementSystem - Move enemies
3. CombatSystem - Towers shoot
4. ProjectileSystem - Move projectiles
5. RewardSystem - Grant rewards
6. LifecycleSystem - Cleanup & life loss

This order ensures:
- Enemies spawn before movement
- Towers shoot at current positions
- Projectiles hit before cleanup
- Rewards given before entities removed
- Cleanup happens last

## Thread Safety

- Each `Game` instance has its own mutex
- `Manager` has a separate mutex for game map
- All public methods are thread-safe
- Multiple games can run concurrently

## Performance Considerations

### Efficient Queries
```go
// O(1) access to entities by type
towers := world.GetTowers()      // Fast
enemies := world.GetEnemies()    // Fast
projectiles := world.GetProjectiles() // Fast
```

### Batch Cleanup
```go
// Cleanup all dead entities at once
removed := world.CleanupDeadEntities()
```

### Lock Granularity
```go
// Fine-grained locking
g.mu.RLock()
towers := g.world.GetTowers()  // Read lock
g.mu.RUnlock()

g.mu.Lock()
g.world.AddEntity(tower)  // Write lock
g.mu.Unlock()
```

## Adding New Features

### New Tower Type

1. Add to `balance.yaml`:
```yaml
towers:
  freeze:
    cost: 100
    range: 80.0
    damage: 5
    fire_rate: 0.5
```

2. Optionally extend entity:
```go
type FreezeTowerEntity struct {
    TowerEntity
    FreezeDuration float64
}
```

3. Update CombatSystem if needed

### New System

1. Create system file:
```go
type BuffSystem struct {
    buffs map[string]Buff
}

func (s *BuffSystem) Update(world *ecs.World, dt float64) {
    // Apply buffs to entities
}
```

2. Register in game:
```go
buffSystem := systems.NewBuffSystem()
systemManager.AddSystem(buffSystem)
```

### New Enemy Type

1. Add to `balance.yaml`:
```yaml
enemies:
  flying:
    hp: 60
    speed: 1.8
    gold_reward: 15
    score_reward: 15
```

2. Update wave composition:
```yaml
waves:
  late_waves:
    basic: 40
    fast: 30
    tank: 20
    flying: 10
```

## Testing

```go
func TestGame(t *testing.T) {
    cfg := config.MustLoad()
    game := game.NewGame("test", cfg)
    
    game.Start()
    defer game.Stop()
    
    // Place tower
    err := game.AddTower("basic", 100, 100)
    assert.NoError(t, err)
    
    // Check state
    state := game.GetState()
    assert.Equal(t, 1, len(state.Towers))
}
```

## Backward Compatibility

The `EngineAdapter` provides compatibility with old code:

```go
// Old interface
eng := game.NewEngineAdapter(gameInstance)
eng.Start()
eng.AddTower(x, y)  // Assumes "basic" type
oldState := eng.GetState()  // Returns domain.GameState
```

This allows gradual migration from the old engine.

## Best Practices

1. **Always use the factory** to create entities
2. **Query World efficiently** using type-specific methods
3. **Don't hold locks long** - get data and release
4. **Systems should be stateless** when possible
5. **Configuration over code** - use YAML for balance
6. **Test systems in isolation** - mock World if needed

## Future Roadmap

- [ ] Entity pooling for performance
- [ ] Spatial partitioning for large maps
- [ ] Replay system using repository
- [ ] Network synchronization for multiplayer
- [ ] Save/load from database
- [ ] Tower upgrade system
- [ ] Special abilities system
- [ ] Achievement system
