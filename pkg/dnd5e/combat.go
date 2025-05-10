package dnd5e

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

// CombatRules handles D&D 5e combat rules
type CombatRules struct {
	diceRoller *DiceRoller
}

// NewCombatRules creates a new combat rules engine
func NewCombatRules(diceRoller *DiceRoller) *CombatRules {
	return &CombatRules{
		diceRoller: diceRoller,
	}
}

// AttackResult represents the result of an attack
type AttackResult struct {
	AttackRoll     int
	TotalAttack    int
	HitResult      bool
	IsCritical     bool
	IsCriticalMiss bool
	Damage         int
	DamageType     string
	Description    string
}

// MeleeAttack simulates a melee attack
func (c *CombatRules) MeleeAttack(attackerName, targetName, weaponName string, attackBonus, damageBonus int, targetAC int, damageDice string, damageType string) AttackResult {
	return c.attack(attackerName, targetName, weaponName, attackBonus, damageBonus, targetAC, damageDice, damageType, false)
}

// RangedAttack simulates a ranged attack
func (c *CombatRules) RangedAttack(attackerName, targetName, weaponName string, attackBonus, damageBonus int, targetAC int, damageDice string, damageType string) AttackResult {
	return c.attack(attackerName, targetName, weaponName, attackBonus, damageBonus, targetAC, damageDice, damageType, true)
}

// attack is a helper function for both melee and ranged attacks
func (c *CombatRules) attack(attackerName, targetName, weaponName string, attackBonus, damageBonus int, targetAC int, damageDice string, damageType string, isRanged bool) AttackResult {
	// Roll attack
	attackRoll := c.diceRoller.Roll(1, 20)
	totalAttack := attackRoll + attackBonus
	
	// Check for critical hit or miss
	isCritical := attackRoll == 20
	isCritMiss := attackRoll == 1
	
	result := AttackResult{
		AttackRoll:     attackRoll,
		TotalAttack:    totalAttack,
		IsCritical:     isCritical,
		IsCriticalMiss: isCritMiss,
		DamageType:     damageType,
	}
	
	// Critical miss always fails
	if isCritMiss {
		result.HitResult = false
		result.Description = fmt.Sprintf("%s critically misses their %s attack with %s against %s!", 
			attackerName, isRanged ? "ranged" : "melee", weaponName, targetName)
		return result
	}
	
	// Check if attack hits
	hits := isCritical || totalAttack >= targetAC
	result.HitResult = hits
	
	if !hits {
		result.Description = fmt.Sprintf("%s attacks %s with %s but misses! (Rolled %d + %d = %d vs AC %d)", 
			attackerName, targetName, weaponName, attackRoll, attackBonus, totalAttack, targetAC)
		return result
	}
	
	// Roll damage
	damage := c.rollDamage(damageDice, damageBonus, isCritical)
	result.Damage = damage
	
	// Create description
	result.Description = fmt.Sprintf("%s %s %s with %s for %d %s damage!", 
		attackerName, 
		isCritical ? "critically hits" : "hits",
		targetName, 
		weaponName, 
		damage,
		damageType)
	
	return result
}

// rollDamage calculates damage for an attack
func (c *CombatRules) rollDamage(damageDice string, damageBonus int, isCritical bool) int {
	// Parse the damage dice
	parts := strings.Split(damageDice, "d")
	if len(parts) != 2 {
		return damageBonus // Invalid format, return just the bonus
	}
	
	diceCount := 0
	diceValue := 0
	
	fmt.Sscanf(parts[0], "%d", &diceCount)
	fmt.Sscanf(parts[1], "%d", &diceValue)
	
	if diceCount <= 0 || diceValue <= 0 {
		return damageBonus // Invalid values, return just the bonus
	}
	
	// Roll damage
	baseDamage := c.diceRoller.Roll(diceCount, diceValue)
	
	// Double dice for critical hit
	if isCritical {
		baseDamage += c.diceRoller.Roll(diceCount, diceValue)
	}
	
	return baseDamage + damageBonus
}

// SpellCastResult represents the result of a spell cast
type SpellCastResult struct {
	Success      bool
	SpellName    string
	Effect       string
	SaveDC       int
	Damage       int
	DamageType   string
	Healing      int
	TargetEffect string
	Description  string
}

// CastDamageSpell simulates casting a damage-dealing spell
func (c *CombatRules) CastDamageSpell(casterName, targetName, spellName string, damageDice string, damageType string, saveDC int, saveAbilityMod int, halfDamageOnSave bool) SpellCastResult {
	result := SpellCastResult{
		Success:   true,
		SpellName: spellName,
		SaveDC:    saveDC,
		DamageType: damageType,
		Effect:    "damage",
	}
	
	// Roll damage
	damage := c.diceRoller.RollDamage(damageDice)
	result.Damage = damage
	
	// Check if target makes a saving throw
	if saveDC > 0 {
		// Roll saving throw
		saveRoll := c.diceRoller.Roll(1, 20) + saveAbilityMod
		savePassed := saveRoll >= saveDC
		
		if savePassed {
			if halfDamageOnSave {
				// Half damage on successful save
				result.Damage = int(math.Floor(float64(damage) / 2))
				result.TargetEffect = "half damage (save successful)"
			} else {
				// No damage on successful save
				result.Damage = 0
				result.TargetEffect = "no damage (save successful)"
			}
			result.Description = fmt.Sprintf("%s casts %s at %s! %s makes a successful saving throw for %s damage.", 
				casterName, spellName, targetName, targetName, result.TargetEffect)
		} else {
			result.TargetEffect = "full damage (save failed)"
			result.Description = fmt.Sprintf("%s casts %s at %s for %d %s damage! %s fails their saving throw.", 
				casterName, spellName, targetName, damage, damageType, targetName)
		}
	} else {
		// No saving throw required (like Magic Missile)
		result.TargetEffect = "full damage (no save)"
		result.Description = fmt.Sprintf("%s casts %s at %s for %d %s damage!", 
			casterName, spellName, targetName, damage, damageType)
	}
	
	return result
}

// CastHealingSpell simulates casting a healing spell
func (c *CombatRules) CastHealingSpell(casterName, targetName, spellName string, healingDice string, bonus int) SpellCastResult {
	result := SpellCastResult{
		Success:   true,
		SpellName: spellName,
		Effect:    "healing",
	}
	
	// Roll healing
	healing := c.diceRoller.RollDamage(healingDice) + bonus
	result.Healing = healing
	
	result.Description = fmt.Sprintf("%s casts %s on %s, healing %d hit points!", 
		casterName, spellName, targetName, healing)
	
	return result
}

// CastBuffSpell simulates casting a buff spell
func (c *CombatRules) CastBuffSpell(casterName, targetName, spellName, effect string, duration int) SpellCastResult {
	result := SpellCastResult{
		Success:   true,
		SpellName: spellName,
		Effect:    "buff",
	}
	
	var durationText string
	if duration == 1 {
		durationText = "1 round"
	} else {
		durationText = fmt.Sprintf("%d rounds", duration)
	}
	
	result.TargetEffect = effect
	
	if targetName == casterName {
		result.Description = fmt.Sprintf("%s casts %s, gaining %s for %s!", 
			casterName, spellName, effect, durationText)
	} else {
		result.Description = fmt.Sprintf("%s casts %s on %s, granting %s for %s!", 
			casterName, spellName, targetName, effect, durationText)
	}
	
	return result
}

// MovementResult represents the result of a movement action
type MovementResult struct {
	Success     bool
	StartPos    []int
	EndPos      []int
	Distance    int
	Description string
	Error       error
}

// ValidateMovement checks if a movement path is valid
func (c *CombatRules) ValidateMovement(actorName string, startPos []int, endPos []int, movementPath [][]int, maxDistance int, occupiedPositions map[string][]int, obstacles map[string][]int) MovementResult {
	result := MovementResult{
		Success:  true,
		StartPos: startPos,
		EndPos:   endPos,
	}
	
	// Check if path is empty
	if len(movementPath) == 0 {
		result.Success = false
		result.Error = errors.New("movement path is empty")
		result.Description = fmt.Sprintf("%s attempts to move but provides no path", actorName)
		return result
	}
	
	// Check if start position matches actor's position
	if startPos[0] != movementPath[0][0] || startPos[1] != movementPath[0][1] {
		result.Success = false
		result.Error = errors.New("movement path doesn't start at actor's position")
		result.Description = fmt.Sprintf("%s attempts to move from an invalid starting position", actorName)
		return result
	}
	
	// Check if end position matches the last position in the path
	lastPos := movementPath[len(movementPath)-1]
	if endPos[0] != lastPos[0] || endPos[1] != lastPos[1] {
		result.Success = false
		result.Error = errors.New("movement path doesn't end at the specified end position")
		result.Description = fmt.Sprintf("%s's movement path doesn't lead to the intended destination", actorName)
		return result
	}
	
	// Calculate total distance
	totalDistance := 0
	currentPos := startPos
	
	for i, pos := range movementPath {
		if i == 0 {
			// Skip the starting position
			continue
		}
		
		// Check for diagonal movement (costs 1.5 times as much in D&D 5e)
		dx := abs(pos[0] - currentPos[0])
		dy := abs(pos[1] - currentPos[1])
		
		if dx > 1 || dy > 1 {
			result.Success = false
			result.Error = errors.New("invalid movement: can only move to adjacent squares")
			result.Description = fmt.Sprintf("%s attempts an invalid movement step from [%d,%d] to [%d,%d]", 
				actorName, currentPos[0], currentPos[1], pos[0], pos[1])
			return result
		}
		
		// Check for obstacles
		posKey := fmt.Sprintf("%d,%d", pos[0], pos[1])
		if _, exists := obstacles[posKey]; exists {
			result.Success = false
			result.Error = errors.New("movement path is blocked by an obstacle")
			result.Description = fmt.Sprintf("%s's path is blocked by an obstacle at [%d,%d]", 
				actorName, pos[0], pos[1])
			return result
		}
		
		// Check for other combatants
		for name, position := range occupiedPositions {
			if name != actorName && position[0] == pos[0] && position[1] == pos[1] {
				result.Success = false
				result.Error = errors.New("movement path is blocked by another combatant")
				result.Description = fmt.Sprintf("%s's path is blocked by %s at [%d,%d]", 
					actorName, name, pos[0], pos[1])
				return result
			}
		}
		
		// Calculate movement cost (diagonal = 1.5 squares, round down)
		moveCost := 1
		if dx == 1 && dy == 1 {
			// Diagonal movement costs more
			moveCost = 1
			// In a more accurate implementation, you could track half-squares
			// and charge 2 squares for every other diagonal move
		}
		
		totalDistance += moveCost
		currentPos = pos
	}
	
	// Check if movement exceeds maximum distance
	if totalDistance > maxDistance {
		result.Success = false
		result.Error = fmt.Errorf("movement exceeds maximum distance (used %d, max %d)", totalDistance, maxDistance)
		result.Description = fmt.Sprintf("%s attempts to move too far (used %d squares, max %d)", 
			actorName, totalDistance, maxDistance)
		return result
	}
	
	result.Distance = totalDistance
	result.Description = fmt.Sprintf("%s moves from [%d,%d] to [%d,%d] (%d squares)", 
		actorName, startPos[0], startPos[1], endPos[0], endPos[1], totalDistance)
	
	return result
}

// Condition represents a temporary condition affecting a combatant
type Condition struct {
	Name       string
	Duration   int
	Effect     string
	SaveDC     int
	SaveType   string
	EndOfTurn  bool
	StartOfTurn bool
}

// ProcessCondition handles condition effects and duration
func (c *CombatRules) ProcessCondition(condition *Condition, actorName string, abilityMod int) (bool, string) {
	// Reduce duration
	condition.Duration--
	
	// Check if condition has expired
	if condition.Duration <= 0 {
		return true, fmt.Sprintf("%s is no longer affected by %s", actorName, condition.Name)
	}
	
	// Process saving throws if needed
	if condition.SaveDC > 0 {
		// Roll saving throw
		saveRoll := c.diceRoller.Roll(1, 20) + abilityMod
		savePassed := saveRoll >= condition.SaveDC
		
		if savePassed {
			return true, fmt.Sprintf("%s succeeds on a %s saving throw and is no longer affected by %s", 
				actorName, condition.SaveType, condition.Name)
		} else {
			return false, fmt.Sprintf("%s fails on a %s saving throw and remains affected by %s (duration: %d rounds)", 
				actorName, condition.SaveType, condition.Name, condition.Duration)
		}
	}
	
	// No saving throw needed, just update duration
	return false, fmt.Sprintf("%s remains affected by %s (duration: %d rounds)", 
		actorName, condition.Name, condition.Duration)
}

// Helper functions

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
