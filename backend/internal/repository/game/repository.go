package repository

import (
	"encoding/json"
	"errors"
	"time"
)

var (
	ErrSaveNotFound = errors.New("save not found")
	ErrInvalidData  = errors.New("invalid save data")
)

// GameSave represents a saved game state
type GameSave struct {
	ID        string    `json:"id"`
	GameID    string    `json:"game_id"`
	Data      []byte    `json:"data"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Repository defines the interface for game persistence
type Repository interface {
	// Save stores a game state
	Save(gameID string, data []byte) (string, error)
	
	// Load retrieves a game state by save ID
	Load(saveID string) (*GameSave, error)
	
	// LoadLatest retrieves the latest save for a game
	LoadLatest(gameID string) (*GameSave, error)
	
	// List returns all saves for a game
	List(gameID string) ([]*GameSave, error)
	
	// Delete removes a save
	Delete(saveID string) error
	
	// DeleteAll removes all saves for a game
	DeleteAll(gameID string) error
}

// SaveMetadata contains metadata about a game save
type SaveMetadata struct {
	Wave     int       `json:"wave"`
	Gold     int       `json:"gold"`
	Lives    int       `json:"lives"`
	Score    int       `json:"score"`
	GameOver bool      `json:"game_over"`
	SavedAt  time.Time `json:"saved_at"`
}

// ExtractMetadata extracts metadata from save data
func ExtractMetadata(data []byte) (*SaveMetadata, error) {
	var state struct {
		Wave     int  `json:"wave"`
		Gold     int  `json:"gold"`
		Lives    int  `json:"lives"`
		Score    int  `json:"score"`
		GameOver bool `json:"gameOver"`
	}
	
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	
	return &SaveMetadata{
		Wave:     state.Wave,
		Gold:     state.Gold,
		Lives:    state.Lives,
		Score:    state.Score,
		GameOver: state.GameOver,
		SavedAt:  time.Now(),
	}, nil
}
