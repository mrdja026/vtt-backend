<img src="https://r2cdn.perplexity.ai/pplx-full-logo-primary-dark%402x.png" class="logo" width="120"/>

# Comprehensive D\&D Turn-Based Combat Backend Requirements Document

This document provides a detailed framework for developing a Dungeons \& Dragons turn-based combat system backend that integrates with the D\&D 5e SRD API while following Go Blueprint best practices. The requirements outlined below will guide Replit in implementing a robust, efficient, and maintainable backend service.

## Current Requirements Analysis

The initial draft outlines a solid foundation for a D\&D combat API with core endpoints for user management, character creation, game sessions, and combat actions[^1]. The proposed structure follows modern API design practices with versioned endpoints and clear authentication requirements. The Go Blueprint-inspired structure provides a good starting point, with separation of concerns through the cmd/internal/pkg organization[^1].

However, several enhancements can be made to improve integration with external services, optimize the code structure, and follow the latest best practices for Go API development.

### Core Functionality Matrix (Enhanced)

| Feature | HTTP Method | Endpoint | Auth Required | Description |
| :-- | :-- | :-- | :-- | :-- |
| User Registration | POST | /api/v1/auth/register | No | Create new user account |
| User Login | POST | /api/v1/auth/login | No | Authenticate user and return JWT |
| Character Create | POST | /api/v1/characters | Yes | Create a new D\&D character |
| Character Retrieve | GET | /api/v1/characters/{id} | Yes | Get character details |
| Character List | GET | /api/v1/characters | Yes | List user's characters |
| Game Session | POST | /api/v1/games | Yes | Create a new game session |
| Combat Initiate | POST | /api/v1/combat | Yes | Start combat encounter |
| Combat Action | POST | /api/v1/combat/{id}/move | Yes | Execute combat action |
| Combat Status | GET | /api/v1/combat/{id} | Yes | Get current combat state |

## D\&D 5e SRD API Integration

The D\&D 5e SRD API provides an extensive collection of official 5th Edition SRD content that can be leveraged to enhance our combat system[^9][^13]. Rather than implementing D\&D rules from scratch, we can use this API to access canonical data and rules.

### API Integration Layer

```go
// pkg/dnd5e/srd_client.go
package dnd5e

import (
    "encoding/json"
    "fmt"
    "net/http"
)

const BaseURL = "https://www.dnd5eapi.co/api/2014"

type SRDClient struct {
    httpClient *http.Client
    cache      Cache
}

func NewSRDClient(client *http.Client, cache Cache) *SRDClient {
    return &SRDClient{
        httpClient: client,
        cache:      cache,
    }
}

func (c *SRDClient) GetClass(index string) (*ClassData, error) {
    // Check cache first
    cacheKey := fmt.Sprintf("class:%s", index)
    if data, found := c.cache.Get(cacheKey); found {
        return data.(*ClassData), nil
    }
    
    // Fetch from API
    resp, err := c.httpClient.Get(fmt.Sprintf("%s/classes/%s", BaseURL, index))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var classData ClassData
    if err := json.NewDecoder(resp.Body).Decode(&classData); err != nil {
        return nil, err
    }
    
    // Store in cache
    c.cache.Set(cacheKey, &classData)
    return &classData, nil
}

// Additional methods for monsters, spells, equipment, etc.
```


### Resource Types to Integrate

Based on the D\&D 5e SRD API documentation[^3][^6][^9][^13], we should integrate the following resources:

1. Characters: Classes, subclasses, races
2. Combat: Abilities, skills, spell effects
3. Equipment: Weapons, armor, magic items
4. Monsters: Stats, abilities, challenge ratings

## Enhanced Go Blueprint Structure

The Go Blueprint CLI tool provides an excellent foundation for Go project structure[^7][^8][^15]. We'll enhance our implementation using more advanced features provided by the tool.

### Project Generation Command

```bash
go-blueprint create \
  --name dnd-combat \
  --framework gin \
  --driver sqlite \
  --advanced \
  --feature docker \
  --feature githubaction \
  --git commit
```

The addition of the websocket feature will enable real-time updates for multiple players in a combat session[^12][^15].

### Expanded Project Architecture

```
/.
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── auth/                    # Authentication logic
│   │   ├── handler.go
│   │   ├── middleware.go
│   │   ├── repository.go
│   │   └── service.go
│   ├── character/               # Character management
│   │   ├── handler.go
│   │   ├── repository.go
│   │   └── service.go
│   ├── combat/                  # Combat mechanics
│   │   ├── handler.go
│   │   ├── repository.go
│   │   └── service.go
│   ├── game/                    # Game session management
│   │   ├── handler.go
│   │   ├── repository.go
│   │   └── service.go
│   └── models/                  # Domain models
│       ├── character.go
│       ├── combat.go
│       └── user.go
├── pkg/
│   ├── dnd5e/                   # D&D rules implementation
│   │   ├── srd_client.go        # API client for D&D 5e SRD
│   │   ├── dice.go              # Dice rolling utilities
│   │   └── combat.go            # Combat utilities
│   ├── middleware/              # Shared middleware
│   │   ├── logging.go
│   │   └── auth.go
│   ├── database/                # Database connection
│   │   └── sqlite.go
│   └── websocket/               # Websocket implementation
│       └── hub.go
├── api/
│   └── v1/
│       └── routes.go            # API route definitions
├── config/
│   └── config.yaml              # Application configuration
├── tests/
│   ├── unit/                    # Unit tests
│   └── integration/             # Integration tests
└── docs/
    └── swagger/                 # API documentation
```


## API Documentation Best Practices

Following Microsoft's API design best practices[^10] and Stoplight's API documentation guide[^5], we should implement:

1. **Clear Pagination**: Implement offset/limit pagination with predictable defaults:

```
GET /api/v1/characters?limit=25&offset=0
```

2. **Field Selection**: Allow clients to specify only needed fields:

```
GET /api/v1/characters?fields=id,name,class,level
```

3. **Documented Error Responses**: Standardize error formats and document status codes[^5][^17]

```json
{
  "error": {
    "code": "COMBAT_NOT_FOUND",
    "message": "Combat session not found",
    "details": "Combat session with ID 123 does not exist or has ended"
  }
}
```

4. **OpenAPI Documentation**: Generate comprehensive API documentation using OpenAPI/Swagger[^10][^17]

## D\&D Combat Implementation Details

### Combat Initialization with SRD Integration

```go
// internal/combat/handler.go
func (h *Handler) InitiateCombat(c *gin.Context) {
    var req struct {
        ParticipantIDs []string `json:"participants"`
        Environment    string   `json:"environment"`
        MonsterIDs     []string `json:"monster_ids"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "Invalid request"})
        return
    }
    
    // Fetch character data from database
    characters, err := h.characterRepo.GetMultiple(req.ParticipantIDs)
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to retrieve characters"})
        return
    }
    
    // Fetch monster data from SRD API
    monsters := make([]Monster, 0, len(req.MonsterIDs))
    for _, monsterID := range req.MonsterIDs {
        monster, err := h.srdClient.GetMonster(monsterID)
        if err != nil {
            h.logger.Error("failed to fetch monster", "id", monsterID, "error", err)
            continue
        }
        monsters = append(monsters, *monster)
    }
    
    // Calculate initiative for all participants
    initiative := h.combatService.CalculateInitiative(characters, monsters)
    
    // Create and store combat session
    combat, err := h.combatService.CreateCombat(characters, monsters, initiative, req.Environment)
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to create combat"})
        return
    }
    
    // Notify clients via websocket
    h.websocketHub.BroadcastToSession(combat.ID, "combat_started", combat)
    
    c.JSON(200, combat)
}
```


## Database Schema Implementation

```go
// internal/models/character.go
type Character struct {
    ID              string    `json:"id" db:"id"`
    UserID          string    `json:"user_id" db:"user_id"`
    Name            string    `json:"name" db:"name"`
    Race            string    `json:"race" db:"race"`
    Class           string    `json:"class" db:"class"`
    Level           int       `json:"level" db:"level"`
    Strength        int       `json:"strength" db:"strength"`
    Dexterity       int       `json:"dexterity" db:"dexterity"`
    Constitution    int       `json:"constitution" db:"constitution"`
    Intelligence    int       `json:"intelligence" db:"intelligence"`
    Wisdom          int       `json:"wisdom" db:"wisdom"`
    Charisma        int       `json:"charisma" db:"charisma"`
    HitPoints       int       `json:"hit_points" db:"hit_points"`
    MaxHitPoints    int       `json:"max_hit_points" db:"max_hit_points"`
    ArmorClass      int       `json:"armor_class" db:"armor_class"`
    EquipmentJSON   string    `json:"-" db:"equipment_json"`
    SpellsJSON      string    `json:"-" db:"spells_json"`
    CreatedAt       time.Time `json:"created_at" db:"created_at"`
    UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}
```


## Replit Implementation Recommendations

To maximize the effectiveness of this document for Replit development, the following additional sections should be included:

### 1. Local Development Setup

```bash
# Clone the repository
git clone https://github.com/your-org/dnd-combat.git

# Install Go Blueprint
go install github.com/melkeydev/go-blueprint@latest

# Generate project structure
go-blueprint create --name dnd-combat --framework gin --driver sqlite --advanced --feature docker --feature githubaction --feature websocket --git commit

# Install dependencies
cd dnd-combat
go mod tidy

# Run the application
make run
```


### 2. Environment Variables

```
# .env
DB_PATH=./data/dnd_combat.db
JWT_SECRET=your-secret-key-here
SRD_API_BASE_URL=https://www.dnd5eapi.co/api/2014
PORT=8080
ENV=development
```


### 3. Testing Strategy

```go
// tests/integration/combat_test.go
func TestCombatInitialization(t *testing.T) {
    // Set up test environment
    // ...
    
    // Create test characters
    chars := createTestCharacters(t, api)
    
    // Initialize combat
    resp, err := http.Post(
        fmt.Sprintf("%s/api/v1/combat", testServer.URL),
        "application/json",
        strings.NewReader(`{"participants": ["char1", "char2"], "monster_ids": ["goblin", "orc"]}`),
    )
    
    // Assert response
    assert.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
    
    // Parse response
    var combat Combat
    assert.NoError(t, json.NewDecoder(resp.Body).Decode(&combat))
    
    // Verify combat state
    assert.NotEmpty(t, combat.ID)
    assert.Equal(t, 4, len(combat.TurnOrder))
    assert.Equal(t, 1, combat.RoundNumber)
}
```


## Implementation Checklist

1. Scaffold project using Go Blueprint CLI with specified features
2. Set up database schema and migrations
3. Implement D\&D 5e SRD API client with caching mechanisms
4. Create authentication flow with JWT
5. Implement character creation and management endpoints
6. Build combat state machine with initiative tracking
7. Integrate websocket for real-time combat updates
8. Implement combat actions (attack, spell casting, movement)
9. Add environment effects and combat modifiers
10. Create comprehensive test suite for all components
11. Generate OpenAPI/Swagger documentation
12. Set up Docker deployment for Replit

## Conclusion

This enhanced requirements document provides a comprehensive roadmap for implementing a D\&D turn-based combat backend system using Go Blueprint best practices and integrating with the D\&D 5e SRD API. By following this structure and implementing the suggested components, Replit will be able to create a robust, scalable, and maintainable backend service that adheres to modern development practices and provides a rich feature set for D\&D combat mechanics.

The integration with the D\&D 5e SRD API ensures accuracy and completeness of game rules, while the Go Blueprint structure provides a solid foundation for code organization and maintainability. The addition of websocket support enables real-time multiplayer experiences, making the combat system more engaging and interactive.

<div style="text-align: center">⁂</div>

[^1]: format-the-requirements-I-as-a-developer-want-to.md

[^2]: https://5e-bits.github.io/docs/

[^3]: https://5e-bits.github.io/docs/api/get-list-of-all-available-resources-for-an-endpoint

[^4]: https://www.reddit.com/r/DnD/comments/4cpxpi/anyone_know_of_any_dd_related_datasets/

[^5]: https://stoplight.io/api-documentation-guide

[^6]: https://5e-bits.github.io/docs/api/get-a-class-by-index

[^7]: https://github.com/Melkeydev/go-blueprint

[^8]: https://go-blueprint.dev

[^9]: https://github.com/5e-bits/5e-srd-api

[^10]: https://learn.microsoft.com/en-us/azure/architecture/best-practices/api-design

[^11]: https://5e-bits.github.io/docs/api/get-a-subclass-by-index

[^12]: https://raw.githubusercontent.com/Melkeydev/go-blueprint/main/README.md

[^13]: https://5e-bits.github.io/docs/

[^14]: https://www.reddit.com/r/ProductManagement/comments/15llhnx/documenting_api_requirements/

[^15]: https://docs.go-blueprint.dev/installation/

[^16]: https://www.dnd5eapi.co

[^17]: https://konghq.com/blog/learning-center/guide-to-api-documentation

[^18]: https://www.freepublicapis.com/dungeons-and-dragons

[^19]: https://liblab.com/blog/api-documentation-best-practices

[^20]: https://www.reddit.com/r/DnD/comments/5jmv0k/dnd5eapi_a_rest_api_for_dd_5th_edition/

[^21]: https://media.wizards.com/2016/downloads/DND/SRD-OGL_V5.1.pdf

[^22]: https://github.com/5e-bits/awesome-5e-srd

[^23]: https://5e-bits.github.io/docs/api/get-an-ability-score-by-index

[^24]: https://library.wiremock.org/catalog/api/d/dnd5eapi.co/dnd5eapi-co/

[^25]: https://github.com/5e-bits/docs

[^26]: https://www.5esrd.com/gamemastering/5th-edition-options/

[^27]: https://www.reddit.com/r/DnD/comments/1292p7v/searchable_spreadsheet_of_srd_51_creatures/

[^28]: https://rpg.stackexchange.com/questions/114955/difference-between-the-srd-and-the-basic-rules

[^29]: https://swagger.io/resources/articles/best-practices-in-api-design/

[^30]: https://www.apimatic.io/blog/2022/11/14-best-practices-to-write-openapi-for-better-api-consumption

[^31]: https://www.akana.com/blog/api-requirements-what-consider

[^32]: https://swagger.io/resources/articles/difference-between-api-documentation-specification/

[^33]: https://www.pandium.com/blogs/mastering-api-documentation-3-proven-best-practices-for-success

[^34]: https://stoplight.io/api-style-guides-guidelines-and-best-practices

[^35]: https://swagger.io/specification/

[^36]: https://cloud.google.com/apis/design

[^37]: https://www.altexsoft.com/blog/api-documentation/

[^38]: https://www.reddit.com/r/Eberron/comments/16lo5qn/what_is_a_schema/

[^39]: https://mojobob.com/roleplay/dnd5e/3e-5e_monster-conversion.html

[^40]: https://rpgbot.net/dnd5/characters/classes/artificer/

[^41]: https://rpgbot.net/dnd5/characters/classes/warlock/spells/

[^42]: https://rpgbot.net/dnd5/tools/monsterizer/

[^43]: https://rpgbot.net/dnd5/characters/classes/fighter/

[^44]: https://2minutetabletop.com/bits-and-bugs-1-dd-5e/

[^45]: https://5e-bits.github.io/docs/

[^46]: https://5e-bits.github.io/docs/api/get-a-spell-by-index

[^47]: https://www.aidedd.org/dnd-filters/monsters.php

[^48]: https://5e-bits.github.io/docs/tutorials/advanced/react-spell-cards

[^49]: https://www.reddit.com/r/golang/comments/1dezksp/golang_folder_structure_for_mid_size_project/

[^50]: https://victorpierre.dev/blog/five-go-interfaces-best-practices/

[^51]: https://docs.go-blueprint.dev

[^52]: https://martinheinz.dev/blog/5

[^53]: https://www.reddit.com/r/golang/comments/1dhvbkd/how_many_of_you_using_goblueprintdev/

[^54]: https://www.linkedin.com/posts/donald-lutz-5a9b0b2_github-melkeydevgo-blueprint-go-blueprint-activity-7282176084723347456-4jxr

[^55]: https://www.youtube.com/watch?v=1ZbQS6pOlSQ

[^56]: https://dev.epicgames.com/documentation/en-us/unreal-engine/blueprint-best-practices-in-unreal-engine

[^57]: https://www.linkedin.com/posts/roberthelmstetter_github-melkeydevgo-blueprint-go-blueprint-activity-7266823620704755714-L0l8

[^58]: https://www.youtube.com/watch?v=dxPakeBsgl4

[^59]: https://github.com/goauthentik/authentik/issues/8548

[^60]: https://google.github.io/styleguide/go/best-practices.html

[^61]: https://open5e.com/api-docs

[^62]: https://pypi.org/project/dnd5epy/

[^63]: https://swagger.io/blog/api-documentation/best-practices-in-api-documentation/

[^64]: https://docs.mulesoft.com/general/api-led-design

[^65]: https://learn.openapis.org/best-practices.html

[^66]: https://www.reddit.com/r/DnD/comments/g6885u/5e_class_and_subclass_compilation_chart/

[^67]: https://lesley.edu/article/empowering-students-the-5e-model-explained

[^68]: https://arcane.org/dd-5e-spell-component-database/

[^69]: https://www.reddit.com/r/UnearthedArcana/comments/8zvr6s/the_great_dd5e_monster_spreadsheet/

[^70]: https://www.dnd5eapi.co

[^71]: https://github.com/Melkeydev/go-blueprint/blob/main/main.go

[^72]: https://dev.to/hieunguyendev/initialize-a-go-project-with-go-blueprint-2bd1

[^73]: https://docs.cloudify.co/latest/bestpractices/agiledevelopmentbp/

[^74]: https://www.youtube.com/watch?v=-YeAt08picE

