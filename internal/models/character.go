package models

import (
	"time"
)

// Character represents a D&D character
type Character struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	Name         string    `json:"name"`
	Race         string    `json:"race"`
	Class        string    `json:"class"`
	Level        int       `json:"level"`
	Strength     int       `json:"strength"`
	Dexterity    int       `json:"dexterity"`
	Constitution int       `json:"constitution"`
	Intelligence int       `json:"intelligence"`
	Wisdom       int       `json:"wisdom"`
	Charisma     int       `json:"charisma"`
	HitPoints    int       `json:"hit_points"`
	MaxHitPoints int       `json:"max_hit_points"`
	ArmorClass   int       `json:"armor_class"`
	Equipment    []string  `json:"equipment"`
	Spells       []string  `json:"spells"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// GetAbilityModifier calculates the ability modifier for a given ability score
func GetAbilityModifier(score int) int {
	return (score - 10) / 2
}

// GetProficiencyBonus calculates the proficiency bonus based on character level
func GetProficiencyBonus(level int) int {
	return ((level - 1) / 4) + 2
}
