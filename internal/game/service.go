package game

import (
	"dnd-combat/internal/models"
)

// Service handles game business logic
type Service struct {
	repo *Repository
}

// NewService creates a new game service
func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// Create creates a new game session
func (s *Service) Create(game *models.Game) error {
	return s.repo.Create(game)
}

// GetByID retrieves a game by ID
func (s *Service) GetByID(id string) (*models.Game, error) {
	return s.repo.GetByID(id)
}

// GetByUserID retrieves all games for a user
func (s *Service) GetByUserID(userID string) ([]*models.Game, error) {
	return s.repo.GetByUserID(userID)
}

// Update updates a game
func (s *Service) Update(game *models.Game) error {
	return s.repo.Update(game)
}
