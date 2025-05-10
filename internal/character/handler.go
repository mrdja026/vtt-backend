package character

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dnd-combat/internal/models"
)

// Handler handles character-related HTTP requests
type Handler struct {
	service *Service
}

// NewHandler creates a new character handler
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// CreateRequest represents the request body for character creation
type CreateRequest struct {
	Name         string `json:"name" binding:"required"`
	Race         string `json:"race" binding:"required"`
	Class        string `json:"class" binding:"required"`
	Level        int    `json:"level" binding:"required,min=1,max=20"`
	Strength     int    `json:"strength" binding:"required,min=3,max=20"`
	Dexterity    int    `json:"dexterity" binding:"required,min=3,max=20"`
	Constitution int    `json:"constitution" binding:"required,min=3,max=20"`
	Intelligence int    `json:"intelligence" binding:"required,min=3,max=20"`
	Wisdom       int    `json:"wisdom" binding:"required,min=3,max=20"`
	Charisma     int    `json:"charisma" binding:"required,min=3,max=20"`
	HitPoints    int    `json:"hit_points" binding:"required,min=1"`
	ArmorClass   int    `json:"armor_class" binding:"required,min=1"`
	Equipment    []string `json:"equipment"`
	Spells       []string `json:"spells"`
}

// Create handles character creation
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
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

	character := &models.Character{
		UserID:       userID.(string),
		Name:         req.Name,
		Race:         req.Race,
		Class:        req.Class,
		Level:        req.Level,
		Strength:     req.Strength,
		Dexterity:    req.Dexterity,
		Constitution: req.Constitution,
		Intelligence: req.Intelligence,
		Wisdom:       req.Wisdom,
		Charisma:     req.Charisma,
		HitPoints:    req.HitPoints,
		MaxHitPoints: req.HitPoints,
		ArmorClass:   req.ArmorClass,
		Equipment:    req.Equipment,
		Spells:       req.Spells,
	}

	if err := h.service.Create(character); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create character"})
		return
	}

	c.JSON(http.StatusCreated, character)
}

// Get retrieves a character by ID
func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Character ID is required"})
		return
	}

	// Get the user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	character, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve character"})
		return
	}

	if character == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
		return
	}

	// Check if the character belongs to the authenticated user
	if character.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to view this character"})
		return
	}

	c.JSON(http.StatusOK, character)
}

// List retrieves all characters owned by the authenticated user
func (h *Handler) List(c *gin.Context) {
	// Get the user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	characters, err := h.service.GetByUserID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve characters"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"characters": characters})
}
