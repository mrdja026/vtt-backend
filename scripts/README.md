# Database Setup Scripts

This directory contains scripts to set up the PostgreSQL database for the D&D Combat API.

## Requirements

- PostgreSQL 12+ installed and running
- `psql` command-line tool available in your PATH
- Superuser access to your PostgreSQL instance

## Setup Instructions

### Windows

1. Open Command Prompt or PowerShell
2. Navigate to the project root directory
3. Run the setup script:
   ```batch
   cd vtt_ttrpg_golang_backend
   scripts\setup_database.bat
   ```

### Linux/macOS

1. Open Terminal
2. Navigate to the project root directory
3. Make the script executable and run it:
   ```bash
   cd vtt_ttrpg_golang_backend
   chmod +x scripts/setup_database.sh
   ./scripts/setup_database.sh
   ```

## What These Scripts Do

1. Create a new database called `dnd_combat`
2. Create a new database user `dnd_app_user` with password `dnd_password`
3. Grant necessary permissions to the user
4. Set environment variables required by the application

## Manual Setup

If you prefer to set up the database manually, you can:

1. Run the SQL script directly:
   ```
   psql -U postgres -f scripts/create_database.sql
   ```

2. Set the environment variables:
   - Windows:
     ```batch
     set DB_TYPE=postgres
     set DATABASE_URL=postgres://dnd_app_user:dnd_password@localhost:5432/dnd_combat?sslmode=disable
     ```
   - Linux/macOS:
     ```bash
     export DB_TYPE=postgres
     export DATABASE_URL=postgres://dnd_app_user:dnd_password@localhost:5432/dnd_combat?sslmode=disable
     ```

## Customizing Database Configuration

If you want to use different database credentials or settings:

1. Edit the `create_database.sql` file to change database name, user, or password
2. Update the connection string in the setup scripts accordingly
3. Update your application's environment variables to match the changes

## Running Migrations

After setting up the database, you can apply migrations:

```bash
go run cmd/migrate/main.go -up
```

This will create all the tables defined in the migration files. 