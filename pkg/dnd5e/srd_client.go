package dnd5e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

const DefaultBaseURL = "https://www.dnd5eapi.co/api"

// Cache interface for SRD API responses
type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
	Delete(key string)
}

// InMemoryCache implements a simple in-memory cache for SRD API responses
type InMemoryCache struct {
	items map[string]cacheItem
	mu    sync.RWMutex
}

type cacheItem struct {
	value      interface{}
	expiration time.Time
}

// NewInMemoryCache creates a new in-memory cache
func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		items: make(map[string]cacheItem),
	}
}

// Get retrieves an item from the cache
func (c *InMemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	// Check if item has expired
	if time.Now().After(item.expiration) {
		delete(c.items, key)
		return nil, false
	}

	return item.value, true
}

// Set adds an item to the cache with a default expiration of 1 hour
func (c *InMemoryCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = cacheItem{
		value:      value,
		expiration: time.Now().Add(1 * time.Hour),
	}
}

// Delete removes an item from the cache
func (c *InMemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

// SRDClient is a client for the D&D 5e SRD API
type SRDClient struct {
	httpClient *http.Client
	baseURL    string
	cache      Cache
}

// NewSRDClient creates a new SRD API client
func NewSRDClient(httpClient *http.Client, cache Cache) *SRDClient {
	baseURL := os.Getenv("SRD_API_BASE_URL")
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	return &SRDClient{
		httpClient: httpClient,
		baseURL:    baseURL,
		cache:      cache,
	}
}

// ClassData represents a character class from the SRD API
type ClassData struct {
	Index       string   `json:"index"`
	Name        string   `json:"name"`
	HitDie      int      `json:"hit_die"`
	ProfChoices []string `json:"prof_choices"`
	Proficiencies []string `json:"proficiencies"`
	SavingThrows []string `json:"saving_throws"`
	StartingEquipment []string `json:"starting_equipment"`
	ClassLevels  string   `json:"class_levels"`
	Subclasses   []struct {
		Name  string `json:"name"`
		URL   string `json:"url"`
	} `json:"subclasses"`
}

// RaceData represents a character race from the SRD API
type RaceData struct {
	Index          string   `json:"index"`
	Name           string   `json:"name"`
	Speed          int      `json:"speed"`
	AbilityBonuses []struct {
		AbilityScore struct {
			Index string `json:"index"`
			Name  string `json:"name"`
		} `json:"ability_score"`
		Bonus int `json:"bonus"`
	} `json:"ability_bonuses"`
	Alignment       string   `json:"alignment"`
	Age             string   `json:"age"`
	Size            string   `json:"size"`
	SizeDescription string   `json:"size_description"`
	Languages       []string `json:"languages"`
	Traits          []string `json:"traits"`
	Subraces        []struct {
		Name  string `json:"name"`
		URL   string `json:"url"`
	} `json:"subraces"`
}

// GetClass fetches a class from the SRD API
func (c *SRDClient) GetClass(index string) (*ClassData, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("class:%s", index)
	if data, found := c.cache.Get(cacheKey); found {
		return data.(*ClassData), nil
	}

	// Fetch from API
	url := fmt.Sprintf("%s/classes/%s", c.baseURL, index)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching class data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-OK status: %d", resp.StatusCode)
	}

	var classData ClassData
	if err := json.NewDecoder(resp.Body).Decode(&classData); err != nil {
		return nil, fmt.Errorf("error decoding class data: %w", err)
	}

	// Store in cache
	c.cache.Set(cacheKey, &classData)
	return &classData, nil
}

// GetRace fetches a race from the SRD API
func (c *SRDClient) GetRace(index string) (*RaceData, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("race:%s", index)
	if data, found := c.cache.Get(cacheKey); found {
		return data.(*RaceData), nil
	}

	// Fetch from API
	url := fmt.Sprintf("%s/races/%s", c.baseURL, index)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching race data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-OK status: %d", resp.StatusCode)
	}

	var raceData RaceData
	if err := json.NewDecoder(resp.Body).Decode(&raceData); err != nil {
		return nil, fmt.Errorf("error decoding race data: %w", err)
	}

	// Store in cache
	c.cache.Set(cacheKey, &raceData)
	return &raceData, nil
}

// GetMonster fetches a monster from the SRD API
func (c *SRDClient) GetMonster(index string) (*Monster, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("monster:%s", index)
	if data, found := c.cache.Get(cacheKey); found {
		return data.(*Monster), nil
	}

	// Fetch from API
	url := fmt.Sprintf("%s/monsters/%s", c.baseURL, index)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching monster data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-OK status: %d", resp.StatusCode)
	}

	var apiResponse struct {
		Index        string  `json:"index"`
		Name         string  `json:"name"`
		Size         string  `json:"size"`
		Type         string  `json:"type"`
		Alignment    string  `json:"alignment"`
		ArmorClass   int     `json:"armor_class"`
		HitPoints    int     `json:"hit_points"`
		HitDice      string  `json:"hit_dice"`
		Speed        struct {
			Walk   string `json:"walk"`
			Swim   string `json:"swim,omitempty"`
			Fly    string `json:"fly,omitempty"`
			Climb  string `json:"climb,omitempty"`
			Burrow string `json:"burrow,omitempty"`
		} `json:"speed"`
		Strength     int     `json:"strength"`
		Dexterity    int     `json:"dexterity"`
		Constitution int     `json:"constitution"`
		Intelligence int     `json:"intelligence"`
		Wisdom       int     `json:"wisdom"`
		Charisma     int     `json:"charisma"`
		Actions      []struct {
			Name        string `json:"name"`
			Description string `json:"desc"`
			AttackBonus int    `json:"attack_bonus,omitempty"`
			Damage      []struct {
				DamageDice string `json:"damage_dice"`
				DamageType struct {
					Name string `json:"name"`
				} `json:"damage_type"`
			} `json:"damage,omitempty"`
		} `json:"actions"`
		ChallengeRating float64 `json:"challenge_rating"`
		XP              int     `json:"xp"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("error decoding monster data: %w", err)
	}

	// Convert to our internal Monster struct
	monster := &Monster{
		Index:        apiResponse.Index,
		Name:         apiResponse.Name,
		Size:         apiResponse.Size,
		Type:         apiResponse.Type,
		Alignment:    apiResponse.Alignment,
		ArmorClass:   apiResponse.ArmorClass,
		HitDice:      apiResponse.HitDice,
		Speed: MonsterSpeed{
			Walk:   parseSpeed(apiResponse.Speed.Walk),
			Swim:   parseSpeed(apiResponse.Speed.Swim),
			Fly:    parseSpeed(apiResponse.Speed.Fly),
			Climb:  parseSpeed(apiResponse.Speed.Climb),
			Burrow: parseSpeed(apiResponse.Speed.Burrow),
		},
		Strength:     apiResponse.Strength,
		Dexterity:    apiResponse.Dexterity,
		Constitution: apiResponse.Constitution,
		Intelligence: apiResponse.Intelligence,
		Wisdom:       apiResponse.Wisdom,
		Charisma:     apiResponse.Charisma,
		StrengthMod:  (apiResponse.Strength - 10) / 2,
		DexterityMod: (apiResponse.Dexterity - 10) / 2,
		ConMod:       (apiResponse.Constitution - 10) / 2,
		IntMod:       (apiResponse.Intelligence - 10) / 2,
		WisdomMod:    (apiResponse.Wisdom - 10) / 2,
		CharismaMod:  (apiResponse.Charisma - 10) / 2,
		Actions:      make([]MonsterAction, 0, len(apiResponse.Actions)),
		ChallengeRating: apiResponse.ChallengeRating,
		XP:           apiResponse.XP,
	}

	// Process actions
	for _, action := range apiResponse.Actions {
		monsterAction := MonsterAction{
			Name:        action.Name,
			Description: action.Description,
			AttackBonus: action.AttackBonus,
		}

		// Process damage
		if len(action.Damage) > 0 {
			diceCount, diceValue, bonus := parseDamageDice(action.Damage[0].DamageDice)
			monsterAction.Damage = DamageInfo{
				DiceCount: diceCount,
				DiceValue: diceValue,
				Bonus:     bonus,
				Type:      action.Damage[0].DamageType.Name,
			}
		}

		monster.Actions = append(monster.Actions, monsterAction)
	}

	// Store in cache
	c.cache.Set(cacheKey, monster)
	return monster, nil
}

// GetSpell fetches a spell from the SRD API
func (c *SRDClient) GetSpell(index string) (*Spell, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("spell:%s", index)
	if data, found := c.cache.Get(cacheKey); found {
		return data.(*Spell), nil
	}

	// Fetch from API
	url := fmt.Sprintf("%s/spells/%s", c.baseURL, index)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching spell data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-OK status: %d", resp.StatusCode)
	}

	var apiResponse struct {
		Index       string   `json:"index"`
		Name        string   `json:"name"`
		Level       int      `json:"level"`
		School      struct {
			Name string `json:"name"`
		} `json:"school"`
		CastingTime string   `json:"casting_time"`
		Range       string   `json:"range"`
		Components  []string `json:"components"`
		Duration    string   `json:"duration"`
		Description []string `json:"desc"`
		HigherLevel []string `json:"higher_level,omitempty"`
		Classes     []struct {
			Name string `json:"name"`
		} `json:"classes"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("error decoding spell data: %w", err)
	}

	// Convert to our internal Spell struct
	spell := &Spell{
		Index:       apiResponse.Index,
		Name:        apiResponse.Name,
		Level:       apiResponse.Level,
		School:      apiResponse.School.Name,
		CastingTime: apiResponse.CastingTime,
		Range:       apiResponse.Range,
		Components:  apiResponse.Components,
		Duration:    apiResponse.Duration,
		Description: combineStringSlice(apiResponse.Description),
		Classes:     make([]string, 0, len(apiResponse.Classes)),
	}

	if len(apiResponse.HigherLevel) > 0 {
		spell.HigherLevel = combineStringSlice(apiResponse.HigherLevel)
	}

	for _, class := range apiResponse.Classes {
		spell.Classes = append(spell.Classes, class.Name)
	}

	// Store in cache
	c.cache.Set(cacheKey, spell)
	return spell, nil
}

// Helper functions

// parseSpeed converts a speed string like "30 ft." to an integer
func parseSpeed(speed string) int {
	var value int
	fmt.Sscanf(speed, "%d", &value)
	return value
}

// parseDamageDice parses a damage dice string like "2d6+2"
func parseDamageDice(dice string) (count, value, bonus int) {
	fmt.Sscanf(dice, "%dd%d+%d", &count, &value, &bonus)
	return
}

// combineStringSlice joins a string slice into a single string
func combineStringSlice(slice []string) string {
	result := ""
	for _, str := range slice {
		result += str + " "
	}
	return result
}

// Monster is defined in the models package, but we need it here for the client
// This is imported from internal/models/combat.go
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
