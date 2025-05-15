@echo off
echo Setting up PostgreSQL database for D&D Combat API...

REM Check if psql is available
where psql >nul 2>nul
if %ERRORLEVEL% neq 0 (
    echo ERROR: PostgreSQL command-line tools not found.
    echo Please ensure PostgreSQL is installed and psql is in your PATH.
    exit /b 1
)

REM Run the SQL script to create the database
echo Creating database and user...
psql -U postgres -f scripts\create_database.sql

if %ERRORLEVEL% neq 0 (
    echo ERROR: Failed to create database. See error message above.
    exit /b 1
)

REM Set environment variable for the application
echo Setting environment variables...
set DATABASE_URL=postgres://dnd_app_user:dnd_password@localhost:5432/dnd_combat?sslmode=disable
set DB_TYPE=postgres

echo.
echo Database setup complete!
echo.
echo You can now run the application with:
echo   go run cmd/api/main.go
echo.
echo Or apply migrations manually with:
echo   go run cmd/migrate/main.go -up
echo.
echo Connection string (set as DATABASE_URL):
echo   %DATABASE_URL%
echo.
echo For running PostgreSQL migrations during development,
echo you may need to set the following environment variables:
echo   set DB_TYPE=postgres
echo   set DATABASE_URL=postgres://dnd_app_user:dnd_password@localhost:5432/dnd_combat?sslmode=disable 