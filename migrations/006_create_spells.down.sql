-- Drop the trigger first
DROP TRIGGER IF EXISTS update_spells_modtime ON spells;

-- Drop the indexes
DROP INDEX IF EXISTS idx_spells_name;
DROP INDEX IF EXISTS idx_spells_level;
DROP INDEX IF EXISTS idx_spells_school;
DROP INDEX IF EXISTS idx_spells_concentration;
DROP INDEX IF EXISTS idx_spells_classes;

-- Drop the table
DROP TABLE IF EXISTS spells; 