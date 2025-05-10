package models

import (
	"time"
)

// User represents a user account
type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Game represents a D&D game session
type Game struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	DMUserID    string    `json:"dm_user_id"`
	PlayerIDs   []string  `json:"player_ids"`
	Status      string    `json:"status"` // active, completed, paused
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
