#!/bin/bash
echo "Setting up PostgreSQL database for D&D Combat API..."

# Check if psql is available
if ! command -v psql &> /dev/null; then
    echo "ERROR: PostgreSQL command-line tools not found."
    echo "Please ensure PostgreSQL is installed and psql is in your PATH."
    exit 1
fi

# Run the SQL script to create the database
echo "Creating database and user..."
psql -U postgres -f scripts/create_database.sql

if [ $? -ne 0 ]; then
    echo "ERROR: Failed to create database. See error message above."
    exit 1
fi

# Set environment variable for the application
echo "Setting environment variables..."
export DATABASE_URL="postgres://dnd_app_user:dnd_password@localhost:5432/dnd_combat?sslmode=disable"
export DB_TYPE="postgres"

echo
echo "Database setup complete!"
echo
echo "You can now run the application with:"
echo "  go run cmd/api/main.go"
echo
echo "Or apply migrations manually with:"
echo "  go run cmd/migrate/main.go -up"
echo
echo "Connection string (set as DATABASE_URL):"
echo "  $DATABASE_URL"
echo
echo "For running PostgreSQL migrations during development,"
echo "you may need to set the following environment variables:"
echo "  export DB_TYPE=postgres"
echo "  export DATABASE_URL=postgres://dnd_app_user:dnd_password@localhost:5432/dnd_combat?sslmode=disable" 