-- Characters table for storing player characters
CREATE TABLE IF NOT EXISTS characters (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
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
    speed INTEGER NOT NULL DEFAULT 30,
    initiative_bonus INTEGER NOT NULL DEFAULT 0,
    equipment_json JSONB NOT NULL DEFAULT '[]',
    spells_json JSONB NOT NULL DEFAULT '[]',
    features_json JSONB NOT NULL DEFAULT '[]',
    portrait_url TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Create index for user_id lookups (to find characters belonging to a user)
CREATE INDEX IF NOT EXISTS idx_characters_user_id ON characters(user_id);

-- Create trigger to automatically update the updated_at timestamp
CREATE TRIGGER update_characters_modtime
    BEFORE UPDATE ON characters
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column(); 