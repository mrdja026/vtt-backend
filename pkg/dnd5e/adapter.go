package dnd5e

import (
	"dnd-combat/internal/models"
)

// SRDClientAdapter adapts the SRDClient to implement the combat.SRDClient interface
type SRDClientAdapter struct {
	client *SRDClient
}

// NewSRDClientAdapter creates a new adapter for SRDClient
func NewSRDClientAdapter(client *SRDClient) *SRDClientAdapter {
	return &SRDClientAdapter{
		client: client,
	}
}

// GetMonster fetches a monster from the SRD API and converts it to models.Monster
func (a *SRDClientAdapter) GetMonster(index string) (*models.Monster, error) {
	monster, err := a.client.GetMonster(index)
	if err != nil {
		return nil, err
	}

	// Convert from pkg/dnd5e.Monster to models.Monster
	actions := make([]models.MonsterAction, 0, len(monster.Actions))
	for _, action := range monster.Actions {
		actions = append(actions, models.MonsterAction{
			Name:        action.Name,
			Description: action.Description,
			AttackBonus: action.AttackBonus,
			Damage: models.DamageInfo{
				DiceCount: action.Damage.DiceCount,
				DiceValue: action.Damage.DiceValue,
				Bonus:     action.Damage.Bonus,
				Type:      action.Damage.Type,
			},
		})
	}

	return &models.Monster{
		Index:        monster.Index,
		Name:         monster.Name,
		Size:         monster.Size,
		Type:         monster.Type,
		Alignment:    monster.Alignment,
		ArmorClass:   monster.ArmorClass,
		HitDice:      monster.HitDice,
		Speed: models.MonsterSpeed{
			Walk:   monster.Speed.Walk,
			Swim:   monster.Speed.Swim,
			Fly:    monster.Speed.Fly,
			Climb:  monster.Speed.Climb,
			Burrow: monster.Speed.Burrow,
		},
		Strength:       monster.Strength,
		Dexterity:      monster.Dexterity,
		Constitution:   monster.Constitution,
		Intelligence:   monster.Intelligence,
		Wisdom:         monster.Wisdom,
		Charisma:       monster.Charisma,
		StrengthMod:    monster.StrengthMod,
		DexterityMod:   monster.DexterityMod,
		ConMod:         monster.ConMod,
		IntMod:         monster.IntMod,
		WisdomMod:      monster.WisdomMod,
		CharismaMod:    monster.CharismaMod,
		Actions:        actions,
		ChallengeRating: monster.ChallengeRating,
		XP:             monster.XP,
	}, nil
}

// GetSpell fetches a spell from the SRD API and converts it to models.Spell
func (a *SRDClientAdapter) GetSpell(index string) (*models.Spell, error) {
	spell, err := a.client.GetSpell(index)
	if err != nil {
		return nil, err
	}

	return &models.Spell{
		Index:       spell.Index,
		Name:        spell.Name,
		Level:       spell.Level,
		School:      spell.School,
		CastingTime: spell.CastingTime,
		Range:       spell.Range,
		Components:  spell.Components,
		Duration:    spell.Duration,
		Description: spell.Description,
		HigherLevel: spell.HigherLevel,
		Classes:     spell.Classes,
	}, nil
}