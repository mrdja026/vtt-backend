-- First remove the foreign key constraint from games table
ALTER TABLE games
DROP CONSTRAINT IF EXISTS fk_games_active_combat;

-- Drop the triggers
DROP TRIGGER IF EXISTS update_combats_modtime ON combats;

-- Drop the indexes
DROP INDEX IF EXISTS idx_combats_dm_user_id;
DROP INDEX IF EXISTS idx_combats_game_id;
DROP INDEX IF EXISTS idx_combats_status;

-- Drop the table
DROP TABLE IF EXISTS combats; 