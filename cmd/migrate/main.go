package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"

	"dnd-combat/config"
)

func main() {
	// Flags
	migrationsPath := flag.String("path", "migrations", "Path to migration files")
	upCommand := flag.Bool("up", false, "Run migrations up")
	downCommand := flag.Bool("down", false, "Roll back all migrations")
	stepFlag := flag.Int("steps", 0, "Number of migrations to apply (up) or roll back (down)")
	versionCommand := flag.Bool("version", false, "Show current migration version")
	flag.Parse()

	// Load environment variables
	godotenv.Load()

	// Initialize configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Validate database type is postgres
	if cfg.DBType != "postgres" {
		log.Fatalf("Migrations currently only supported for PostgreSQL databases")
	}

	// Resolve migrations path
	absPath, err := filepath.Abs(*migrationsPath)
	if err != nil {
		log.Fatalf("Failed to resolve migrations path: %v", err)
	}

	sourceURL := fmt.Sprintf("file://%s", absPath)
	databaseURL := cfg.PostgresConnStr

	// Create a new migrate instance
	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		log.Fatalf("Migration failed to initialize: %v", err)
	}
	defer m.Close()

	// Execute the command
	switch {
	case *versionCommand:
		version, dirty, err := m.Version()
		if err != nil {
			if err == migrate.ErrNilVersion {
				fmt.Println("No migrations applied yet")
				return
			}
			log.Fatalf("Failed to get migration version: %v", err)
		}
		fmt.Printf("Current migration version: %d (dirty: %t)\n", version, dirty)

	case *upCommand && *stepFlag > 0:
		// Apply a specific number of migrations
		if err := m.Steps(*stepFlag); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to apply migrations: %v", err)
		}
		fmt.Printf("Applied %d migrations\n", *stepFlag)

	case *downCommand && *stepFlag > 0:
		// Roll back a specific number of migrations
		if err := m.Steps(-*stepFlag); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to roll back migrations: %v", err)
		}
		fmt.Printf("Rolled back %d migrations\n", *stepFlag)

	case *upCommand:
		// Apply all migrations
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to apply migrations: %v", err)
		}
		fmt.Println("Applied all migrations")

	case *downCommand:
		// Roll back all migrations
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to roll back migrations: %v", err)
		}
		fmt.Println("Rolled back all migrations")

	default:
		flag.Usage()
		os.Exit(1)
	}
}
