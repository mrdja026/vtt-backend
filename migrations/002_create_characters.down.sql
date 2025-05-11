-- Drop the trigger first
DROP TRIGGER IF EXISTS update_characters_modtime ON characters;

-- Drop the index
DROP INDEX IF EXISTS idx_characters_user_id;

-- Drop the table
DROP TABLE IF EXISTS characters; 