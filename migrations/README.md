# Database Migrations

This directory contains SQL migration files for the PostgreSQL database used by the D&D Combat API.

## Migration Files

Migrations are organized in numbered pairs of files:
- `NNN_description.up.sql`: Contains SQL statements to apply the migration
- `NNN_description.down.sql`: Contains SQL statements to roll back the migration

Current migrations:
- `001_create_users`: Creates the users table and related indexes
- `002_create_characters`: Creates the characters table and related indexes
- `003_create_games`: Creates the games table and related indexes
- `004_create_combats`: Creates the combats table and related indexes
- `005_create_combat_actions`: Creates the combat_actions table and related indexes
- `006_create_spells`: Creates the spells table and related indexes
- `007_create_templates`: Creates the templates table and related indexes

## Running Migrations

Migrations are automatically applied when the server starts with PostgreSQL database type.

To manually run or rollback migrations, use the migrate tool:

```bash
# Apply all migrations
go run cmd/migrate/main.go -up

# Apply a specific number of migrations
go run cmd/migrate/main.go -up -steps 2

# Rollback all migrations
go run cmd/migrate/main.go -down

# Rollback a specific number of migrations
go run cmd/migrate/main.go -down -steps 1

# Show current migration version
go run cmd/migrate/main.go -version
```

## Creating New Migrations

To create a new migration, add a pair of files with the next sequential number:

```bash
touch migrations/008_description.up.sql migrations/008_description.down.sql
```

Follow these guidelines:
1. Always create both up and down migrations
2. Use meaningful descriptions
3. Include indexes for fields that will be frequently queried
4. Add foreign key constraints where appropriate
5. Add automated triggers for updated_at timestamps

## Migration Best Practices

1. **Idempotent Scripts**: Migrations should be idempotent where possible (using `IF EXISTS` and `IF NOT EXISTS`)
2. **Atomic Changes**: Each migration should represent a single logical change
3. **Test Migrations**: Test both up and down migrations before committing
4. **Documenting Changes**: Add comments to complex migrations
5. **Avoid Data Loss**: Migrations should preserve existing data