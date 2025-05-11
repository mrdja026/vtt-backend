-- Templates table for storing reusable content templates
CREATE TABLE IF NOT EXISTS templates (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    type TEXT NOT NULL, -- 'character', 'monster', 'spell', etc.
    content_json JSONB NOT NULL,
    is_public BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Create index for user_id lookups (to find templates belonging to a user)
CREATE INDEX IF NOT EXISTS idx_templates_user_id ON templates(user_id);

-- Create index for type lookups (to find templates of a specific type)
CREATE INDEX IF NOT EXISTS idx_templates_type ON templates(type);

-- Create index for public templates
CREATE INDEX IF NOT EXISTS idx_templates_public ON templates(is_public) WHERE is_public = TRUE;

-- Create composite index for user+name (users should be able to find their templates by name quickly)
CREATE INDEX IF NOT EXISTS idx_templates_user_name ON templates(user_id, name);

-- Create trigger to automatically update the updated_at timestamp
CREATE TRIGGER update_templates_modtime
    BEFORE UPDATE ON templates
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column(); 