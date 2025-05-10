package models

import (
        "time"
)

// Monster represents a D&D monster
type Monster struct {
        Index        string       `json:"index"`
        Name         string       `json:"name"`
        Size         string       `json:"size"`
        Type         string       `json:"type"`
        Alignment    string       `json:"alignment"`
        ArmorClass   int          `json:"armor_class"`
        HitDice      string       `json:"hit_dice"`
        Speed        MonsterSpeed `json:"speed"`
        Strength     int          `json:"strength"`
        Dexterity    int          `json:"dexterity"`
        Constitution int          `json:"constitution"`
        Intelligence int          `json:"intelligence"`
        Wisdom       int          `json:"wisdom"`
        Charisma     int          `json:"charisma"`
        StrengthMod  int          `json:"strength_mod"`
        DexterityMod int          `json:"dexterity_mod"`
        ConMod       int          `json:"constitution_mod"`
        IntMod       int          `json:"intelligence_mod"`
        WisdomMod    int          `json:"wisdom_mod"`
        CharismaMod  int          `json:"charisma_mod"`
        Actions      []MonsterAction `json:"actions"`
        ChallengeRating float64   `json:"challenge_rating"`
        XP           int          `json:"xp"`
}

// MonsterSpeed represents monster movement speeds
type MonsterSpeed struct {
        Walk   int `json:"walk"`
        Swim   int `json:"swim,omitempty"`
        Fly    int `json:"fly,omitempty"`
        Climb  int `json:"climb,omitempty"`
        Burrow int `json:"burrow,omitempty"`
}

// DamageInfo represents damage information for an attack
type DamageInfo struct {
        DiceCount int    `json:"dice_count"`
        DiceValue int    `json:"dice_value"`
        Bonus     int    `json:"bonus"`
        Type      string `json:"type"`
}

// MonsterAction represents an action a monster can take
type MonsterAction struct {
        Name        string     `json:"name"`
        Description string     `json:"description"`
        AttackBonus int        `json:"attack_bonus"`
        Range       int        `json:"range,omitempty"`
        Damage      DamageInfo `json:"damage,omitempty"`
}

// Spell represents a D&D spell
type Spell struct {
        Index       string   `json:"index"`
        Name        string   `json:"name"`
        Level       int      `json:"level"`
        School      string   `json:"school"`
        CastingTime string   `json:"casting_time"`
        Range       string   `json:"range"`
        Components  []string `json:"components"`
        Duration    string   `json:"duration"`
        Description string   `json:"description"`
        HigherLevel string   `json:"higher_level,omitempty"`
        Classes     []string `json:"classes"`
}

// Combat represents a D&D combat encounter
type Combat struct {
        ID              string           `json:"id"`
        DMUserID        string           `json:"dm_user_id"`
        CurrentTurnIndex int             `json:"current_turn_index"` // Renamed from CurrentTurnIdx for consistency
        RoundNumber     int              `json:"round_number"`
        Status          string           `json:"status"`
        Initiative      []InitiativeItem `json:"initiative"`
        Participants    []Combatant      `json:"participants"`
        Battlefield     Battlefield      `json:"battlefield"`
        Environment     string           `json:"environment"`
        CreatedAt       time.Time        `json:"created_at"`
        UpdatedAt       time.Time        `json:"updated_at"`
}

// InitiativeItem represents a participant's initiative order
type InitiativeItem struct {
        ID          string `json:"id"`
        Name        string `json:"name"`
        Initiative  int    `json:"initiative"`
        Dexterity   int    `json:"dexterity"`
        IsCharacter bool   `json:"is_character"`
}

// Combatant represents a participant in combat
type Combatant struct {
        ID           string      `json:"id"`
        UserID       string      `json:"user_id,omitempty"` // Added for characters owned by users
        CharacterID  string      `json:"character_id,omitempty"` // ID reference to character
        MonsterID    string      `json:"monster_id,omitempty"` // ID reference to monster
        Name         string      `json:"name"`
        Type         string      `json:"type"` // "character" or "monster"
        HP           int         `json:"hp"`
        MaxHP        int         `json:"max_hp"`
        AC           int         `json:"ac"`
        Initiative   int         `json:"initiative"`
        Position     [2]int      `json:"position"`
        Conditions   []string    `json:"conditions"`
        Stats        interface{} `json:"stats,omitempty"` // Character or Monster
}

// Battlefield represents the combat area
type Battlefield struct {
        Width     int                `json:"width"`
        Height    int                `json:"height"`
        Grid      map[string]string  `json:"grid"` // Map of "x,y" to content
        Terrain   map[string]string  `json:"terrain"`
        Obstacles map[string]bool    `json:"obstacles"`
}

// CombatAction represents an action taken in combat
type CombatAction struct {
        ID               string                 `json:"id"`
        CombatID         string                 `json:"combat_id"`
        ActorID          string                 `json:"actor_id"`
        Type             string                 `json:"type"` // "attack", "cast", "move", etc.
        TargetIDs        []string               `json:"target_ids,omitempty"`
        SpellID          string                 `json:"spell_id,omitempty"`
        WeaponName       string                 `json:"weapon_name,omitempty"`
        MovementPath     [][2]int               `json:"movement_path,omitempty"`
        ExtraData        map[string]interface{} `json:"extra_data,omitempty"`
        ResultDescription string                `json:"result_description,omitempty"`
        CreatedAt        time.Time              `json:"created_at"`
}

// ActionResult represents the result of a combat action
type ActionResult struct {
        Success      bool         `json:"success"`
        Description  string       `json:"description"`
        Damage       int          `json:"damage,omitempty"`
        DamageType   string       `json:"damage_type,omitempty"`
        Healing      int          `json:"healing,omitempty"`
        TargetEffect string       `json:"target_effect,omitempty"`
        Errors       []string     `json:"errors,omitempty"`
}