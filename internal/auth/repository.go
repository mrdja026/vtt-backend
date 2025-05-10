package auth

import (
        "database/sql"
        "errors"

        "dnd-combat/internal/models"
        "dnd-combat/pkg/database"
)

// Repository handles database operations for authentication
type Repository struct {
        db *database.DB
}

// NewRepository creates a new auth repository
func NewRepository(db *database.DB) *Repository {
        return &Repository{
                db: db,
        }
}

// Create stores a new user in the database
func (r *Repository) Create(user *models.User) error {
        query := `
                INSERT INTO users (username, email, password_hash, created_at, updated_at)
                VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
                RETURNING id
        `
        
        err := r.db.QueryRow(
                query,
                user.Username,
                user.Email,
                user.PasswordHash,
        ).Scan(&user.ID)

        return err
}

// GetByUsername finds a user by username
func (r *Repository) GetByUsername(username string) (*models.User, error) {
        query := `
                SELECT id, username, email, password_hash, created_at, updated_at
                FROM users
                WHERE username = ?
                LIMIT 1
        `
        
        user := &models.User{}
        err := r.db.QueryRow(query, username).Scan(
                &user.ID,
                &user.Username,
                &user.Email,
                &user.PasswordHash,
                &user.CreatedAt,
                &user.UpdatedAt,
        )

        if err != nil {
                if errors.Is(err, sql.ErrNoRows) {
                        return nil, nil
                }
                return nil, err
        }

        return user, nil
}

// GetByEmail finds a user by email
func (r *Repository) GetByEmail(email string) (*models.User, error) {
        query := `
                SELECT id, username, email, password_hash, created_at, updated_at
                FROM users
                WHERE email = ?
                LIMIT 1
        `
        
        user := &models.User{}
        err := r.db.QueryRow(query, email).Scan(
                &user.ID,
                &user.Username,
                &user.Email,
                &user.PasswordHash,
                &user.CreatedAt,
                &user.UpdatedAt,
        )

        if err != nil {
                if errors.Is(err, sql.ErrNoRows) {
                        return nil, nil
                }
                return nil, err
        }

        return user, nil
}

// GetByID finds a user by ID
func (r *Repository) GetByID(id string) (*models.User, error) {
        query := `
                SELECT id, username, email, password_hash, created_at, updated_at
                FROM users
                WHERE id = ?
                LIMIT 1
        `
        
        user := &models.User{}
        err := r.db.QueryRow(query, id).Scan(
                &user.ID,
                &user.Username,
                &user.Email,
                &user.PasswordHash,
                &user.CreatedAt,
                &user.UpdatedAt,
        )

        if err != nil {
                if errors.Is(err, sql.ErrNoRows) {
                        return nil, nil
                }
                return nil, err
        }

        return user, nil
}
