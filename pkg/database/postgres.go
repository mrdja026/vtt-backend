package database

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// PostgresDB is a wrapper around sql.DB with additional methods
type PostgresDB struct {
	*sql.DB
}

// Make sure PostgresDB implements DBInterface
var _ DBInterface = (*PostgresDB)(nil)

// NewPostgres establishes a new PostgreSQL database connection using the provided connection string.
// It configures connection pooling parameters and verifies connectivity by pinging the database.
// Returns a PostgresDB instance on success or an error if the connection cannot be established.
func NewPostgres(connString string) (*PostgresDB, error) {
	// Open database connection
	sqlDB, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresDB{DB: sqlDB}, nil
}

// RunMigrations applies database migrations from the specified directory
func (db *PostgresDB) RunMigrations(migrationsPath string) error {
	// Create migration instance
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	// Resolve absolute path to migrations
	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to resolve migrations path: %w", err)
	}

	// Initialize migrations
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", absPath),
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize migrations: %w", err)
	}

	// Apply all migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Println("Database migrations applied successfully")
	return nil
}

// Close closes the database connection
func (db *PostgresDB) Close() error {
	return db.DB.Close()
}
