package dnd5e

import (
        "fmt"
        "math/rand"
        "regexp"
        "strconv"
        "strings"
        "time"
)

// DiceRoller handles dice rolling operations for D&D
type DiceRoller struct {
        rng *rand.Rand
}

// NewDiceRoller creates a new dice roller with a seeded random source
func NewDiceRoller() *DiceRoller {
        return &DiceRoller{
                rng: rand.New(rand.NewSource(time.Now().UnixNano())),
        }
}

// Roll rolls a specified number of dice with the given sides
func (d *DiceRoller) Roll(count, sides int) int {
        if count <= 0 || sides <= 0 {
                return 0
        }

        total := 0
        for i := 0; i < count; i++ {
                total += d.rng.Intn(sides) + 1
        }
        return total
}

// RollWithAdvantage rolls a d20 with advantage (roll twice, take the higher value)
func (d *DiceRoller) RollWithAdvantage() int {
        roll1 := d.Roll(1, 20)
        roll2 := d.Roll(1, 20)
        if roll1 > roll2 {
                return roll1
        }
        return roll2
}

// RollWithDisadvantage rolls a d20 with disadvantage (roll twice, take the lower value)
func (d *DiceRoller) RollWithDisadvantage() int {
        roll1 := d.Roll(1, 20)
        roll2 := d.Roll(1, 20)
        if roll1 < roll2 {
                return roll1
        }
        return roll2
}

// RollHitPoints calculates hit points based on a hit dice string (e.g., "3d8+4")
func (d *DiceRoller) RollHitPoints(hitDice string) int {
        // Parse hit dice string
        re := regexp.MustCompile(`(\d+)d(\d+)(?:\+(\d+))?`)
        match := re.FindStringSubmatch(hitDice)
        
        if len(match) < 3 {
                // Invalid format, return default value
                return 10
        }
        
        count, _ := strconv.Atoi(match[1])
        sides, _ := strconv.Atoi(match[2])
        
        bonus := 0
        if len(match) > 3 && match[3] != "" {
                bonus, _ = strconv.Atoi(match[3])
        }
        
        // Roll the dice and add the bonus
        return d.Roll(count, sides) + bonus
}

// RollDamage calculates damage based on a damage formula (e.g., "2d6+3")
func (d *DiceRoller) RollDamage(damageFormula string) int {
        return d.parseAndRoll(damageFormula)
}

// parseAndRoll parses a dice formula and rolls it
func (d *DiceRoller) parseAndRoll(formula string) int {
        formula = strings.TrimSpace(formula)
        
        // Handle empty formula
        if formula == "" {
                return 0
        }
        
        // Split by addition
        parts := strings.Split(formula, "+")
        total := 0
        
        for _, part := range parts {
                part = strings.TrimSpace(part)
                
                // Check if it's a dice roll (e.g., "2d6")
                if strings.Contains(part, "d") {
                        diceParts := strings.Split(part, "d")
                        if len(diceParts) != 2 {
                                continue
                        }
                        
                        count, err := strconv.Atoi(diceParts[0])
                        if err != nil {
                                count = 1 // Default to 1 if not specified
                        }
                        
                        sides, err := strconv.Atoi(diceParts[1])
                        if err != nil {
                                continue
                        }
                        
                        total += d.Roll(count, sides)
                } else {
                        // It's a static bonus
                        bonus, err := strconv.Atoi(part)
                        if err != nil {
                                continue
                        }
                        
                        total += bonus
                }
        }
        
        return total
}

// RollSavingThrow simulates a saving throw against a DC
func (d *DiceRoller) RollSavingThrow(abilityMod int, dc int, hasAdvantage bool, hasDisadvantage bool) bool {
        var roll int
        
        if hasAdvantage && !hasDisadvantage {
                roll = d.RollWithAdvantage()
        } else if hasDisadvantage && !hasAdvantage {
                roll = d.RollWithDisadvantage()
        } else {
                roll = d.Roll(1, 20)
        }
        
        total := roll + abilityMod
        return total >= dc
}

// RollInitiative simulates an initiative roll
func (d *DiceRoller) RollInitiative(dexMod int, hasAdvantage bool) int {
        if hasAdvantage {
                return d.RollWithAdvantage() + dexMod
        }
        return d.Roll(1, 20) + dexMod
}

// RollAttack simulates an attack roll
func (d *DiceRoller) RollAttack(attackBonus int, hasAdvantage bool, hasDisadvantage bool) (int, bool) {
        var roll int
        
        if hasAdvantage && !hasDisadvantage {
                roll = d.RollWithAdvantage()
        } else if hasDisadvantage && !hasAdvantage {
                roll = d.RollWithDisadvantage()
        } else {
                roll = d.Roll(1, 20)
        }
        
        // Check for critical hit
        isCritical := roll == 20
        
        total := roll + attackBonus
        return total, isCritical
}

// RollAbilityCheck simulates an ability check
func (d *DiceRoller) RollAbilityCheck(abilityMod int, profBonus int, isProficient bool, hasAdvantage bool, hasDisadvantage bool) int {
        var roll int
        
        if hasAdvantage && !hasDisadvantage {
                roll = d.RollWithAdvantage()
        } else if hasDisadvantage && !hasAdvantage {
                roll = d.RollWithDisadvantage()
        } else {
                roll = d.Roll(1, 20)
        }
        
        total := roll + abilityMod
        if isProficient {
                total += profBonus
        }
        
        return total
}

// FormatRollResult formats a roll result as a string
func (d *DiceRoller) FormatRollResult(roll int, modifier int) string {
        if modifier == 0 {
                return fmt.Sprintf("%d", roll)
        } else if modifier > 0 {
                return fmt.Sprintf("%d + %d = %d", roll, modifier, roll+modifier)
        } else {
                return fmt.Sprintf("%d - %d = %d", roll, -modifier, roll+modifier)
        }
}
