package game

import (
	"database/sql"
	"encoding/json"
	"errors"

	"dnd-combat/internal/models"
)

// Repository handles database operations for games
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new game repository
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Create stores a new game in the database
func (r *Repository) Create(game *models.Game) error {
	// Convert player IDs to JSON
	playerIDsJSON, err := json.Marshal(game.PlayerIDs)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO games (
			name, description, dm_user_id, player_ids_json, status,
			created_at, updated_at
		)
		VALUES (
			?, ?, ?, ?, ?,
			CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		)
		RETURNING id
	`
	
	err = r.db.QueryRow(
		query,
		game.Name,
		game.Description,
		game.DMUserID,
		string(playerIDsJSON),
		game.Status,
	).Scan(&game.ID)

	return err
}

// GetByID retrieves a game by ID
func (r *Repository) GetByID(id string) (*models.Game, error) {
	query := `
		SELECT 
			id, name, description, dm_user_id, player_ids_json, status,
			created_at, updated_at
		FROM games
		WHERE id = ?
		LIMIT 1
	`
	
	game := &models.Game{}
	var playerIDsJSON string

	err := r.db.QueryRow(query, id).Scan(
		&game.ID,
		&game.Name,
		&game.Description,
		&game.DMUserID,
		&playerIDsJSON,
		&game.Status,
		&game.CreatedAt,
		&game.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	// Parse player IDs JSON
	if playerIDsJSON != "" {
		if err := json.Unmarshal([]byte(playerIDsJSON), &game.PlayerIDs); err != nil {
			return nil, err
		}
	}

	return game, nil
}

// GetByUserID retrieves all games for a user (either as DM or player)
func (r *Repository) GetByUserID(userID string) ([]*models.Game, error) {
	// Query games where user is DM
	query := `
		SELECT 
			id, name, description, dm_user_id, player_ids_json, status,
			created_at, updated_at
		FROM games
		WHERE dm_user_id = ?
		OR player_ids_json LIKE ?
		ORDER BY updated_at DESC
	`
	
	rows, err := r.db.Query(query, userID, "%"+userID+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []*models.Game

	for rows.Next() {
		game := &models.Game{}
		var playerIDsJSON string

		err := rows.Scan(
			&game.ID,
			&game.Name,
			&game.Description,
			&game.DMUserID,
			&playerIDsJSON,
			&game.Status,
			&game.CreatedAt,
			&game.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		// Parse player IDs JSON
		if playerIDsJSON != "" {
			if err := json.Unmarshal([]byte(playerIDsJSON), &game.PlayerIDs); err != nil {
				return nil, err
			}
		}

		// Only add games where user is actually a player
		isDM := game.DMUserID == userID
		isPlayer := false
		for _, playerID := range game.PlayerIDs {
			if playerID == userID {
				isPlayer = true
				break
			}
		}

		if isDM || isPlayer {
			games = append(games, game)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return games, nil
}

// Update updates a game in the database
func (r *Repository) Update(game *models.Game) error {
	// Convert player IDs to JSON
	playerIDsJSON, err := json.Marshal(game.PlayerIDs)
	if err != nil {
		return err
	}

	query := `
		UPDATE games
		SET
			name = ?,
			description = ?,
			player_ids_json = ?,
			status = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	
	result, err := r.db.Exec(
		query,
		game.Name,
		game.Description,
		string(playerIDsJSON),
		game.Status,
		game.ID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("game not found")
	}

	return nil
}
