package repository

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// MemoryRepository implements in-memory game persistence
type MemoryRepository struct {
	mu    sync.RWMutex
	saves map[string]*GameSave
	index map[string][]string // gameID -> []saveID
}

// NewMemoryRepository creates a new in-memory repository
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		saves: make(map[string]*GameSave),
		index: make(map[string][]string),
	}
}

// Save stores a game state
func (r *MemoryRepository) Save(gameID string, data []byte) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	saveID := uuid.New().String()
	now := time.Now()
	
	save := &GameSave{
		ID:        saveID,
		GameID:    gameID,
		Data:      make([]byte, len(data)),
		CreatedAt: now,
		UpdatedAt: now,
	}
	copy(save.Data, data)
	
	r.saves[saveID] = save
	r.index[gameID] = append(r.index[gameID], saveID)
	
	return saveID, nil
}

// Load retrieves a game state by save ID
func (r *MemoryRepository) Load(saveID string) (*GameSave, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	save, exists := r.saves[saveID]
	if !exists {
		return nil, ErrSaveNotFound
	}
	
	// Return a copy
	result := &GameSave{
		ID:        save.ID,
		GameID:    save.GameID,
		Data:      make([]byte, len(save.Data)),
		CreatedAt: save.CreatedAt,
		UpdatedAt: save.UpdatedAt,
	}
	copy(result.Data, save.Data)
	
	return result, nil
}

// LoadLatest retrieves the latest save for a game
func (r *MemoryRepository) LoadLatest(gameID string) (*GameSave, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	saveIDs, exists := r.index[gameID]
	if !exists || len(saveIDs) == 0 {
		return nil, ErrSaveNotFound
	}
	
	// Get the last save (most recent)
	latestID := saveIDs[len(saveIDs)-1]
	save := r.saves[latestID]
	
	// Return a copy
	result := &GameSave{
		ID:        save.ID,
		GameID:    save.GameID,
		Data:      make([]byte, len(save.Data)),
		CreatedAt: save.CreatedAt,
		UpdatedAt: save.UpdatedAt,
	}
	copy(result.Data, save.Data)
	
	return result, nil
}

// List returns all saves for a game
func (r *MemoryRepository) List(gameID string) ([]*GameSave, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	saveIDs, exists := r.index[gameID]
	if !exists {
		return []*GameSave{}, nil
	}
	
	result := make([]*GameSave, 0, len(saveIDs))
	for _, saveID := range saveIDs {
		if save, exists := r.saves[saveID]; exists {
			// Return a copy
			copied := &GameSave{
				ID:        save.ID,
				GameID:    save.GameID,
				Data:      make([]byte, len(save.Data)),
				CreatedAt: save.CreatedAt,
				UpdatedAt: save.UpdatedAt,
			}
			copy(copied.Data, save.Data)
			result = append(result, copied)
		}
	}
	
	return result, nil
}

// Delete removes a save
func (r *MemoryRepository) Delete(saveID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	save, exists := r.saves[saveID]
	if !exists {
		return ErrSaveNotFound
	}
	
	// Remove from saves map
	delete(r.saves, saveID)
	
	// Remove from index
	gameID := save.GameID
	if saveIDs, exists := r.index[gameID]; exists {
		newIDs := make([]string, 0, len(saveIDs)-1)
		for _, id := range saveIDs {
			if id != saveID {
				newIDs = append(newIDs, id)
			}
		}
		r.index[gameID] = newIDs
	}
	
	return nil
}

// DeleteAll removes all saves for a game
func (r *MemoryRepository) DeleteAll(gameID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	saveIDs, exists := r.index[gameID]
	if !exists {
		return nil
	}
	
	// Remove all saves
	for _, saveID := range saveIDs {
		delete(r.saves, saveID)
	}
	
	// Clear index
	delete(r.index, gameID)
	
	return nil
}

// GetStats returns statistics about the repository
func (r *MemoryRepository) GetStats() RepositoryStats {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	totalSize := 0
	for _, save := range r.saves {
		totalSize += len(save.Data)
	}
	
	return RepositoryStats{
		TotalSaves: len(r.saves),
		TotalGames: len(r.index),
		TotalBytes: totalSize,
	}
}

// RepositoryStats contains statistics about the repository
type RepositoryStats struct {
	TotalSaves int `json:"total_saves"`
	TotalGames int `json:"total_games"`
	TotalBytes int `json:"total_bytes"`
}
