-- Drop the trigger first
DROP TRIGGER IF EXISTS update_games_modtime ON games;

-- Drop the indexes
DROP INDEX IF EXISTS idx_games_dm_user_id;
DROP INDEX IF EXISTS idx_games_player_ids;

-- Drop the table
DROP TABLE IF EXISTS games; 