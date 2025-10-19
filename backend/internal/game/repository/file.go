package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
)

// FileRepository implements file-based game persistence
type FileRepository struct {
	mu      sync.RWMutex
	baseDir string
}

// NewFileRepository creates a new file-based repository
func NewFileRepository(baseDir string) (*FileRepository, error) {
	// Create base directory if it doesn't exist
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create base directory: %w", err)
	}
	
	return &FileRepository{
		baseDir: baseDir,
	}, nil
}

// Save stores a game state to disk
func (r *FileRepository) Save(gameID string, data []byte) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	saveID := uuid.New().String()
	now := time.Now()
	
	save := &GameSave{
		ID:        saveID,
		GameID:    gameID,
		Data:      data,
		CreatedAt: now,
		UpdatedAt: now,
	}
	
	// Create game directory if it doesn't exist
	gameDir := filepath.Join(r.baseDir, gameID)
	if err := os.MkdirAll(gameDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create game directory: %w", err)
	}
	
	// Write save file
	savePath := filepath.Join(gameDir, fmt.Sprintf("%s.json", saveID))
	saveData, err := json.Marshal(save)
	if err != nil {
		return "", fmt.Errorf("failed to marshal save: %w", err)
	}
	
	if err := os.WriteFile(savePath, saveData, 0644); err != nil {
		return "", fmt.Errorf("failed to write save file: %w", err)
	}
	
	return saveID, nil
}

// Load retrieves a game state from disk
func (r *FileRepository) Load(saveID string) (*GameSave, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	// Search for save file in all game directories
	var savePath string
	err := filepath.Walk(r.baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Base(path) == fmt.Sprintf("%s.json", saveID) {
			savePath = path
			return filepath.SkipAll
		}
		return nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to search for save: %w", err)
	}
	
	if savePath == "" {
		return nil, ErrSaveNotFound
	}
	
	// Read and unmarshal save
	data, err := os.ReadFile(savePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read save file: %w", err)
	}
	
	var save GameSave
	if err := json.Unmarshal(data, &save); err != nil {
		return nil, fmt.Errorf("failed to unmarshal save: %w", err)
	}
	
	return &save, nil
}

// LoadLatest retrieves the latest save for a game
func (r *FileRepository) LoadLatest(gameID string) (*GameSave, error) {
	saves, err := r.List(gameID)
	if err != nil {
		return nil, err
	}
	
	if len(saves) == 0 {
		return nil, ErrSaveNotFound
	}
	
	// Find the latest save
	var latest *GameSave
	for _, save := range saves {
		if latest == nil || save.UpdatedAt.After(latest.UpdatedAt) {
			latest = save
		}
	}
	
	return latest, nil
}

// List returns all saves for a game
func (r *FileRepository) List(gameID string) ([]*GameSave, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	gameDir := filepath.Join(r.baseDir, gameID)
	
	// Check if game directory exists
	if _, err := os.Stat(gameDir); os.IsNotExist(err) {
		return []*GameSave{}, nil
	}
	
	// Read all save files
	entries, err := os.ReadDir(gameDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read game directory: %w", err)
	}
	
	saves := make([]*GameSave, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		
		savePath := filepath.Join(gameDir, entry.Name())
		data, err := os.ReadFile(savePath)
		if err != nil {
			continue // Skip files we can't read
		}
		
		var save GameSave
		if err := json.Unmarshal(data, &save); err != nil {
			continue // Skip invalid saves
		}
		
		saves = append(saves, &save)
	}
	
	return saves, nil
}

// Delete removes a save from disk
func (r *FileRepository) Delete(saveID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// Search for and delete save file
	var savePath string
	err := filepath.Walk(r.baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Base(path) == fmt.Sprintf("%s.json", saveID) {
			savePath = path
			return filepath.SkipAll
		}
		return nil
	})
	
	if err != nil {
		return fmt.Errorf("failed to search for save: %w", err)
	}
	
	if savePath == "" {
		return ErrSaveNotFound
	}
	
	if err := os.Remove(savePath); err != nil {
		return fmt.Errorf("failed to delete save file: %w", err)
	}
	
	return nil
}

// DeleteAll removes all saves for a game
func (r *FileRepository) DeleteAll(gameID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	gameDir := filepath.Join(r.baseDir, gameID)
	
	// Check if game directory exists
	if _, err := os.Stat(gameDir); os.IsNotExist(err) {
		return nil
	}
	
	// Remove entire game directory
	if err := os.RemoveAll(gameDir); err != nil {
		return fmt.Errorf("failed to delete game directory: %w", err)
	}
	
	return nil
}
