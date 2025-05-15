package database

import (
	"database/sql"
	"fmt"
)

// DBInterface defines the common interface for database implementations
type DBInterface interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Prepare(query string) (*sql.Stmt, error)
	Begin() (*sql.Tx, error)
	Close() error
}

// MigratableDB defines an interface for databases that support migrations
type MigratableDB interface {
	DBInterface
	RunMigrations(migrationsPath string) error
}

// NewDatabase creates a new database connection based on the provided configuration
func NewDatabase(dbType, sqlitePath, postgresConnStr string) (DBInterface, error) {
	switch dbType {
	case "sqlite":
		return New(sqlitePath)
	case "postgres":
		return NewPostgres(postgresConnStr)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}
