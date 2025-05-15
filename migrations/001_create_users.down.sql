-- Drop the trigger first
DROP TRIGGER IF EXISTS update_users_modtime ON users;

-- Drop the function
DROP FUNCTION IF EXISTS update_modified_column();

-- Drop the indexes
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_email;

-- Drop the table
DROP TABLE IF EXISTS users; 