package game

import (
	"fmt"
	"sync"

	"tower-defense/internal/game/config"
	"tower-defense/internal/logging"
	"github.com/google/uuid"
)

// Manager manages multiple game instances (multi-room support)
type Manager struct {
	mu     sync.RWMutex
	games  map[string]*Game
	config *config.GameConfig
}

// NewManager creates a new game manager
func NewManager(cfg *config.GameConfig) *Manager {
	return &Manager{
		games:  make(map[string]*Game),
		config: cfg,
	}
}

// CreateGame creates a new game instance with a unique ID
func (m *Manager) CreateGame() (*Game, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	gameID := uuid.New().String()
	game := NewGame(gameID, m.config)
	m.games[gameID] = game
	
	logging.Infow("game_created", "game_id", gameID, "total_games", len(m.games))
	
	return game, nil
}

// GetGame retrieves a game by ID
func (m *Manager) GetGame(gameID string) (*Game, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	game, exists := m.games[gameID]
	if !exists {
		return nil, ErrGameNotFound
	}
	
	return game, nil
}

// GetOrCreateDefault gets the default game or creates it if it doesn't exist
func (m *Manager) GetOrCreateDefault() *Game {
	defaultID := "default"
	
	m.mu.RLock()
	game, exists := m.games[defaultID]
	m.mu.RUnlock()
	
	if exists {
		return game
	}
	
	// Create default game
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Double-check after acquiring write lock
	if game, exists := m.games[defaultID]; exists {
		return game
	}
	
	game = NewGame(defaultID, m.config)
	m.games[defaultID] = game
	
	logging.Infow("default_game_created", "game_id", defaultID)
	
	return game
}

// ReplaceDefaultGame replaces the default game instance
func (m *Manager) ReplaceDefaultGame(newGame *Game) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	defaultID := "default"
	
	// Stop old game if exists
	if oldGame, exists := m.games[defaultID]; exists {
		oldGame.Stop()
	}
	
	m.games[defaultID] = newGame
	logging.Infow("default_game_replaced", "game_id", defaultID)
}

// RemoveGame removes a game instance
func (m *Manager) RemoveGame(gameID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	game, exists := m.games[gameID]
	if !exists {
		return ErrGameNotFound
	}
	
	// Stop the game first
	game.Stop()
	
	delete(m.games, gameID)
	
	logging.Infow("game_removed", "game_id", gameID, "remaining_games", len(m.games))
	
	return nil
}

// ListGames returns all active game IDs
func (m *Manager) ListGames() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	ids := make([]string, 0, len(m.games))
	for id := range m.games {
		ids = append(ids, id)
	}
	
	return ids
}

// GetGameCount returns the number of active games
func (m *Manager) GetGameCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	return len(m.games)
}

// Shutdown stops all games and cleans up
func (m *Manager) Shutdown() {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	logging.Infow("manager_shutdown", "game_count", len(m.games))
	
	for _, game := range m.games {
		game.Stop()
	}
	
	m.games = make(map[string]*Game)
}

// GetStats returns statistics about all games
func (m *Manager) GetStats() ManagerStats {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	stats := ManagerStats{
		TotalGames: len(m.games),
		Games:      make([]GameStats, 0, len(m.games)),
	}
	
	for id, game := range m.games {
		state := game.GetState()
		stats.Games = append(stats.Games, GameStats{
			ID:       id,
			Wave:     state.Wave,
			Lives:    state.Lives,
			Score:    state.Score,
			GameOver: state.GameOver,
		})
	}
	
	return stats
}

// ManagerStats contains statistics about the game manager
type ManagerStats struct {
	TotalGames int         `json:"total_games"`
	Games      []GameStats `json:"games"`
}

// GameStats contains statistics about a single game
type GameStats struct {
	ID       string `json:"id"`
	Wave     int    `json:"wave"`
	Lives    int    `json:"lives"`
	Score    int    `json:"score"`
	GameOver bool   `json:"game_over"`
}

// ValidateGameID checks if a game ID is valid
func (m *Manager) ValidateGameID(gameID string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	if _, exists := m.games[gameID]; !exists {
		return fmt.Errorf("%w: %s", ErrGameNotFound, gameID)
	}
	
	return nil
}
