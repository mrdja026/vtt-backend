package combat

import (
        "database/sql"
        "encoding/json"
        "errors"

        "dnd-combat/internal/models"
)

// Repository handles database operations for combat
type Repository struct {
        db *sql.DB
}

// NewRepository creates a new combat repository
func NewRepository(db *sql.DB) *Repository {
        return &Repository{
                db: db,
        }
}

// Create creates a new combat session
func (r *Repository) Create(combat *models.Combat) error {
        // Convert initiative order to JSON
        initiativeJSON, err := json.Marshal(combat.Initiative)
        if err != nil {
                return err
        }

        // Convert participants to JSON
        participantsJSON, err := json.Marshal(combat.Participants)
        if err != nil {
                return err
        }

        // Convert battlefield to JSON
        battlefieldJSON, err := json.Marshal(combat.Battlefield)
        if err != nil {
                return err
        }

        query := `
                INSERT INTO combats (
                        dm_user_id, current_turn_index, round_number, status, 
                        initiative_json, participants_json, battlefield_json, environment,
                        created_at, updated_at
                )
                VALUES (
                        ?, ?, ?, ?, 
                        ?, ?, ?, ?,
                        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
                )
                RETURNING id
        `
        
        err = r.db.QueryRow(
                query,
                combat.DMUserID,
                combat.CurrentTurnIndex,
                combat.RoundNumber,
                combat.Status,
                string(initiativeJSON),
                string(participantsJSON),
                string(battlefieldJSON),
                combat.Environment,
        ).Scan(&combat.ID)

        return err
}

// GetByID retrieves a combat by ID
func (r *Repository) GetByID(id string) (*models.Combat, error) {
        query := `
                SELECT 
                        id, dm_user_id, current_turn_index, round_number, status, 
                        initiative_json, participants_json, battlefield_json, environment,
                        created_at, updated_at
                FROM combats
                WHERE id = ?
                LIMIT 1
        `
        
        combat := &models.Combat{}
        var initiativeJSON, participantsJSON, battlefieldJSON string

        err := r.db.QueryRow(query, id).Scan(
                &combat.ID,
                &combat.DMUserID,
                &combat.CurrentTurnIndex,
                &combat.RoundNumber,
                &combat.Status,
                &initiativeJSON,
                &participantsJSON,
                &battlefieldJSON,
                &combat.Environment,
                &combat.CreatedAt,
                &combat.UpdatedAt,
        )

        if err != nil {
                if errors.Is(err, sql.ErrNoRows) {
                        return nil, nil
                }
                return nil, err
        }

        // Parse initiative JSON
        if initiativeJSON != "" {
                if err := json.Unmarshal([]byte(initiativeJSON), &combat.Initiative); err != nil {
                        return nil, err
                }
        }

        // Parse participants JSON
        if participantsJSON != "" {
                if err := json.Unmarshal([]byte(participantsJSON), &combat.Participants); err != nil {
                        return nil, err
                }
        }

        // Parse battlefield JSON
        if battlefieldJSON != "" {
                if err := json.Unmarshal([]byte(battlefieldJSON), &combat.Battlefield); err != nil {
                        return nil, err
                }
        }

        return combat, nil
}

// Update updates a combat session
func (r *Repository) Update(combat *models.Combat) error {
        // Convert initiative order to JSON
        initiativeJSON, err := json.Marshal(combat.Initiative)
        if err != nil {
                return err
        }

        // Convert participants to JSON
        participantsJSON, err := json.Marshal(combat.Participants)
        if err != nil {
                return err
        }

        // Convert battlefield to JSON
        battlefieldJSON, err := json.Marshal(combat.Battlefield)
        if err != nil {
                return err
        }

        query := `
                UPDATE combats
                SET
                        current_turn_index = ?,
                        round_number = ?,
                        status = ?,
                        initiative_json = ?,
                        participants_json = ?,
                        battlefield_json = ?,
                        updated_at = CURRENT_TIMESTAMP
                WHERE id = ?
        `
        
        result, err := r.db.Exec(
                query,
                combat.CurrentTurnIndex,
                combat.RoundNumber,
                combat.Status,
                string(initiativeJSON),
                string(participantsJSON),
                string(battlefieldJSON),
                combat.ID,
        )

        if err != nil {
                return err
        }

        rows, err := result.RowsAffected()
        if err != nil {
                return err
        }

        if rows == 0 {
                return errors.New("combat not found")
        }

        return nil
}

// SaveAction records a combat action
func (r *Repository) SaveAction(action *models.CombatAction) error {
        // Convert extra data to JSON if it exists
        var extraDataJSON sql.NullString
        if action.ExtraData != nil {
                extraData, err := json.Marshal(action.ExtraData)
                if err != nil {
                        return err
                }
                extraDataJSON.String = string(extraData)
                extraDataJSON.Valid = true
        }

        // Convert movement path to JSON if it exists
        var movementPathJSON sql.NullString
        if action.MovementPath != nil && len(action.MovementPath) > 0 {
                movementPath, err := json.Marshal(action.MovementPath)
                if err != nil {
                        return err
                }
                movementPathJSON.String = string(movementPath)
                movementPathJSON.Valid = true
        }

        // Convert target IDs to JSON if they exist
        var targetIDsJSON sql.NullString
        if action.TargetIDs != nil && len(action.TargetIDs) > 0 {
                targetIDs, err := json.Marshal(action.TargetIDs)
                if err != nil {
                        return err
                }
                targetIDsJSON.String = string(targetIDs)
                targetIDsJSON.Valid = true
        }

        query := `
                INSERT INTO combat_actions (
                        combat_id, actor_id, type, target_ids_json, 
                        spell_id, weapon_name, movement_path_json, extra_data_json,
                        result_description, created_at
                )
                VALUES (
                        ?, ?, ?, ?, 
                        ?, ?, ?, ?,
                        ?, CURRENT_TIMESTAMP
                )
                RETURNING id
        `
        
        var spellID, weaponName sql.NullString
        if action.SpellID != "" {
                spellID.String = action.SpellID
                spellID.Valid = true
        }
        if action.WeaponName != "" {
                weaponName.String = action.WeaponName
                weaponName.Valid = true
        }

        return r.db.QueryRow(
                query,
                action.CombatID,
                action.ActorID,
                action.Type,
                targetIDsJSON,
                spellID,
                weaponName,
                movementPathJSON,
                extraDataJSON,
                action.ResultDescription,
        ).Scan(&action.ID)
}

// GetActionsByCombatID retrieves all actions for a combat
func (r *Repository) GetActionsByCombatID(combatID string) ([]*models.CombatAction, error) {
        query := `
                SELECT 
                        id, combat_id, actor_id, type, target_ids_json, 
                        spell_id, weapon_name, movement_path_json, extra_data_json,
                        result_description, created_at
                FROM combat_actions
                WHERE combat_id = ?
                ORDER BY created_at
        `
        
        rows, err := r.db.Query(query, combatID)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var actions []*models.CombatAction

        for rows.Next() {
                action := &models.CombatAction{}
                var targetIDsJSON, movementPathJSON, extraDataJSON sql.NullString
                var spellID, weaponName sql.NullString

                err := rows.Scan(
                        &action.ID,
                        &action.CombatID,
                        &action.ActorID,
                        &action.Type,
                        &targetIDsJSON,
                        &spellID,
                        &weaponName,
                        &movementPathJSON,
                        &extraDataJSON,
                        &action.ResultDescription,
                        &action.CreatedAt,
                )

                if err != nil {
                        return nil, err
                }

                // Parse target IDs JSON
                if targetIDsJSON.Valid {
                        if err := json.Unmarshal([]byte(targetIDsJSON.String), &action.TargetIDs); err != nil {
                                return nil, err
                        }
                }

                // Parse movement path JSON
                if movementPathJSON.Valid {
                        if err := json.Unmarshal([]byte(movementPathJSON.String), &action.MovementPath); err != nil {
                                return nil, err
                        }
                }

                // Parse extra data JSON
                if extraDataJSON.Valid {
                        var extraData map[string]interface{}
                        if err := json.Unmarshal([]byte(extraDataJSON.String), &extraData); err != nil {
                                return nil, err
                        }
                        action.ExtraData = extraData
                }

                if spellID.Valid {
                        action.SpellID = spellID.String
                }

                if weaponName.Valid {
                        action.WeaponName = weaponName.String
                }

                actions = append(actions, action)
        }

        if err := rows.Err(); err != nil {
                return nil, err
        }

        return actions, nil
}
