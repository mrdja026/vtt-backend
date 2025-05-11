-- Combats table for managing combat encounters
CREATE TABLE IF NOT EXISTS combats (
    id SERIAL PRIMARY KEY,
    game_id INTEGER,
    dm_user_id INTEGER NOT NULL,
    name TEXT NOT NULL DEFAULT 'Combat Encounter',
    current_turn_index INTEGER NOT NULL DEFAULT 0,
    round_number INTEGER NOT NULL DEFAULT 1,
    status TEXT NOT NULL,
    initiative_json JSONB NOT NULL DEFAULT '[]',
    participants_json JSONB NOT NULL DEFAULT '[]',
    battlefield_json JSONB NOT NULL DEFAULT '{}',
    environment TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (dm_user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (game_id) REFERENCES games (id) ON DELETE SET NULL
);

-- Create index for dm_user_id lookups
CREATE INDEX IF NOT EXISTS idx_combats_dm_user_id ON combats(dm_user_id);

-- Create index for game_id lookups
CREATE INDEX IF NOT EXISTS idx_combats_game_id ON combats(game_id);

-- Create index for status lookups (to find active combats)
CREATE INDEX IF NOT EXISTS idx_combats_status ON combats(status);

-- Create trigger to automatically update the updated_at timestamp
CREATE TRIGGER update_combats_modtime
    BEFORE UPDATE ON combats
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();

-- Add foreign key reference from games to combats for active_combat_id
ALTER TABLE games
ADD CONSTRAINT fk_games_active_combat
FOREIGN KEY (active_combat_id) 
REFERENCES combats (id)
ON DELETE SET NULL; 