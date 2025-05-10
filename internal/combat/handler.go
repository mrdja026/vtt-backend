package combat

import (
        "net/http"

        "github.com/gin-gonic/gin"

        "dnd-combat/internal/models"
        "dnd-combat/pkg/websocket"
)

// Handler handles combat-related HTTP requests
type Handler struct {
        service      *Service
        characterSvc CharacterService
        srdClient    SRDClient
        wsHub        *websocket.Hub
}

// CharacterService defines the interface for character operations
type CharacterService interface {
        GetByID(id string) (*models.Character, error)
        GetMultiple(ids []string) ([]*models.Character, error)
        Update(character *models.Character) error
}

// SRDClient defines the interface for the D&D 5e SRD API client
type SRDClient interface {
        GetMonster(index string) (*models.Monster, error)
        GetSpell(index string) (*models.Spell, error)
}

// NewHandler creates a new combat handler
func NewHandler(service *Service, characterSvc CharacterService, srdClient SRDClient, wsHub *websocket.Hub) *Handler {
        return &Handler{
                service:      service,
                characterSvc: characterSvc,
                srdClient:    srdClient,
                wsHub:        wsHub,
        }
}

// InitiateCombatRequest represents the request body for starting a combat
type InitiateCombatRequest struct {
        ParticipantIDs []string `json:"participants" binding:"required"`
        MonsterIDs     []string `json:"monster_ids"`
        Environment    string   `json:"environment"`
}

// InitiateCombat starts a new combat encounter
func (h *Handler) InitiateCombat(c *gin.Context) {
        var req InitiateCombatRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
                return
        }

        // Get user ID from context (set by auth middleware)
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        // Fetch characters
        characters, err := h.characterSvc.GetMultiple(req.ParticipantIDs)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve characters"})
                return
        }

        // Verify that the user owns the characters or is the DM
        for _, char := range characters {
                if char.UserID != userID.(string) {
                        c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to use these characters"})
                        return
                }
        }

        // Fetch monsters from SRD API
        monsters := make([]*models.Monster, 0, len(req.MonsterIDs))
        for _, monsterID := range req.MonsterIDs {
                monster, err := h.srdClient.GetMonster(monsterID)
                if err != nil {
                        c.JSON(http.StatusInternalServerError, gin.H{
                                "error": "Failed to fetch monster data",
                                "details": err.Error(),
                                "monster_id": monsterID,
                        })
                        return
                }
                monsters = append(monsters, monster)
        }

        // Create combat session
        combat, err := h.service.CreateCombat(characters, monsters, req.Environment, userID.(string))
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create combat session"})
                return
        }

        // Broadcast combat state to websocket clients
        h.wsHub.BroadcastToRoom(combat.ID, websocket.Message{
                Type: "combat_initiated",
                Data: combat,
        })

        c.JSON(http.StatusCreated, combat)
}

// GetCombat retrieves the current state of a combat
func (h *Handler) GetCombat(c *gin.Context) {
        id := c.Param("id")
        if id == "" {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Combat ID is required"})
                return
        }

        // Get user ID from context (set by auth middleware)
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        // Get combat session
        combat, err := h.service.GetCombat(id)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve combat session"})
                return
        }

        if combat == nil {
                c.JSON(http.StatusNotFound, gin.H{"error": "Combat session not found"})
                return
        }

        // Check if user is involved in the combat (DM or has a character)
        if !h.service.IsUserInCombat(combat, userID.(string)) {
                c.JSON(http.StatusForbidden, gin.H{"error": "You don't have access to this combat session"})
                return
        }

        c.JSON(http.StatusOK, combat)
}

// CombatActionRequest represents the request body for a combat action
type CombatActionRequest struct {
        ActionType   string                 `json:"action_type" binding:"required"`
        ActorID      string                 `json:"actor_id" binding:"required"`
        TargetIDs    []string               `json:"target_ids"`
        SpellID      string                 `json:"spell_id"`
        WeaponName   string                 `json:"weapon_name"`
        MovementPath [][2]int               `json:"movement_path"`
        ExtraData    map[string]interface{} `json:"extra_data"`
}

// PerformAction executes a combat action
func (h *Handler) PerformAction(c *gin.Context) {
        id := c.Param("id")
        if id == "" {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Combat ID is required"})
                return
        }

        var req CombatActionRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
                return
        }

        // Get user ID from context (set by auth middleware)
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        // Get combat session
        combat, err := h.service.GetCombat(id)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve combat session"})
                return
        }

        if combat == nil {
                c.JSON(http.StatusNotFound, gin.H{"error": "Combat session not found"})
                return
        }

        // Check if it's the actor's turn
        if !h.service.IsActorsTurn(combat, req.ActorID) {
                c.JSON(http.StatusBadRequest, gin.H{"error": "It's not this actor's turn"})
                return
        }

        // Check if user controls the actor
        if !h.service.UserControlsActor(combat, userID.(string), req.ActorID) {
                c.JSON(http.StatusForbidden, gin.H{"error": "You don't control this actor"})
                return
        }

        // Execute the action
        result, err := h.service.ExecuteAction(combat, &models.CombatAction{
                CombatID:     id,
                Type:         req.ActionType,
                ActorID:      req.ActorID,
                TargetIDs:    req.TargetIDs,
                SpellID:      req.SpellID,
                WeaponName:   req.WeaponName,
                MovementPath: req.MovementPath,
                ExtraData:    req.ExtraData,
        })

        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to execute action", "details": err.Error()})
                return
        }

        // Update combat state
        combat, err = h.service.GetCombat(id)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated combat session"})
                return
        }

        // Broadcast updated combat state to websocket clients
        h.wsHub.BroadcastToRoom(combat.ID, websocket.Message{
                Type: "combat_updated",
                Data: combat,
        })

        // Also broadcast the action result
        h.wsHub.BroadcastToRoom(combat.ID, websocket.Message{
                Type: "action_result",
                Data: result,
        })

        c.JSON(http.StatusOK, gin.H{
                "action_result": result,
                "combat": combat,
        })
}

// EndTurnRequest represents the request to end a participant's turn
type EndTurnRequest struct {
        ActorID string `json:"actor_id" binding:"required"`
}

// EndTurn processes the end of a turn
func (h *Handler) EndTurn(c *gin.Context) {
        id := c.Param("id")
        if id == "" {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Combat ID is required"})
                return
        }

        var req EndTurnRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
                return
        }

        // Get user ID from context (set by auth middleware)
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        // Get combat session
        combat, err := h.service.GetCombat(id)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve combat session"})
                return
        }

        if combat == nil {
                c.JSON(http.StatusNotFound, gin.H{"error": "Combat session not found"})
                return
        }

        // Check if it's the actor's turn
        if !h.service.IsActorsTurn(combat, req.ActorID) {
                c.JSON(http.StatusBadRequest, gin.H{"error": "It's not this actor's turn"})
                return
        }

        // Check if user controls the actor
        if !h.service.UserControlsActor(combat, userID.(string), req.ActorID) {
                c.JSON(http.StatusForbidden, gin.H{"error": "You don't control this actor"})
                return
        }

        // End the turn
        if err := h.service.EndTurn(combat); err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to end turn"})
                return
        }

        // Get updated combat
        combat, err = h.service.GetCombat(id)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated combat session"})
                return
        }

        // Broadcast updated combat state to websocket clients
        h.wsHub.BroadcastToRoom(combat.ID, websocket.Message{
                Type: "combat_updated",
                Data: combat,
        })

        c.JSON(http.StatusOK, combat)
}

// WebSocketHandler handles websocket connections for a specific combat
func (h *Handler) WebSocketHandler(c *gin.Context) {
        id := c.Param("id")
        if id == "" {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Combat ID is required"})
                return
        }

        // Get user ID from context (set by auth middleware)
        userID, exists := c.Get("userID")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
                return
        }

        // Get combat session
        combat, err := h.service.GetCombat(id)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve combat session"})
                return
        }

        if combat == nil {
                c.JSON(http.StatusNotFound, gin.H{"error": "Combat session not found"})
                return
        }

        // Check if user is involved in the combat
        if !h.service.IsUserInCombat(combat, userID.(string)) {
                c.JSON(http.StatusForbidden, gin.H{"error": "You don't have access to this combat session"})
                return
        }

        // Upgrade connection to websocket
        h.wsHub.ServeWs(c.Writer, c.Request, id, userID.(string))
}
