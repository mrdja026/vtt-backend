# D&D 5e Turn-Based Combat API

A comprehensive Go-based Dungeons & Dragons 5th Edition turn-based combat backend with robust SRD API integration. This API provides a scalable and modular platform for tabletop role-playing game combat management.

## Features

- Character management with D&D 5e stats and abilities
- Combat system with initiative tracking and turn management
- Battlefield representation with grid, terrain, and obstacles
- Integration with D&D 5e SRD API for spells, monsters, and game rules
- WebSocket support for real-time combat updates
- Authentication and game session management

## Tech Stack

- Go programming language
- Gin web framework
- SQLite database (development) / PostgreSQL (production)
- Clean architecture (handlers, services, repositories)
- D&D 5e SRD API integration

## Prerequisites

- Go 1.18 or later
- Git
- SQLite (for development)
- PostgreSQL (for production)

## Local Setup

1. Clone the repository:

```bash
git clone https://github.com/yourusername/dnd-combat.git
cd dnd-combat
```

2. Set up environment variables:

```bash
# Create a .env file
touch .env

# Add the following variables to the .env file
DB_PATH=./data/dnd_combat.db
JWT_SECRET=your-secret-key-here
SRD_API_BASE_URL=https://www.dnd5eapi.co/api
PORT=8000
ENV=development
```

3. Initialize the database:

```bash
# Create the data directory if it doesn't exist
mkdir -p data
```

4. Install dependencies:

```bash
go mod tidy
```

5. Run the application:

```bash
go run cmd/api/main.go
```

The server should now be running at `http://localhost:8000`.

## Directory Structure

```
/.
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── auth/                    # Authentication logic
│   ├── character/               # Character management
│   ├── combat/                  # Combat mechanics
│   ├── game/                    # Game session management
│   └── models/                  # Domain models
├── pkg/
│   ├── dnd5e/                   # D&D rules implementation
│   ├── middleware/              # Shared middleware
│   ├── database/                # Database connection
│   └── websocket/               # Websocket implementation
├── api/
│   └── v1/                      # API route definitions
├── config/                      # Configuration
├── data/                        # SQLite database
```

## API Testing with cURL

Below are curl examples to test each endpoint of the API. For the authenticated endpoints, you'll need to replace `YOUR_JWT_TOKEN` with an actual token obtained from the login endpoint.

### Debug Endpoint

```bash
# Test API health
curl -X GET http://localhost:8000/api/v1/debug/ping
```

### Authentication Endpoints

```bash
# Register a new user
curl -X POST http://localhost:8000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "securepassword123"
  }'

# Login with user credentials
curl -X POST http://localhost:8000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "securepassword123"
  }'
```

### Character Management Endpoints

```bash
# Create a new character (requires auth token)
curl -X POST http://localhost:8000/api/v1/characters \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Aragorn",
    "race": "human",
    "class": "fighter",
    "level": 5,
    "strength": 16,
    "dexterity": 14,
    "constitution": 15,
    "intelligence": 12,
    "wisdom": 13,
    "charisma": 14,
    "hit_points": 45,
    "armor_class": 16
  }'

# Get a specific character by ID
curl -X GET http://localhost:8000/api/v1/characters/character_id_here \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# List all user's characters
curl -X GET http://localhost:8000/api/v1/characters \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Game Session Endpoints

```bash
# Create a new game session
curl -X POST http://localhost:8000/api/v1/games \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Epic Adventure",
    "description": "A journey through the forgotten realms",
    "player_ids": ["player1_id", "player2_id"]
  }'

# Get a specific game by ID
curl -X GET http://localhost:8000/api/v1/games/game_id_here \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# List all games
curl -X GET http://localhost:8000/api/v1/games \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Update a game
curl -X PUT http://localhost:8000/api/v1/games/game_id_here \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Updated Epic Adventure",
    "description": "A new journey through the forgotten realms",
    "player_ids": ["player1_id", "player2_id", "player3_id"]
  }'
```

### Combat Endpoints

```bash
# Initiate a combat encounter
curl -X POST http://localhost:8000/api/v1/combat \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "participants": ["character_id1", "character_id2"],
    "monster_ids": ["goblin", "orc"],
    "environment": "forest"
  }'

# Get combat state
curl -X GET http://localhost:8000/api/v1/combat/combat_id_here \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Perform an attack action
curl -X POST http://localhost:8000/api/v1/combat/combat_id_here/action \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "actor_id": "character_id1",
    "type": "attack",
    "target_ids": ["monster_id1"],
    "weapon_name": "longsword"
  }'

# Perform a spell casting action
curl -X POST http://localhost:8000/api/v1/combat/combat_id_here/action \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "actor_id": "character_id2",
    "type": "cast_spell",
    "target_ids": ["monster_id1", "monster_id2"],
    "spell_id": "fireball"
  }'

# Perform a movement action
curl -X POST http://localhost:8000/api/v1/combat/combat_id_here/action \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "actor_id": "character_id1",
    "type": "move",
    "movement_path": [[3,4], [4,4], [5,4]]
  }'

# End the current turn
curl -X POST http://localhost:8000/api/v1/combat/combat_id_here/end-turn \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "actor_id": "character_id1"
  }'
```

## Testing Flow

For a complete test flow, follow these steps:

1. Register a user
2. Login to get a JWT token
3. Create characters
4. Create a game session
5. Start a combat encounter with the created characters
6. Perform combat actions (attack, move, cast spells)
7. End turns and observe combat progression

## WebSocket Connection

For real-time combat updates, connect to the WebSocket endpoint:

```javascript
// JavaScript example
const socket = new WebSocket('ws://localhost:8000/api/v1/ws/combat/combat_id_here?token=YOUR_JWT_TOKEN');

socket.onopen = function(e) {
  console.log('Connection established');
};

socket.onmessage = function(event) {
  console.log('Data received:', event.data);
};
```

## Battlefield Implementation

The battlefield is represented as a 2D grid with:

- Coordinates for positioning characters and monsters
- Terrain types (normal, difficult, water)
- Obstacles (walls, trees, rocks)
- Movement validation based on D&D 5e rules

Movement is restricted by:
- Character/monster movement speed
- Obstacles and other combatants
- Valid adjacent squares (no diagonal jumping)

## D&D 5e Rules Implementation

This API implements core D&D 5e combat rules:

- Initiative determination using d20 + DEX modifier
- Attack rolls with advantage/disadvantage
- Spell casting with appropriate ranges and effects
- Saving throws against effects
- Combat actions (attack, dodge, help, hide, dash, disengage)

## Development

**Important Notice**: This is a proprietary codebase with all rights reserved. No contributions, modifications, forks, or derivative works are permitted without explicit written permission from the owner.

The code is provided for reference and educational purposes only. Any collaboration or development on this codebase requires formal authorization.

Please refer to the LICENSE file for complete terms and restrictions.

## License and Usage Restrictions

**PROPRIETARY AND CONFIDENTIAL**

© 2025 All Rights Reserved

This codebase is provided for reference and educational purposes only. All rights are reserved and no part of this codebase may be reproduced, distributed, or transmitted in any form or by any means, including photocopying, recording, or other electronic or mechanical methods, without the prior written permission of the owner.

**RESTRICTED ACCESS**
- No copying, modification, or distribution of this code is permitted
- No use in commercial or non-commercial applications without explicit written permission
- No derivative works may be created based on this codebase
- No license is granted to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of this software

Unauthorized use, reproduction, or distribution of this code or any portion of it may result in severe civil and criminal penalties.

## Acknowledgments

- D&D 5e SRD API: https://www.dnd5eapi.co
- Wizards of the Coast for Dungeons & Dragons