package game

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dnd-combat/internal/models"
)

// Handler handles game-related HTTP requests
type Handler struct {
	service *Service
}

// NewHandler creates a new game handler
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// CreateGameRequest represents the request body for game creation
type CreateGameRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	PlayerIDs   []string `json:"player_ids"`
}

// Create handles game session creation
func (h *Handler) Create(c *gin.Context) {
	var req CreateGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Get the user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	game := &models.Game{
		Name:        req.Name,
		Description: req.Description,
		DMUserID:    userID.(string),
		PlayerIDs:   req.PlayerIDs,
		Status:      "active",
	}

	if err := h.service.Create(game); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create game session"})
		return
	}

	c.JSON(http.StatusCreated, game)
}

// Get retrieves a game session by ID
func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Game ID is required"})
		return
	}

	// Get the user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	game, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve game session"})
		return
	}

	if game == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game session not found"})
		return
	}

	// Check if the user is the DM or a player in the game
	if game.DMUserID != userID.(string) && !contains(game.PlayerIDs, userID.(string)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to view this game"})
		return
	}

	c.JSON(http.StatusOK, game)
}

// List retrieves all game sessions for a user
func (h *Handler) List(c *gin.Context) {
	// Get the user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	games, err := h.service.GetByUserID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve game sessions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"games": games})
}

// Update updates a game session
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Game ID is required"})
		return
	}

	var req CreateGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Get the user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get the game first to check permissions
	game, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve game session"})
		return
	}

	if game == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game session not found"})
		return
	}

	// Only the DM can update the game
	if game.DMUserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only the DM can update the game"})
		return
	}

	// Update the game
	game.Name = req.Name
	game.Description = req.Description
	game.PlayerIDs = req.PlayerIDs

	if err := h.service.Update(game); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update game session"})
		return
	}

	c.JSON(http.StatusOK, game)
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
