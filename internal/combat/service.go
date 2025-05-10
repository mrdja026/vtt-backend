package combat

import (
        "errors"
        "fmt"
        "time"

        "dnd-combat/internal/models"
        "dnd-combat/pkg/dnd5e"
)

// Service handles combat business logic
type Service struct {
        repo        *Repository
        diceRoller  *dnd5e.DiceRoller
        combatRules *dnd5e.CombatRules
}

// NewService creates a new combat service
func NewService(repo *Repository, diceRoller *dnd5e.DiceRoller, combatRules *dnd5e.CombatRules) *Service {
        return &Service{
                repo:        repo,
                diceRoller:  diceRoller,
                combatRules: combatRules,
        }
}

// CreateCombat initializes a new combat session
func (s *Service) CreateCombat(characters []*models.Character, monsters []*models.Monster, environment string, dmUserID string) (*models.Combat, error) {
        // Create participants from characters and monsters
        participants := make([]*models.Combatant, 0, len(characters)+len(monsters))
        
        // Add characters as combatants
        for _, char := range characters {
                participants = append(participants, &models.Combatant{
                        ID:          char.ID,
                        Type:        "character",
                        Name:        char.Name,
                        Initiative:  0, // Will be calculated later
                        HP:          char.HitPoints,
                        MaxHP:       char.MaxHitPoints,
                        AC:          char.ArmorClass,
                        Position:    [2]int{0, 0}, // Default position, will be updated later
                        UserID:      char.UserID,
                        CharacterID: char.ID,
                        Stats:       char, // Store character data for reference
                        Conditions:  []string{},
                })
        }
        
        // Add monsters as combatants
        for i, monster := range monsters {
                monsterID := fmt.Sprintf("monster_%s_%d", monster.Index, i)
                hp := s.diceRoller.RollHitPoints(monster.HitDice)
                
                participants = append(participants, &models.Combatant{
                        ID:          monsterID,
                        Type:        "monster",
                        Name:        monster.Name,
                        Initiative:  0, // Will be calculated later
                        HP:          hp,
                        MaxHP:       hp,
                        AC:          monster.ArmorClass,
                        Position:    [2]int{0, 0}, // Default position, will be updated later
                        UserID:      dmUserID,    // DM controls all monsters
                        MonsterID:   monster.Index,
                        Stats:       monster, // Store monster data for reference
                        Conditions:  []string{},
                })
        }
        
        // Roll initiative for all participants
        initiative := s.rollInitiative(participants)
        
        // Create combat session
        // Convert []*models.Combatant to []models.Combatant
        combatantSlice := make([]models.Combatant, len(participants))
        for i, p := range participants {
                combatantSlice[i] = *p
        }
        
        // Get battlefield
        battlefield := s.createBattlefield(environment, participants)
        
        combat := &models.Combat{
                DMUserID:         dmUserID,
                Initiative:       initiative,
                Participants:     combatantSlice,
                CurrentTurnIndex: 0,
                RoundNumber:      1,
                Status:           "active",
                Environment:      environment,
                Battlefield:      *battlefield,
                CreatedAt:        time.Now(),
                UpdatedAt:        time.Now(),
        }
        
        // Position participants on the battlefield
        s.positionParticipants(combat)
        
        // Save to database
        if err := s.repo.Create(combat); err != nil {
                return nil, err
        }
        
        return combat, nil
}

// GetCombat retrieves a combat session by ID
func (s *Service) GetCombat(id string) (*models.Combat, error) {
        return s.repo.GetByID(id)
}

// IsUserInCombat checks if a user is participating in a combat
func (s *Service) IsUserInCombat(combat *models.Combat, userID string) bool {
        // Check if user is the DM
        if combat.DMUserID == userID {
                return true
        }
        
        // Check if user has a character in combat
        for _, participant := range combat.Participants {
                if participant.UserID == userID {
                        return true
                }
        }
        
        return false
}

// IsActorsTurn checks if it's the actor's turn
func (s *Service) IsActorsTurn(combat *models.Combat, actorID string) bool {
        if combat.CurrentTurnIndex < 0 || combat.CurrentTurnIndex >= len(combat.Initiative) {
                return false
        }
        
        return combat.Initiative[combat.CurrentTurnIndex].ID == actorID
}

// UserControlsActor checks if a user controls a specific actor
func (s *Service) UserControlsActor(combat *models.Combat, userID string, actorID string) bool {
        // DM controls all monsters
        if combat.DMUserID == userID {
                for _, participant := range combat.Participants {
                        if participant.ID == actorID && participant.Type == "monster" {
                                return true
                        }
                }
        }
        
        // Check if user owns the character
        for _, participant := range combat.Participants {
                if participant.ID == actorID && participant.UserID == userID {
                        return true
                }
        }
        
        return false
}

// ExecuteAction processes a combat action and returns the result
func (s *Service) ExecuteAction(combat *models.Combat, action *models.CombatAction) (*models.ActionResult, error) {
        // Validate the action
        if err := s.validateAction(combat, action); err != nil {
                return nil, err
        }
        
        // Get the actor
        actor := s.getCombatant(combat, action.ActorID)
        if actor == nil {
                return nil, errors.New("actor not found")
        }
        
        // Process different action types
        var result *models.ActionResult
        var err error
        
        switch action.Type {
        case "attack":
                result, err = s.processAttack(combat, action, actor)
        case "cast_spell":
                result, err = s.processSpellCast(combat, action, actor)
        case "move":
                result, err = s.processMovement(combat, action, actor)
        case "dodge":
                result, err = s.processDodge(combat, action, actor)
        case "help":
                result, err = s.processHelp(combat, action, actor)
        case "hide":
                result, err = s.processHide(combat, action, actor)
        case "disengage":
                result, err = s.processDisengage(combat, action, actor)
        case "dash":
                result, err = s.processDash(combat, action, actor)
        case "use_item":
                result, err = s.processItemUse(combat, action, actor)
        default:
                return nil, fmt.Errorf("unknown action type: %s", action.Type)
        }
        
        if err != nil {
                return nil, err
        }
        
        // Update combat with the action result
        if err := s.applyActionResult(combat, result); err != nil {
                return nil, err
        }
        
        // Save action to database
        action.ResultDescription = result.Description
        if err := s.repo.SaveAction(action); err != nil {
                return nil, err
        }
        
        // Update combat session in database
        if err := s.repo.Update(combat); err != nil {
                return nil, err
        }
        
        return result, nil
}

// EndTurn advances to the next participant's turn
func (s *Service) EndTurn(combat *models.Combat) error {
        // Move to the next participant in initiative order
        combat.CurrentTurnIndex++
        
        // If we've gone through everyone, start a new round
        if combat.CurrentTurnIndex >= len(combat.Initiative) {
                combat.CurrentTurnIndex = 0
                combat.RoundNumber++
                
                // Process end-of-round effects (like saving throws against conditions)
                s.processEndOfRound(combat)
        }
        
        // Update in database
        return s.repo.Update(combat)
}

// Helper methods

// rollInitiative calculates initiative order for all participants
func (s *Service) rollInitiative(participants []*models.Combatant) []models.InitiativeItem {
        // Calculate initiative scores
        for _, participant := range participants {
                // Get the dexterity modifier
                var dexMod int
                if participant.Type == "character" {
                        // For characters, calculate from dexterity
                        char := participant.Stats.(*models.Character)
                        dexMod = (char.Dexterity - 10) / 2
                } else if participant.Type == "monster" {
                        // For monsters, use their initiative modifier
                        monster := participant.Stats.(*models.Monster)
                        dexMod = monster.DexterityMod
                }
                
                // Roll initiative: 1d20 + DEX modifier
                initiative := s.diceRoller.Roll(1, 20) + dexMod
                participant.Initiative = initiative
        }
        
        // Sort participants by initiative (descending)
        // Using a simple bubble sort for clarity
        n := len(participants)
        for i := 0; i < n-1; i++ {
                for j := 0; j < n-i-1; j++ {
                        if participants[j].Initiative < participants[j+1].Initiative {
                                participants[j], participants[j+1] = participants[j+1], participants[j]
                        }
                }
        }
        
        // Create the initiative order list
        initiativeItems := make([]models.InitiativeItem, len(participants))
        for i, participant := range participants {
                initiativeItems[i] = models.InitiativeItem{
                        ID:        participant.ID,
                        Name:      participant.Name,
                        Initiative: participant.Initiative,
                        Dexterity: 0, // This is a placeholder, should be filled from the character/monster stats
                        IsCharacter: participant.Type == "character",
                }
        }
        
        return initiativeItems
}

// createBattlefield initializes the combat battlefield
func (s *Service) createBattlefield(environment string, participants []*models.Combatant) *models.Battlefield {
        // Create a standard 10x10 grid battlefield
        // Size can be adjusted based on environment and number of participants
        
        // Initialize maps
        grid := make(map[string]string)
        terrain := make(map[string]string)
        obstacles := make(map[string]bool)
        
        // Create the battlefield
        battlefield := &models.Battlefield{
                Width:     10,
                Height:    10,
                Grid:      grid,
                Terrain:   terrain,
                Obstacles: obstacles,
        }
        
        // Initialize the grid
        for i := 0; i < 10; i++ {
                for j := 0; j < 10; j++ {
                        // Convert coordinates to string keys
                        key := fmt.Sprintf("%d,%d", i, j)
                        // Set default values
                        battlefield.Terrain[key] = "normal"
                        battlefield.Grid[key] = "empty"
                }
        }
        
        // Customize based on environment (simple example)
        switch environment {
        case "forest":
                // Add some trees as obstacles
                treePositions := [][2]int{{1, 1}, {3, 4}, {6, 7}, {8, 2}}
                for _, pos := range treePositions {
                        key := fmt.Sprintf("%d,%d", pos[0], pos[1])
                        battlefield.Grid[key] = "tree"
                        battlefield.Obstacles[key] = true
                }
                
                // Add difficult terrain
                difficultTerrainPositions := [][2]int{{2, 2}, {2, 3}, {3, 2}}
                for _, pos := range difficultTerrainPositions {
                        key := fmt.Sprintf("%d,%d", pos[0], pos[1])
                        battlefield.Terrain[key] = "difficult"
                }
                // Add one more difficult terrain
                battlefield.Terrain[fmt.Sprintf("%d,%d", 3, 3)] = "difficult"
        case "dungeon":
                // Add some walls as obstacles
                wallPositions := [][2]int{{0, 5}, {1, 5}, {2, 5}, {3, 5}}
                for _, pos := range wallPositions {
                        key := fmt.Sprintf("%d,%d", pos[0], pos[1])
                        battlefield.Grid[key] = "wall"
                        battlefield.Obstacles[key] = true
                }
                
                // Add a trap
                trapKey := fmt.Sprintf("%d,%d", 5, 5)
                battlefield.Terrain[trapKey] = "trap"
        case "cave":
                // Add some rocks as obstacles
                rockPositions := [][2]int{{2, 3}, {7, 6}}
                for _, pos := range rockPositions {
                        key := fmt.Sprintf("%d,%d", pos[0], pos[1])
                        battlefield.Grid[key] = "rock"
                        battlefield.Obstacles[key] = true
                }
                
                // Add some water terrain
                waterPositions := [][2]int{{4, 4}, {4, 5}, {5, 4}, {5, 5}}
                for _, pos := range waterPositions {
                        key := fmt.Sprintf("%d,%d", pos[0], pos[1])
                        battlefield.Terrain[key] = "water"
                }
        }
        
        return battlefield
}

// positionParticipants places participants on the battlefield
func (s *Service) positionParticipants(combat *models.Combat) {
        // Simple positioning algorithm: characters on one side, monsters on the other
        characterCount := 0
        monsterCount := 0
        
        for _, participant := range combat.Participants {
                if participant.Type == "character" {
                        // Position characters on the left side
                        x := 2
                        y := 2 + characterCount
                        
                        // Adjust to stay within bounds
                        if y >= combat.Battlefield.Height {
                                y = characterCount % combat.Battlefield.Height
                                x = 1
                        }
                        
                        participant.Position = [2]int{x, y}
                        characterCount++
                } else {
                        // Position monsters on the right side
                        x := combat.Battlefield.Width - 3
                        y := 2 + monsterCount
                        
                        // Adjust to stay within bounds
                        if y >= combat.Battlefield.Height {
                                y = monsterCount % combat.Battlefield.Height
                                x = combat.Battlefield.Width - 2
                        }
                        
                        participant.Position = [2]int{x, y}
                        monsterCount++
                }
        }
}

// getCombatant finds a combatant by ID
func (s *Service) getCombatant(combat *models.Combat, id string) *models.Combatant {
        for i, participant := range combat.Participants {
                if participant.ID == id {
                        // Need to return a pointer to the combatant, not a copy
                        return &combat.Participants[i]
                }
        }
        return nil
}

// validateAction checks if an action is valid
func (s *Service) validateAction(combat *models.Combat, action *models.CombatAction) error {
        // Check if action is allowed in the current state
        if combat.Status != "active" {
                return errors.New("combat is not active")
        }
        
        // Check if it's the actor's turn
        if !s.IsActorsTurn(combat, action.ActorID) {
                return errors.New("it's not this actor's turn")
        }
        
        // Actor must exist in combat
        actor := s.getCombatant(combat, action.ActorID)
        if actor == nil {
                return errors.New("actor not found in combat")
        }
        
        // Dead actors can't take actions
        if actor.HP <= 0 {
                return errors.New("actor is unconscious or dead")
        }
        
        // Validate targets if provided
        if len(action.TargetIDs) > 0 {
                for _, targetID := range action.TargetIDs {
                        target := s.getCombatant(combat, targetID)
                        if target == nil {
                                return errors.New("target not found in combat")
                        }
                }
        }
        
        // Validate action-specific requirements
        switch action.Type {
        case "attack":
                return s.validateAttack(combat, action, actor)
        case "cast_spell":
                return s.validateSpellCast(combat, action, actor)
        case "move":
                return s.validateMovement(combat, action, actor)
        }
        
        return nil
}

// validateAttack checks if an attack action is valid
func (s *Service) validateAttack(combat *models.Combat, action *models.CombatAction, actor *models.Combatant) error {
        // Must have a target
        if len(action.TargetIDs) == 0 {
                return errors.New("attack requires a target")
        }
        
        // Must specify a weapon
        if action.WeaponName == "" {
                return errors.New("attack requires a weapon")
        }
        
        // Check if target is in range
        target := s.getCombatant(combat, action.TargetIDs[0])
        
        // Get weapon range
        var weaponRange int
        if actor.Type == "character" {
                // For characters, check equipment
                // Simple implementation - would need to be expanded to check actual equipment
                weaponRange = 5 // Melee range by default
                if action.WeaponName == "longbow" || action.WeaponName == "shortbow" {
                        weaponRange = 80 // Ranged weapons
                }
        } else {
                // For monsters, use action range from stats
                monster := actor.Stats.(*models.Monster)
                // Simple implementation - check monster actions
                weaponRange = 5 // Melee range by default
                for _, monsterAction := range monster.Actions {
                        if monsterAction.Name == action.WeaponName {
                                if monsterAction.Range > 0 {
                                        weaponRange = monsterAction.Range
                                }
                                break
                        }
                }
        }
        
        // Calculate distance
        distance := calculateDistance(actor.Position, target.Position)
        if distance > weaponRange {
                return fmt.Errorf("target is out of range (distance: %d, range: %d)", distance, weaponRange)
        }
        
        return nil
}

// validateSpellCast checks if a spell casting action is valid
func (s *Service) validateSpellCast(combat *models.Combat, action *models.CombatAction, actor *models.Combatant) error {
        // Must have a spell
        if action.SpellID == "" {
                return errors.New("spell casting requires a spell")
        }
        
        // Character-specific validation
        if actor.Type == "character" {
                char := actor.Stats.(*models.Character)
                
                // Check if character knows the spell
                spellFound := false
                for _, spell := range char.Spells {
                        if spell == action.SpellID {
                                spellFound = true
                                break
                        }
                }
                
                if !spellFound {
                        return errors.New("character doesn't know this spell")
                }
                
                // Would need more validation for spell slots, components, etc.
        }
        
        // Target validation depends on the spell
        // Would need to lookup spell targeting requirements from SRD
        
        return nil
}

// validateMovement checks if a movement action is valid
func (s *Service) validateMovement(combat *models.Combat, action *models.CombatAction, actor *models.Combatant) error {
        // Must have a movement path
        if len(action.MovementPath) == 0 {
                return errors.New("movement requires a path")
        }
        
        // Calculate movement speed
        var speed int
        if actor.Type == "character" {
                // Simple implementation - base speed
                speed = 30
        } else {
                monster := actor.Stats.(*models.Monster)
                speed = monster.Speed.Walk
        }
        
        // Convert to grid movement (assuming 5ft per grid)
        maxSquares := speed / 5
        
        // Check if path is too long
        if len(action.MovementPath) > maxSquares {
                return fmt.Errorf("movement path exceeds speed (max: %d squares)", maxSquares)
        }
        
        // Check for valid path
        currentPos := actor.Position
        for _, pos := range action.MovementPath {
                // Check if position is in bounds
                if pos[0] < 0 || pos[0] >= combat.Battlefield.Width || 
                   pos[1] < 0 || pos[1] >= combat.Battlefield.Height {
                        return errors.New("movement path goes out of bounds")
                }
                
                // Check if position is blocked
                posKey := fmt.Sprintf("%d,%d", pos[0], pos[1])
                if combat.Battlefield.Grid[posKey] == "wall" || 
                   combat.Battlefield.Grid[posKey] == "tree" ||
                   combat.Battlefield.Grid[posKey] == "rock" ||
                   combat.Battlefield.Obstacles[posKey] {
                        return errors.New("movement path is blocked by an obstacle")
                }
                
                // Check for other combatants
                for _, other := range combat.Participants {
                        if other.ID != actor.ID && other.Position[0] == pos[0] && other.Position[1] == pos[1] {
                                return errors.New("movement path is blocked by another combatant")
                        }
                }
                
                // Check for valid movement (no diagonal jumps)
                dx := abs(pos[0] - currentPos[0])
                dy := abs(pos[1] - currentPos[1])
                if dx > 1 || dy > 1 || (dx + dy) > 1 {
                        return errors.New("invalid movement: can only move to adjacent squares")
                }
                
                currentPos = pos
        }
        
        return nil
}

// processAttack handles an attack action
func (s *Service) processAttack(combat *models.Combat, action *models.CombatAction, actor *models.Combatant) (*models.ActionResult, error) {
        target := s.getCombatant(combat, action.TargetIDs[0])
        
        result := &models.ActionResult{
                Success:     false,
                Damage:      0,
                Description: "",
        }
        
        // Roll to hit
        var attackBonus int
        var damage int
        var damageType string
        
        if actor.Type == "character" {
                char := actor.Stats.(*models.Character)
                
                // Calculate attack bonus based on weapon and stats
                // Simple implementation
                if action.WeaponName == "longbow" || action.WeaponName == "shortbow" {
                        // Ranged weapons use DEX
                        attackBonus = (char.Dexterity - 10) / 2
                } else {
                        // Melee weapons use STR by default
                        attackBonus = (char.Strength - 10) / 2
                }
                
                // Add proficiency bonus based on level
                profBonus := 2
                if char.Level >= 5 {
                        profBonus = 3
                }
                if char.Level >= 9 {
                        profBonus = 4
                }
                if char.Level >= 13 {
                        profBonus = 5
                }
                if char.Level >= 17 {
                        profBonus = 6
                }
                
                attackBonus += profBonus
                
                // Determine damage based on weapon
                switch action.WeaponName {
                case "longsword":
                        damage = s.diceRoller.Roll(1, 8) + (char.Strength - 10) / 2
                        damageType = "slashing"
                case "longbow":
                        damage = s.diceRoller.Roll(1, 8) + (char.Dexterity - 10) / 2
                        damageType = "piercing"
                case "dagger":
                        damage = s.diceRoller.Roll(1, 4) + (char.Strength - 10) / 2
                        damageType = "piercing"
                default:
                        damage = s.diceRoller.Roll(1, 6) + (char.Strength - 10) / 2
                        damageType = "bludgeoning"
                }
        } else {
                // Monster attack
                monster := actor.Stats.(*models.Monster)
                
                // Find the action for this weapon
                for _, monsterAction := range monster.Actions {
                        if monsterAction.Name == action.WeaponName {
                                attackBonus = monsterAction.AttackBonus
                                
                                // Parse damage (simplified)
                                damage = s.diceRoller.Roll(monsterAction.Damage.DiceCount, monsterAction.Damage.DiceValue) + monsterAction.Damage.Bonus
                                damageType = monsterAction.Damage.Type
                                break
                        }
                }
        }
        
        // Roll attack
        attackRoll := s.diceRoller.Roll(1, 20)
        totalAttack := attackRoll + attackBonus
        
        // Check for critical hit or miss
        isCritical := attackRoll == 20
        isCritMiss := attackRoll == 1
        
        if isCritMiss {
                result.Description = fmt.Sprintf("%s critically misses their attack with %s against %s!", 
                        actor.Name, action.WeaponName, target.Name)
                return result, nil
        }
        
        if isCritical {
                // Double damage dice on critical hit
                if actor.Type == "character" {
                        // Roll extra damage dice based on weapon
                        switch action.WeaponName {
                        case "longsword":
                                damage += s.diceRoller.Roll(1, 8)
                        case "longbow":
                                damage += s.diceRoller.Roll(1, 8)
                        case "dagger":
                                damage += s.diceRoller.Roll(1, 4)
                        default:
                                damage += s.diceRoller.Roll(1, 6)
                        }
                } else {
                        // Monster critical hit
                        monster := actor.Stats.(*models.Monster)
                        
                        for _, monsterAction := range monster.Actions {
                                if monsterAction.Name == action.WeaponName {
                                        damage += s.diceRoller.Roll(monsterAction.Damage.DiceCount, monsterAction.Damage.DiceValue)
                                        break
                                }
                        }
                }
        }
        
        // Check if attack hits
        hits := isCritical || totalAttack >= target.AC
        
        if !hits {
                result.Description = fmt.Sprintf("%s attacks %s with %s but misses! (Rolled %d + %d = %d vs AC %d)", 
                        actor.Name, target.Name, action.WeaponName, attackRoll, attackBonus, totalAttack, target.AC)
                return result, nil
        }
        
        // Apply damage
        result.Success = true
        result.Damage = damage
        
        // Update target HP
        newHP := target.HP - damage
        if newHP < 0 {
                newHP = 0
        }
        target.HP = newHP
        
        // Check if target is defeated
        if target.HP == 0 {
                if target.Type == "character" {
                        // Characters are unconscious at 0 HP
                        target.Conditions = append(target.Conditions, "unconscious")
                        hitDescription := "hits"
                        if isCritical {
                                hitDescription = "critically hits"
                        }
                        result.Description = fmt.Sprintf("%s %s %s with %s for %d %s damage! %s falls unconscious!", 
                                actor.Name, 
                                hitDescription,
                                target.Name, 
                                action.WeaponName, 
                                damage,
                                damageType,
                                target.Name)
                } else {
                        // Monsters are dead at 0 HP
                        hitDescription := "hits"
                        if isCritical {
                                hitDescription = "critically hits"
                        }
                        result.Description = fmt.Sprintf("%s %s %s with %s for %d %s damage! %s is defeated!", 
                                actor.Name, 
                                hitDescription,
                                target.Name, 
                                action.WeaponName, 
                                damage,
                                damageType,
                                target.Name)
                }
        } else {
                hitDescription := "hits"
                if isCritical {
                        hitDescription = "critically hits"
                }
                result.Description = fmt.Sprintf("%s %s %s with %s for %d %s damage! (HP: %d/%d)", 
                        actor.Name, 
                        hitDescription,
                        target.Name, 
                        action.WeaponName, 
                        damage,
                        damageType,
                        target.HP,
                        target.MaxHP)
        }
        
        return result, nil
}

// processSpellCast handles a spell casting action
func (s *Service) processSpellCast(combat *models.Combat, action *models.CombatAction, actor *models.Combatant) (*models.ActionResult, error) {
        result := &models.ActionResult{
                Success:     true,
                Damage:      0,
                Description: "",
        }
        
        // This would typically involve looking up the spell in the SRD API
        // For simplicity, we'll implement a few basic spells
        
        switch action.SpellID {
        case "cure-wounds":
                // Healing spell
                if len(action.TargetIDs) == 0 {
                        return nil, errors.New("cure wounds requires a target")
                }
                
                target := s.getCombatant(combat, action.TargetIDs[0])
                
                // Calculate healing amount
                var spellcastingMod int
                if actor.Type == "character" {
                        char := actor.Stats.(*models.Character)
                        // Use the appropriate spellcasting ability modifier
                        // This is a simplification - different classes use different abilities
                        spellcastingMod = (char.Wisdom - 10) / 2
                } else {
                        spellcastingMod = 3 // Default for monsters
                }
                
                // 1d8 + spellcasting modifier
                healing := s.diceRoller.Roll(1, 8) + spellcastingMod
                
                // Apply healing
                newHP := target.HP + healing
                if newHP > target.MaxHP {
                        newHP = target.MaxHP
                }
                
                oldHP := target.HP
                target.HP = newHP
                
                // Remove unconscious condition if healed from 0 HP
                if oldHP == 0 && newHP > 0 {
                        for i, condition := range target.Conditions {
                                if condition == "unconscious" {
                                        // Remove the condition
                                        target.Conditions = append(target.Conditions[:i], target.Conditions[i+1:]...)
                                        break
                                }
                        }
                }
                
                result.Description = fmt.Sprintf("%s casts Cure Wounds on %s, healing %d damage! (HP: %d/%d)",
                        actor.Name, target.Name, healing, target.HP, target.MaxHP)
                
        case "magic-missile":
                // Damage spell that always hits
                if len(action.TargetIDs) == 0 {
                        return nil, errors.New("magic missile requires a target")
                }
                
                target := s.getCombatant(combat, action.TargetIDs[0])
                
                // 3 missiles, each dealing 1d4+1 force damage
                totalDamage := 0
                for i := 0; i < 3; i++ {
                        damage := s.diceRoller.Roll(1, 4) + 1
                        totalDamage += damage
                }
                
                // Apply damage
                newHP := target.HP - totalDamage
                if newHP < 0 {
                        newHP = 0
                }
                target.HP = newHP
                
                result.Damage = totalDamage
                
                // Check if target is defeated
                if target.HP == 0 {
                        if target.Type == "character" {
                                // Characters are unconscious at 0 HP
                                target.Conditions = append(target.Conditions, "unconscious")
                                result.Description = fmt.Sprintf("%s casts Magic Missile at %s, dealing %d force damage! %s falls unconscious!", 
                                        actor.Name, target.Name, totalDamage, target.Name)
                        } else {
                                // Monsters are dead at 0 HP
                                result.Description = fmt.Sprintf("%s casts Magic Missile at %s, dealing %d force damage! %s is defeated!", 
                                        actor.Name, target.Name, totalDamage, target.Name)
                        }
                } else {
                        result.Description = fmt.Sprintf("%s casts Magic Missile at %s, dealing %d force damage! (HP: %d/%d)", 
                                actor.Name, target.Name, totalDamage, target.HP, target.MaxHP)
                }
                
        case "shield":
                // Defensive spell
                result.Description = fmt.Sprintf("%s casts Shield, granting +5 AC until the start of their next turn!",
                        actor.Name)
                
                // Add temporary condition
                actor.Conditions = append(actor.Conditions, "shield")
                actor.AC += 5 // Temporary AC boost
                
        default:
                return nil, fmt.Errorf("spell '%s' not implemented", action.SpellID)
        }
        
        return result, nil
}

// processMovement handles a movement action
func (s *Service) processMovement(combat *models.Combat, action *models.CombatAction, actor *models.Combatant) (*models.ActionResult, error) {
        // Update actor position to the end of the path
        if len(action.MovementPath) > 0 {
                oldPos := actor.Position
                actor.Position = action.MovementPath[len(action.MovementPath)-1]
                
                return &models.ActionResult{
                        Success:     true,
                        Description: fmt.Sprintf("%s moves from [%d,%d] to [%d,%d]", 
                                actor.Name, oldPos[0], oldPos[1], actor.Position[0], actor.Position[1]),
                }, nil
        }
        
        return &models.ActionResult{
                Success:     false,
                Description: fmt.Sprintf("%s attempts to move but stays in place", actor.Name),
        }, nil
}

// processDodge handles a dodge action
func (s *Service) processDodge(combat *models.Combat, action *models.CombatAction, actor *models.Combatant) (*models.ActionResult, error) {
        // Add dodge condition
        actor.Conditions = append(actor.Conditions, "dodge")
        
        return &models.ActionResult{
                Success:     true,
                Description: fmt.Sprintf("%s takes the Dodge action, giving attackers disadvantage until their next turn", actor.Name),
        }, nil
}

// processHelp handles a help action
func (s *Service) processHelp(combat *models.Combat, action *models.CombatAction, actor *models.Combatant) (*models.ActionResult, error) {
        if len(action.TargetIDs) == 0 {
                return nil, errors.New("help action requires a target")
        }
        
        target := s.getCombatant(combat, action.TargetIDs[0])
        
        // Add help condition to the target
        target.Conditions = append(target.Conditions, "helped")
        
        return &models.ActionResult{
                Success:     true,
                Description: fmt.Sprintf("%s helps %s, giving them advantage on their next ability check or attack roll", 
                        actor.Name, target.Name),
        }, nil
}

// processHide handles a hide action
func (s *Service) processHide(combat *models.Combat, action *models.CombatAction, actor *models.Combatant) (*models.ActionResult, error) {
        // Roll stealth check
        var stealthMod int
        if actor.Type == "character" {
                char := actor.Stats.(*models.Character)
                stealthMod = (char.Dexterity - 10) / 2
        } else {
                monster := actor.Stats.(*models.Monster)
                stealthMod = monster.DexterityMod
        }
        
        stealthRoll := s.diceRoller.Roll(1, 20) + stealthMod
        
        // Add hidden condition with the stealth value
        actor.Conditions = append(actor.Conditions, fmt.Sprintf("hidden:%d", stealthRoll))
        
        return &models.ActionResult{
                Success:     true,
                Description: fmt.Sprintf("%s attempts to hide, rolling a %d for Stealth", actor.Name, stealthRoll),
        }, nil
}

// processDisengage handles a disengage action
func (s *Service) processDisengage(combat *models.Combat, action *models.CombatAction, actor *models.Combatant) (*models.ActionResult, error) {
        // Add disengage condition
        actor.Conditions = append(actor.Conditions, "disengage")
        
        return &models.ActionResult{
                Success:     true,
                Description: fmt.Sprintf("%s takes the Disengage action, preventing opportunity attacks until their next turn", actor.Name),
        }, nil
}

// processDash handles a dash action
func (s *Service) processDash(combat *models.Combat, action *models.CombatAction, actor *models.Combatant) (*models.ActionResult, error) {
        // Add dash condition
        actor.Conditions = append(actor.Conditions, "dash")
        
        return &models.ActionResult{
                Success:     true,
                Description: fmt.Sprintf("%s takes the Dash action, doubling their movement speed until the end of their turn", actor.Name),
        }, nil
}

// processItemUse handles an item use action
func (s *Service) processItemUse(combat *models.Combat, action *models.CombatAction, actor *models.Combatant) (*models.ActionResult, error) {
        // Extract item name from extra data
        itemName, ok := action.ExtraData["item_name"].(string)
        if !ok {
                return nil, errors.New("item name not provided or invalid")
        }
        
        var result *models.ActionResult
        
        // Process different item types
        switch itemName {
        case "healing-potion":
                // Heal 2d4+2 hit points
                healing := s.diceRoller.Roll(2, 4) + 2
                
                // Apply healing
                newHP := actor.HP + healing
                if newHP > actor.MaxHP {
                        newHP = actor.MaxHP
                }
                actor.HP = newHP
                
                result = &models.ActionResult{
                        Success:     true,
                        Description: fmt.Sprintf("%s drinks a Healing Potion, recovering %d hit points! (HP: %d/%d)",
                                actor.Name, healing, actor.HP, actor.MaxHP),
                }
                
        case "antitoxin":
                // Gain advantage on saving throws against poison
                actor.Conditions = append(actor.Conditions, "antitoxin")
                
                result = &models.ActionResult{
                        Success:     true,
                        Description: fmt.Sprintf("%s uses Antitoxin, gaining advantage on saving throws against poison for 1 hour",
                                actor.Name),
                }
                
        default:
                return nil, fmt.Errorf("item '%s' not implemented", itemName)
        }
        
        return result, nil
}

// applyActionResult updates the combat state based on action results
func (s *Service) applyActionResult(combat *models.Combat, result *models.ActionResult) error {
        // Check if any participants are defeated
        allMonstersDead := true
        allPlayersDead := true
        
        for _, participant := range combat.Participants {
                if participant.HP <= 0 {
                        if participant.Type == "monster" {
                                // Remove dead monsters from the battlefield
                                participant.Position = [2]int{-1, -1}
                        }
                } else {
                        if participant.Type == "monster" {
                                allMonstersDead = false
                        } else if participant.Type == "character" {
                                allPlayersDead = false
                        }
                }
        }
        
        // Check for combat end conditions
        if allMonstersDead {
                combat.Status = "victory"
                return nil
        }
        
        if allPlayersDead {
                combat.Status = "defeat"
                return nil
        }
        
        return nil
}

// processEndOfRound handles end-of-round effects
func (s *Service) processEndOfRound(combat *models.Combat) {
        for _, participant := range combat.Participants {
                // Process ongoing conditions
                newConditions := []string{}
                
                for _, condition := range participant.Conditions {
                        // Check for temporary conditions that should expire
                        switch condition {
                        case "shield", "dodge", "disengage", "dash", "helped":
                                // These conditions expire at the end of the round
                                continue
                        default:
                                // Keep other conditions
                                newConditions = append(newConditions, condition)
                        }
                }
                
                // Update conditions
                participant.Conditions = newConditions
                
                // Reset AC if shield condition was removed
                if containsString(participant.Conditions, "shield") && !containsString(newConditions, "shield") {
                        participant.AC -= 5 // Remove Shield spell AC bonus
                }
                
                // Process ongoing damage (not implemented in this simple version)
        }
}

// Helper functions

// calculateDistance calculates the distance between two positions on the grid
func calculateDistance(pos1, pos2 [2]int) int {
        // Manhattan distance for grid movement
        return abs(pos1[0]-pos2[0]) + abs(pos1[1]-pos2[1])
}

// abs returns the absolute value of an integer
func abs(x int) int {
        if x < 0 {
                return -x
        }
        return x
}

// containsString checks if a string slice contains a string
func containsString(slice []string, str string) bool {
        for _, item := range slice {
                if item == str {
                        return true
                }
        }
        return false
}
