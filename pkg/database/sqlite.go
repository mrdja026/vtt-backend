package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// DB is a wrapper around sql.DB with additional methods
type DB struct {
	*sql.DB
}

// Make sure DB implements DBInterface
var _ DBInterface = (*DB)(nil)

// New initializes a new SQLite database at the specified path, ensuring required directories and schema exist.
// Returns a DB instance connected to the database, or an error if setup fails.
func New(dbPath string) (*DB, error) {
	// Ensure the directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database connection
	sqlDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Enable foreign keys
	if _, err := sqlDB.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Create tables
	if err := createTables(sqlDB); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return &DB{DB: sqlDB}, nil
}

// createTables initializes the required tables in the SQLite database if they do not already exist.
// 
// Returns an error if any table creation fails.
func createTables(db *sql.DB) error {
	// Create users table
	if _, err := db.Exec(`
                CREATE TABLE IF NOT EXISTS users (
                        id TEXT PRIMARY KEY,
                        username TEXT UNIQUE NOT NULL,
                        email TEXT UNIQUE NOT NULL,
                        password_hash TEXT NOT NULL,
                        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
                )
        `); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// Create characters table
	if _, err := db.Exec(`
                CREATE TABLE IF NOT EXISTS characters (
                        id TEXT PRIMARY KEY,
                        user_id TEXT NOT NULL,
                        name TEXT NOT NULL,
                        race TEXT NOT NULL,
                        class TEXT NOT NULL,
                        level INTEGER NOT NULL,
                        strength INTEGER NOT NULL,
                        dexterity INTEGER NOT NULL,
                        constitution INTEGER NOT NULL,
                        intelligence INTEGER NOT NULL,
                        wisdom INTEGER NOT NULL,
                        charisma INTEGER NOT NULL,
                        hit_points INTEGER NOT NULL,
                        max_hit_points INTEGER NOT NULL,
                        armor_class INTEGER NOT NULL,
                        equipment_json TEXT NOT NULL,
                        spells_json TEXT NOT NULL,
                        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
                )
        `); err != nil {
		return fmt.Errorf("failed to create characters table: %w", err)
	}

	// Create games table
	if _, err := db.Exec(`
                CREATE TABLE IF NOT EXISTS games (
                        id TEXT PRIMARY KEY,
                        name TEXT NOT NULL,
                        description TEXT NOT NULL,
                        dm_user_id TEXT NOT NULL,
                        player_ids_json TEXT NOT NULL,
                        status TEXT NOT NULL,
                        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        FOREIGN KEY (dm_user_id) REFERENCES users (id) ON DELETE CASCADE
                )
        `); err != nil {
		return fmt.Errorf("failed to create games table: %w", err)
	}

	// Create combats table
	if _, err := db.Exec(`
                CREATE TABLE IF NOT EXISTS combats (
                        id TEXT PRIMARY KEY,
                        dm_user_id TEXT NOT NULL,
                        current_turn_index INTEGER NOT NULL,
                        round_number INTEGER NOT NULL,
                        status TEXT NOT NULL,
                        initiative_json TEXT NOT NULL,
                        participants_json TEXT NOT NULL,
                        battlefield_json TEXT NOT NULL,
                        environment TEXT NOT NULL,
                        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        FOREIGN KEY (dm_user_id) REFERENCES users (id) ON DELETE CASCADE
                )
        `); err != nil {
		return fmt.Errorf("failed to create combats table: %w", err)
	}

	// Create combat_actions table
	if _, err := db.Exec(`
                CREATE TABLE IF NOT EXISTS combat_actions (
                        id TEXT PRIMARY KEY,
                        combat_id TEXT NOT NULL,
                        actor_id TEXT NOT NULL,
                        type TEXT NOT NULL,
                        target_ids_json TEXT,
                        spell_id TEXT,
                        weapon_name TEXT,
                        movement_path_json TEXT,
                        extra_data_json TEXT,
                        result_description TEXT NOT NULL,
                        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        FOREIGN KEY (combat_id) REFERENCES combats (id) ON DELETE CASCADE
                )
        `); err != nil {
		return fmt.Errorf("failed to create combat_actions table: %w", err)
	}

	return nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}
