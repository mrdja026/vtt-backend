-- Spells table for storing spell information
CREATE TABLE IF NOT EXISTS spells (
    id SERIAL PRIMARY KEY,
    index_name TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    level INTEGER NOT NULL,
    school TEXT NOT NULL,
    casting_time TEXT NOT NULL,
    range TEXT NOT NULL,
    components TEXT NOT NULL,
    materials TEXT,
    duration TEXT NOT NULL,
    concentration BOOLEAN NOT NULL DEFAULT FALSE,
    description TEXT NOT NULL,
    higher_level TEXT,
    attack_type TEXT,
    damage_type TEXT,
    damage_at_level_json JSONB,
    dc_type TEXT,
    dc_success TEXT,
    area_of_effect_json JSONB,
    healing_at_level_json JSONB,
    classes_json JSONB NOT NULL DEFAULT '[]',
    subclasses_json JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for common lookups
CREATE INDEX IF NOT EXISTS idx_spells_name ON spells(name);
CREATE INDEX IF NOT EXISTS idx_spells_level ON spells(level);
CREATE INDEX IF NOT EXISTS idx_spells_school ON spells(school);
CREATE INDEX IF NOT EXISTS idx_spells_concentration ON spells(concentration);

-- Create GIN index for classes (to efficiently find spells available to a class)
CREATE INDEX IF NOT EXISTS idx_spells_classes ON spells USING GIN (classes_json);

-- Create trigger to automatically update the updated_at timestamp
CREATE TRIGGER update_spells_modtime
    BEFORE UPDATE ON spells
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column(); 