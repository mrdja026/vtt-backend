-- Drop the indexes
DROP INDEX IF EXISTS idx_combat_actions_combat_id;
DROP INDEX IF EXISTS idx_combat_actions_actor_id;
DROP INDEX IF EXISTS idx_combat_actions_type;

-- Drop the table
DROP TABLE IF EXISTS combat_actions; 