-- Combat actions table for tracking combat history
CREATE TABLE IF NOT EXISTS combat_actions (
    id SERIAL PRIMARY KEY,
    combat_id INTEGER NOT NULL,
    actor_id TEXT NOT NULL,
    type TEXT NOT NULL,
    target_ids_json JSONB,
    spell_id TEXT,
    weapon_name TEXT,
    movement_path_json JSONB,
    damage INTEGER,
    damage_type TEXT,
    healing INTEGER,
    success BOOLEAN,
    roll INTEGER,
    target_effect TEXT,
    extra_data_json JSONB,
    result_description TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (combat_id) REFERENCES combats (id) ON DELETE CASCADE
);

-- Create index for combat_id lookups (to get all actions in a combat)
CREATE INDEX IF NOT EXISTS idx_combat_actions_combat_id ON combat_actions(combat_id);

-- Create index for actor_id lookups (to get all actions by a specific actor)
CREATE INDEX IF NOT EXISTS idx_combat_actions_actor_id ON combat_actions(actor_id);

-- Create index for type lookups (to get all actions of a specific type)
CREATE INDEX IF NOT EXISTS idx_combat_actions_type ON combat_actions(type); 