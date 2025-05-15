-- Games table for managing game sessions
CREATE TABLE IF NOT EXISTS games (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    dm_user_id INTEGER NOT NULL,
    player_ids_json JSONB NOT NULL DEFAULT '[]',
    status TEXT NOT NULL,
    settings_json JSONB NOT NULL DEFAULT '{}',
    active_combat_id INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (dm_user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Create index for dm_user_id lookups (to find games a user is DMing)
CREATE INDEX IF NOT EXISTS idx_games_dm_user_id ON games(dm_user_id);

-- Create GIN index for player_ids_json field (to efficiently find games a player is in)
CREATE INDEX IF NOT EXISTS idx_games_player_ids ON games USING GIN (player_ids_json);

-- Create trigger to automatically update the updated_at timestamp
CREATE TRIGGER update_games_modtime
    BEFORE UPDATE ON games
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column(); 