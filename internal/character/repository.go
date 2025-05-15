package character

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"

	"dnd-combat/internal/models"
	"dnd-combat/pkg/database"
)

// Repository handles database operations for characters
type Repository struct {
	db database.DBInterface
}

// NewRepository creates a new character repository
func NewRepository(db database.DBInterface) *Repository {
	return &Repository{
		db: db,
	}
}

// Create stores a new character in the database
func (r *Repository) Create(character *models.Character) error {
	// Convert equipment to JSON
	equipmentJSON, err := json.Marshal(character.Equipment)
	if err != nil {
		return err
	}

	// Convert spells to JSON
	spellsJSON, err := json.Marshal(character.Spells)
	if err != nil {
		return err
	}

	query := `
                INSERT INTO characters (
                        user_id, name, race, class, level, 
                        strength, dexterity, constitution, intelligence, wisdom, charisma,
                        hit_points, max_hit_points, armor_class, equipment_json, spells_json,
                        created_at, updated_at
                )
                VALUES (
                        ?, ?, ?, ?, ?, 
                        ?, ?, ?, ?, ?, ?,
                        ?, ?, ?, ?, ?,
                        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
                )
                RETURNING id
        `

	err = r.db.QueryRow(
		query,
		character.UserID,
		character.Name,
		character.Race,
		character.Class,
		character.Level,
		character.Strength,
		character.Dexterity,
		character.Constitution,
		character.Intelligence,
		character.Wisdom,
		character.Charisma,
		character.HitPoints,
		character.MaxHitPoints,
		character.ArmorClass,
		string(equipmentJSON),
		string(spellsJSON),
	).Scan(&character.ID)

	return err
}

// GetByID retrieves a character by ID
func (r *Repository) GetByID(id string) (*models.Character, error) {
	query := `
                SELECT 
                        id, user_id, name, race, class, level,
                        strength, dexterity, constitution, intelligence, wisdom, charisma,
                        hit_points, max_hit_points, armor_class, equipment_json, spells_json,
                        created_at, updated_at
                FROM characters
                WHERE id = ?
                LIMIT 1
        `

	character := &models.Character{}
	var equipmentJSON, spellsJSON string

	err := r.db.QueryRow(query, id).Scan(
		&character.ID,
		&character.UserID,
		&character.Name,
		&character.Race,
		&character.Class,
		&character.Level,
		&character.Strength,
		&character.Dexterity,
		&character.Constitution,
		&character.Intelligence,
		&character.Wisdom,
		&character.Charisma,
		&character.HitPoints,
		&character.MaxHitPoints,
		&character.ArmorClass,
		&equipmentJSON,
		&spellsJSON,
		&character.CreatedAt,
		&character.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	// Parse equipment JSON
	if equipmentJSON != "" {
		if err := json.Unmarshal([]byte(equipmentJSON), &character.Equipment); err != nil {
			return nil, err
		}
	}

	// Parse spells JSON
	if spellsJSON != "" {
		if err := json.Unmarshal([]byte(spellsJSON), &character.Spells); err != nil {
			return nil, err
		}
	}

	return character, nil
}

// GetByUserID retrieves all characters for a user
func (r *Repository) GetByUserID(userID string) ([]*models.Character, error) {
	query := `
                SELECT 
                        id, user_id, name, race, class, level,
                        strength, dexterity, constitution, intelligence, wisdom, charisma,
                        hit_points, max_hit_points, armor_class, equipment_json, spells_json,
                        created_at, updated_at
                FROM characters
                WHERE user_id = ?
                ORDER BY name
        `

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var characters []*models.Character

	for rows.Next() {
		character := &models.Character{}
		var equipmentJSON, spellsJSON string

		err := rows.Scan(
			&character.ID,
			&character.UserID,
			&character.Name,
			&character.Race,
			&character.Class,
			&character.Level,
			&character.Strength,
			&character.Dexterity,
			&character.Constitution,
			&character.Intelligence,
			&character.Wisdom,
			&character.Charisma,
			&character.HitPoints,
			&character.MaxHitPoints,
			&character.ArmorClass,
			&equipmentJSON,
			&spellsJSON,
			&character.CreatedAt,
			&character.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		// Parse equipment JSON
		if equipmentJSON != "" {
			if err := json.Unmarshal([]byte(equipmentJSON), &character.Equipment); err != nil {
				return nil, err
			}
		}

		// Parse spells JSON
		if spellsJSON != "" {
			if err := json.Unmarshal([]byte(spellsJSON), &character.Spells); err != nil {
				return nil, err
			}
		}

		characters = append(characters, character)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return characters, nil
}

// GetMultiple retrieves multiple characters by their IDs
func (r *Repository) GetMultiple(ids []string) ([]*models.Character, error) {
	if len(ids) == 0 {
		return []*models.Character{}, nil
	}

	// Create placeholders for the IN clause
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))

	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	query := `
                SELECT 
                        id, user_id, name, race, class, level,
                        strength, dexterity, constitution, intelligence, wisdom, charisma,
                        hit_points, max_hit_points, armor_class, equipment_json, spells_json,
                        created_at, updated_at
                FROM characters
                WHERE id IN (` + placeholders[0] + strings.Repeat(", ?", len(placeholders)-1) + `)
        `

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var characters []*models.Character

	for rows.Next() {
		character := &models.Character{}
		var equipmentJSON, spellsJSON string

		err := rows.Scan(
			&character.ID,
			&character.UserID,
			&character.Name,
			&character.Race,
			&character.Class,
			&character.Level,
			&character.Strength,
			&character.Dexterity,
			&character.Constitution,
			&character.Intelligence,
			&character.Wisdom,
			&character.Charisma,
			&character.HitPoints,
			&character.MaxHitPoints,
			&character.ArmorClass,
			&equipmentJSON,
			&spellsJSON,
			&character.CreatedAt,
			&character.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		// Parse equipment JSON
		if equipmentJSON != "" {
			if err := json.Unmarshal([]byte(equipmentJSON), &character.Equipment); err != nil {
				return nil, err
			}
		}

		// Parse spells JSON
		if spellsJSON != "" {
			if err := json.Unmarshal([]byte(spellsJSON), &character.Spells); err != nil {
				return nil, err
			}
		}

		characters = append(characters, character)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return characters, nil
}

// Update updates a character in the database
func (r *Repository) Update(character *models.Character) error {
	// Convert equipment to JSON
	equipmentJSON, err := json.Marshal(character.Equipment)
	if err != nil {
		return err
	}

	// Convert spells to JSON
	spellsJSON, err := json.Marshal(character.Spells)
	if err != nil {
		return err
	}

	query := `
                UPDATE characters
                SET
                        name = ?,
                        race = ?,
                        class = ?,
                        level = ?,
                        strength = ?,
                        dexterity = ?,
                        constitution = ?,
                        intelligence = ?,
                        wisdom = ?,
                        charisma = ?,
                        hit_points = ?,
                        max_hit_points = ?,
                        armor_class = ?,
                        equipment_json = ?,
                        spells_json = ?,
                        updated_at = CURRENT_TIMESTAMP
                WHERE id = ?
        `

	result, err := r.db.Exec(
		query,
		character.Name,
		character.Race,
		character.Class,
		character.Level,
		character.Strength,
		character.Dexterity,
		character.Constitution,
		character.Intelligence,
		character.Wisdom,
		character.Charisma,
		character.HitPoints,
		character.MaxHitPoints,
		character.ArmorClass,
		string(equipmentJSON),
		string(spellsJSON),
		character.ID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("character not found")
	}

	return nil
}
