package systems

import (
	"tower-defense/internal/game/ecs"
)

// System is the interface for all game systems
type System interface {
	Update(world *ecs.World, dt float64)
}

// SystemManager manages and updates all systems
type SystemManager struct {
	systems []System
}

// NewSystemManager creates a new system manager
func NewSystemManager() *SystemManager {
	return &SystemManager{
		systems: make([]System, 0),
	}
}

// AddSystem adds a system to the manager
func (sm *SystemManager) AddSystem(system System) {
	sm.systems = append(sm.systems, system)
}

// Update updates all systems in order
func (sm *SystemManager) Update(world *ecs.World, dt float64) {
	for _, system := range sm.systems {
		system.Update(world, dt)
	}
}
