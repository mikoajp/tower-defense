package ecs

import (
	"sync"
)

// World manages all entities and provides queries
type World struct {
	mu       sync.RWMutex
	entities map[string]Entity
	
	// Indexed by type for fast queries
	towers      map[string]*TowerEntity
	enemies     map[string]*EnemyEntity
	projectiles map[string]*ProjectileEntity
}

// NewWorld creates a new ECS world
func NewWorld() *World {
	return &World{
		entities:    make(map[string]Entity),
		towers:      make(map[string]*TowerEntity),
		enemies:     make(map[string]*EnemyEntity),
		projectiles: make(map[string]*ProjectileEntity),
	}
}

// AddEntity adds an entity to the world
func (w *World) AddEntity(entity Entity) {
	w.mu.Lock()
	defer w.mu.Unlock()
	
	id := entity.GetID()
	w.entities[id] = entity
	
	// Add to type-specific index
	switch e := entity.(type) {
	case *TowerEntity:
		w.towers[id] = e
	case *EnemyEntity:
		w.enemies[id] = e
	case *ProjectileEntity:
		w.projectiles[id] = e
	}
}

// RemoveEntity removes an entity from the world
func (w *World) RemoveEntity(id string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	
	entity, ok := w.entities[id]
	if !ok {
		return
	}
	
	delete(w.entities, id)
	
	// Remove from type-specific index
	switch entity.GetType() {
	case EntityTypeTower:
		delete(w.towers, id)
	case EntityTypeEnemy:
		delete(w.enemies, id)
	case EntityTypeProjectile:
		delete(w.projectiles, id)
	}
}

// GetEntity retrieves an entity by ID
func (w *World) GetEntity(id string) (Entity, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	entity, ok := w.entities[id]
	return entity, ok
}

// GetTowers returns all tower entities
func (w *World) GetTowers() []*TowerEntity {
	w.mu.RLock()
	defer w.mu.RUnlock()
	
	towers := make([]*TowerEntity, 0, len(w.towers))
	for _, t := range w.towers {
		if t.Alive {
			towers = append(towers, t)
		}
	}
	return towers
}

// GetEnemies returns all enemy entities
func (w *World) GetEnemies() []*EnemyEntity {
	w.mu.RLock()
	defer w.mu.RUnlock()
	
	enemies := make([]*EnemyEntity, 0, len(w.enemies))
	for _, e := range w.enemies {
		if e.Alive {
			enemies = append(enemies, e)
		}
	}
	return enemies
}

// GetProjectiles returns all projectile entities
func (w *World) GetProjectiles() []*ProjectileEntity {
	w.mu.RLock()
	defer w.mu.RUnlock()
	
	projectiles := make([]*ProjectileEntity, 0, len(w.projectiles))
	for _, p := range w.projectiles {
		if p.Alive {
			projectiles = append(projectiles, p)
		}
	}
	return projectiles
}

// GetEnemy retrieves a specific enemy by ID
func (w *World) GetEnemy(id string) (*EnemyEntity, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	enemy, ok := w.enemies[id]
	return enemy, ok
}

// CleanupDeadEntities removes all dead entities
func (w *World) CleanupDeadEntities() []string {
	w.mu.Lock()
	defer w.mu.Unlock()
	
	var removed []string
	
	for id, entity := range w.entities {
		if !entity.IsAlive() {
			delete(w.entities, id)
			
			switch entity.GetType() {
			case EntityTypeTower:
				delete(w.towers, id)
			case EntityTypeEnemy:
				delete(w.enemies, id)
			case EntityTypeProjectile:
				delete(w.projectiles, id)
			}
			
			removed = append(removed, id)
		}
	}
	
	return removed
}

// Clear removes all entities
func (w *World) Clear() {
	w.mu.Lock()
	defer w.mu.Unlock()
	
	w.entities = make(map[string]Entity)
	w.towers = make(map[string]*TowerEntity)
	w.enemies = make(map[string]*EnemyEntity)
	w.projectiles = make(map[string]*ProjectileEntity)
}

// EntityCount returns the total number of entities
func (w *World) EntityCount() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return len(w.entities)
}

// TowerCount returns the number of towers
func (w *World) TowerCount() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	count := 0
	for _, t := range w.towers {
		if t.Alive {
			count++
		}
	}
	return count
}
