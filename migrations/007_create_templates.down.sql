-- Drop the trigger first
DROP TRIGGER IF EXISTS update_templates_modtime ON templates;

-- Drop the indexes
DROP INDEX IF EXISTS idx_templates_user_id;
DROP INDEX IF EXISTS idx_templates_type;
DROP INDEX IF EXISTS idx_templates_public;
DROP INDEX IF EXISTS idx_templates_user_name;

-- Drop the table
DROP TABLE IF EXISTS templates; 