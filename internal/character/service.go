package character

import (
	"dnd-combat/internal/models"
)

// Service handles character business logic
type Service struct {
	repo *Repository
}

// NewService creates a new character service
func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// Create creates a new character
func (s *Service) Create(character *models.Character) error {
	return s.repo.Create(character)
}

// GetByID retrieves a character by ID
func (s *Service) GetByID(id string) (*models.Character, error) {
	return s.repo.GetByID(id)
}

// GetByUserID retrieves all characters for a user
func (s *Service) GetByUserID(userID string) ([]*models.Character, error) {
	return s.repo.GetByUserID(userID)
}

// GetMultiple retrieves multiple characters by their IDs
func (s *Service) GetMultiple(ids []string) ([]*models.Character, error) {
	return s.repo.GetMultiple(ids)
}

// Update updates a character
func (s *Service) Update(character *models.Character) error {
	return s.repo.Update(character)
}
