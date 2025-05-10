package models

import (
	"time"
)

// Combat represents a D&D combat session
type Combat struct {
	ID               string       `json:"id"`
	DMUserID         string       `json:"dm_user_id"`
	Initiative       []string     `json:"initiative"`
	Participants     []*Combatant `json:"participants"`
	CurrentTurnIndex int          `json:"current_turn_index"`
	RoundNumber      int          `json:"round_number"`
	Status           string       `json:"status"` // active, victory, defeat
	Environment      string       `json:"environment"`
	Battlefield      *Battlefield `json:"battlefield"`
	CreatedAt        time.Time    `json:"created_at"`
	UpdatedAt        time.Time    `json:"updated_at"`
}

// Combatant represents a participant in a combat
type Combatant struct {
	ID          string      `json:"id"`
	Type        string      `json:"type"` // character, monster
	Name        string      `json:"name"`
	Initiative  int         `json:"initiative"`
	HP          int         `json:"hp"`
	MaxHP       int         `json:"max_hp"`
	AC          int         `json:"ac"`
	Position    []int       `json:"position"` // [x, y] on the battlefield
	UserID      string      `json:"user_id"`
	CharacterID string      `json:"character_id,omitempty"`
	MonsterID   string      `json:"monster_id,omitempty"`
	Stats       interface{} `json:"stats"` // Either Character or Monster
	Conditions  []string    `json:"conditions"`
}

// Battlefield represents the combat field
type Battlefield struct {
	Width   int        `json:"width"`
	Height  int        `json:"height"`
	Terrain [][]string `json:"terrain"` // normal, difficult, water, etc.
	Objects [][]string `json:"objects"` // wall, tree, rock, etc.
}

// CombatAction represents an action taken in combat
type CombatAction struct {
	ID                string      `json:"id"`
	CombatID          string      `json:"combat_id"`
	ActorID           string      `json:"actor_id"`
	Type              string      `json:"type"` // attack, cast_spell, move, etc.
	TargetIDs         []string    `json:"target_ids,omitempty"`
	SpellID           string      `json:"spell_id,omitempty"`
	WeaponName        string      `json:"weapon_name,omitempty"`
	MovementPath      [][]int     `json:"movement_path,omitempty"`
	ExtraData         interface{} `json:"extra_data,omitempty"`
	ResultDescription string      `json:"result_description"`
	CreatedAt         time.Time   `json:"created_at"`
}

// ActionResult represents the outcome of a combat action
type ActionResult struct {
	Success     bool        `json:"success"`
	Damage      int         `json:"damage,omitempty"`
	Healing     int         `json:"healing,omitempty"`
	Effects     []string    `json:"effects,omitempty"`
	Description string      `json:"description"`
	ExtraData   interface{} `json:"extra_data,omitempty"`
}

// Monster represents a D&D monster
type Monster struct {
	Index        string         `json:"index"`
	Name         string         `json:"name"`
	Size         string         `json:"size"`
	Type         string         `json:"type"`
	Alignment    string         `json:"alignment"`
	ArmorClass   int            `json:"armor_class"`
	HitDice      string         `json:"hit_dice"`
	Speed        MonsterSpeed   `json:"speed"`
	Strength     int            `json:"strength"`
	Dexterity    int            `json:"dexterity"`
	Constitution int            `json:"constitution"`
	Intelligence int            `json:"intelligence"`
	Wisdom       int            `json:"wisdom"`
	Charisma     int            `json:"charisma"`
	StrengthMod  int            `json:"strength_mod"`
	DexterityMod int            `json:"dexterity_mod"`
	ConMod       int            `json:"constitution_mod"`
	IntMod       int            `json:"intelligence_mod"`
	WisdomMod    int            `json:"wisdom_mod"`
	CharismaMod  int            `json:"charisma_mod"`
	Actions      []MonsterAction `json:"actions"`
	ChallengeRating float64     `json:"challenge_rating"`
	XP           int            `json:"xp"`
}

// MonsterSpeed represents a monster's speed capabilities
type MonsterSpeed struct {
	Walk   int `json:"walk"`
	Swim   int `json:"swim,omitempty"`
	Fly    int `json:"fly,omitempty"`
	Climb  int `json:"climb,omitempty"`
	Burrow int `json:"burrow,omitempty"`
}

// MonsterAction represents an action a monster can take
type MonsterAction struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	AttackBonus int         `json:"attack_bonus,omitempty"`
	Range       int         `json:"range,omitempty"`
	Damage      DamageInfo  `json:"damage,omitempty"`
}

// DamageInfo represents damage dealt by an attack or spell
type DamageInfo struct {
	DiceCount int    `json:"dice_count"`
	DiceValue int    `json:"dice_value"`
	Bonus     int    `json:"bonus"`
	Type      string `json:"type"`
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
