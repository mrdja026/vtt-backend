# D&D Combat API Schema

This document provides a comprehensive schema of all API endpoints for the D&D 5e Combat API. It is intended for UI developers to understand the available endpoints, their request and response formats, and authentication requirements.

## Authentication

All authenticated endpoints require a JWT token in the Authorization header.

Format: `Authorization: Bearer <token>`

## Base URL

The base URL for all API endpoints is:

```
https://{server-url}/api/v1
```

## Endpoints

### Debug

#### Ping

Checks if the API is running.

- URL: `/debug/ping`
- Method: `GET`
- Auth required: No

**Response**

```json
{
  "message": "pong"
}
```

### Authentication

#### Register

Registers a new user account.

- URL: `/auth/register`
- Method: `POST`
- Auth required: No

**Request**

```json
{
  "username": "string",
  "email": "string",
  "password": "string"
}
```

**Response**

```json
{
  "id": "string",
  "username": "string",
  "email": "string",
  "created_at": "string"
}
```

**Error Responses**

| Status | Description |
|--------|-------------|
| 400 | Invalid request format |
| 409 | Username or email already exists |

#### Login

Authenticates a user and returns a JWT token.

- URL: `/auth/login`
- Method: `POST`
- Auth required: No

**Request**

```json
{
  "username": "string",
  "password": "string"
}
```

**Response**

```json
{
  "token": "string",
  "user": {
    "id": "string",
    "username": "string",
    "email": "string"
  }
}
```

**Error Responses**

| Status | Description |
|--------|-------------|
| 400 | Invalid request format |
| 401 | Invalid credentials |

### Characters

#### Create Character

Creates a new D&D character.

- URL: `/characters`
- Method: `POST`
- Auth required: Yes

**Request**

```json
{
  "name": "string",
  "race": "string",
  "class": "string",
  "level": "integer",
  "strength": "integer",
  "dexterity": "integer",
  "constitution": "integer",
  "intelligence": "integer",
  "wisdom": "integer",
  "charisma": "integer",
  "hit_points": "integer",
  "max_hit_points": "integer",
  "armor_class": "integer",
  "equipment": [
    {
      "name": "string",
      "type": "string",
      "damage": "string",
      "properties": ["string"]
    }
  ],
  "spells": [
    {
      "id": "string",
      "name": "string",
      "level": "integer"
    }
  ]
}
```

**Response**

```json
{
  "id": "string",
  "user_id": "string",
  "name": "string",
  "race": "string",
  "class": "string",
  "level": "integer",
  "strength": "integer",
  "dexterity": "integer",
  "constitution": "integer",
  "intelligence": "integer",
  "wisdom": "integer",
  "charisma": "integer",
  "hit_points": "integer",
  "max_hit_points": "integer",
  "armor_class": "integer",
  "equipment": [
    {
      "name": "string",
      "type": "string",
      "damage": "string",
      "properties": ["string"]
    }
  ],
  "spells": [
    {
      "id": "string",
      "name": "string",
      "level": "integer"
    }
  ],
  "created_at": "string",
  "updated_at": "string"
}
```

**Error Responses**

| Status | Description |
|--------|-------------|
| 400 | Invalid request format |
| 401 | Unauthorized |

#### Get Character

Retrieves a character by ID.

- URL: `/characters/{id}`
- Method: `GET`
- Auth required: Yes

**URL Parameters**

| Parameter | Description |
|-----------|-------------|
| id | Character ID |

**Response**

```json
{
  "id": "string",
  "user_id": "string",
  "name": "string",
  "race": "string",
  "class": "string",
  "level": "integer",
  "strength": "integer",
  "dexterity": "integer",
  "constitution": "integer",
  "intelligence": "integer",
  "wisdom": "integer",
  "charisma": "integer",
  "hit_points": "integer",
  "max_hit_points": "integer",
  "armor_class": "integer",
  "equipment": [
    {
      "name": "string",
      "type": "string",
      "damage": "string",
      "properties": ["string"]
    }
  ],
  "spells": [
    {
      "id": "string",
      "name": "string",
      "level": "integer"
    }
  ],
  "created_at": "string",
  "updated_at": "string"
}
```

**Error Responses**

| Status | Description |
|--------|-------------|
| 401 | Unauthorized |
| 404 | Character not found |

#### List Characters

Retrieves all characters owned by the authenticated user.

- URL: `/characters`
- Method: `GET`
- Auth required: Yes

**Query Parameters**

| Parameter | Description |
|-----------|-------------|
| limit | Maximum number of characters to return (default: 20) |
| offset | Offset for pagination (default: 0) |

**Response**

```json
{
  "total": "integer",
  "limit": "integer",
  "offset": "integer",
  "characters": [
    {
      "id": "string",
      "user_id": "string",
      "name": "string",
      "race": "string",
      "class": "string",
      "level": "integer",
      "hit_points": "integer",
      "max_hit_points": "integer",
      "armor_class": "integer",
      "created_at": "string",
      "updated_at": "string"
    }
  ]
}
```

**Error Responses**

| Status | Description |
|--------|-------------|
| 401 | Unauthorized |

### Games

#### Create Game

Creates a new game session.

- URL: `/games`
- Method: `POST`
- Auth required: Yes

**Request**

```json
{
  "name": "string",
  "description": "string",
  "player_ids": ["string"],
  "settings": {
    "use_grid": "boolean",
    "fog_of_war": "boolean",
    "difficulty": "string"
  }
}
```

**Response**

```json
{
  "id": "string",
  "name": "string",
  "description": "string",
  "dm_user_id": "string",
  "player_ids": ["string"],
  "settings": {
    "use_grid": "boolean",
    "fog_of_war": "boolean",
    "difficulty": "string"
  },
  "created_at": "string",
  "updated_at": "string"
}
```

**Error Responses**

| Status | Description |
|--------|-------------|
| 400 | Invalid request format |
| 401 | Unauthorized |

#### Get Game

Retrieves a game by ID.

- URL: `/games/{id}`
- Method: `GET`
- Auth required: Yes

**URL Parameters**

| Parameter | Description |
|-----------|-------------|
| id | Game ID |

**Response**

```json
{
  "id": "string",
  "name": "string",
  "description": "string",
  "dm_user_id": "string",
  "player_ids": ["string"],
  "players": [
    {
      "id": "string",
      "username": "string",
      "character_id": "string",
      "character_name": "string"
    }
  ],
  "settings": {
    "use_grid": "boolean",
    "fog_of_war": "boolean",
    "difficulty": "string"
  },
  "active_combat_id": "string",
  "created_at": "string",
  "updated_at": "string"
}
```

**Error Responses**

| Status | Description |
|--------|-------------|
| 401 | Unauthorized |
| 404 | Game not found |

#### List Games

Retrieves all games where the authenticated user is a DM or player.

- URL: `/games`
- Method: `GET`
- Auth required: Yes

**Query Parameters**

| Parameter | Description |
|-----------|-------------|
| limit | Maximum number of games to return (default: 20) |
| offset | Offset for pagination (default: 0) |
| dm_only | If true, only return games where user is DM (default: false) |

**Response**

```json
{
  "total": "integer",
  "limit": "integer",
  "offset": "integer",
  "games": [
    {
      "id": "string",
      "name": "string",
      "description": "string",
      "dm_user_id": "string",
      "player_count": "integer",
      "created_at": "string",
      "updated_at": "string"
    }
  ]
}
```

**Error Responses**

| Status | Description |
|--------|-------------|
| 401 | Unauthorized |

#### Update Game

Updates an existing game.

- URL: `/games/{id}`
- Method: `PUT`
- Auth required: Yes

**URL Parameters**

| Parameter | Description |
|-----------|-------------|
| id | Game ID |

**Request**

```json
{
  "name": "string",
  "description": "string",
  "player_ids": ["string"],
  "settings": {
    "use_grid": "boolean",
    "fog_of_war": "boolean",
    "difficulty": "string"
  }
}
```

**Response**

```json
{
  "id": "string",
  "name": "string",
  "description": "string",
  "dm_user_id": "string",
  "player_ids": ["string"],
  "settings": {
    "use_grid": "boolean",
    "fog_of_war": "boolean",
    "difficulty": "string"
  },
  "created_at": "string",
  "updated_at": "string"
}
```

**Error Responses**

| Status | Description |
|--------|-------------|
| 400 | Invalid request format |
| 401 | Unauthorized |
| 403 | User is not DM of this game |
| 404 | Game not found |

### Combat

#### Initiate Combat

Starts a new combat encounter.

- URL: `/combat`
- Method: `POST`
- Auth required: Yes

**Request**

```json
{
  "game_id": "string",
  "participants": ["string"],
  "monster_ids": ["string"],
  "environment": "string"
}
```

**Response**

```json
{
  "id": "string",
  "dm_user_id": "string",
  "current_turn_index": "integer",
  "round_number": "integer",
  "status": "string",
  "initiative": [
    {
      "id": "string",
      "name": "string",
      "initiative": "integer",
      "dexterity": "integer",
      "is_character": "boolean"
    }
  ],
  "participants": [
    {
      "id": "string",
      "user_id": "string",
      "character_id": "string",
      "monster_id": "string",
      "name": "string",
      "type": "string",
      "hp": "integer",
      "max_hp": "integer",
      "ac": "integer",
      "initiative": "integer",
      "position": [0, 0],
      "conditions": ["string"]
    }
  ],
  "battlefield": {
    "width": "integer",
    "height": "integer",
    "grid": {
      "x,y": "string"
    },
    "terrain": {
      "x,y": "string"
    },
    "obstacles": {
      "x,y": "boolean"
    }
  },
  "environment": "string",
  "created_at": "string",
  "updated_at": "string"
}
```

**Error Responses**

| Status | Description |
|--------|-------------|
| 400 | Invalid request format |
| 401 | Unauthorized |
| 403 | User is not DM of this game |

#### Get Combat

Retrieves the current state of a combat encounter.

- URL: `/combat/{id}`
- Method: `GET`
- Auth required: Yes

**URL Parameters**

| Parameter | Description |
|-----------|-------------|
| id | Combat ID |

**Response**

```json
{
  "id": "string",
  "dm_user_id": "string",
  "current_turn_index": "integer",
  "round_number": "integer",
  "status": "string",
  "initiative": [
    {
      "id": "string",
      "name": "string",
      "initiative": "integer",
      "dexterity": "integer",
      "is_character": "boolean"
    }
  ],
  "participants": [
    {
      "id": "string",
      "user_id": "string",
      "character_id": "string",
      "monster_id": "string",
      "name": "string",
      "type": "string",
      "hp": "integer",
      "max_hp": "integer",
      "ac": "integer",
      "initiative": "integer",
      "position": [0, 0],
      "conditions": ["string"]
    }
  ],
  "battlefield": {
    "width": "integer",
    "height": "integer",
    "grid": {
      "x,y": "string"
    },
    "terrain": {
      "x,y": "string"
    },
    "obstacles": {
      "x,y": "boolean"
    }
  },
  "environment": "string",
  "current_actor": {
    "id": "string",
    "name": "string",
    "type": "string"
  },
  "created_at": "string",
  "updated_at": "string"
}
```

**Error Responses**

| Status | Description |
|--------|-------------|
| 401 | Unauthorized |
| 404 | Combat not found |

#### Perform Action

Performs a combat action.

- URL: `/combat/{id}/action`
- Method: `POST`
- Auth required: Yes

**URL Parameters**

| Parameter | Description |
|-----------|-------------|
| id | Combat ID |

**Request**

```json
{
  "actor_id": "string",
  "type": "string",
  "target_ids": ["string"],
  "spell_id": "string",
  "weapon_name": "string",
  "movement_path": [[0, 0], [1, 0], [1, 1]],
  "extra_data": {
    "key": "value"
  }
}
```

The `type` field can be one of: `attack`, `cast_spell`, `move`, `dodge`, `help`, `hide`, `disengage`, `dash`, `use_item`

**Response**

```json
{
  "success": "boolean",
  "description": "string",
  "damage": "integer",
  "damage_type": "string",
  "healing": "integer",
  "target_effect": "string",
  "combat": {
    "id": "string",
    "current_turn_index": "integer",
    "round_number": "integer",
    "status": "string",
    "participants": [
      {
        "id": "string",
        "name": "string",
        "hp": "integer",
        "max_hp": "integer",
        "position": [0, 0],
        "conditions": ["string"]
      }
    ]
  }
}
```

**Error Responses**

| Status | Description |
|--------|-------------|
| 400 | Invalid request format |
| 401 | Unauthorized |
| 403 | Not actor's turn or user does not control actor |
| 404 | Combat not found |
| 409 | Invalid action |

#### End Turn

Ends the current actor's turn.

- URL: `/combat/{id}/end-turn`
- Method: `POST`
- Auth required: Yes

**URL Parameters**

| Parameter | Description |
|-----------|-------------|
| id | Combat ID |

**Request**

```json
{
  "actor_id": "string"
}
```

**Response**

```json
{
  "success": "boolean",
  "description": "string",
  "next_actor": {
    "id": "string",
    "name": "string",
    "type": "string"
  },
  "round_number": "integer",
  "is_new_round": "boolean"
}
```

**Error Responses**

| Status | Description |
|--------|-------------|
| 400 | Invalid request format |
| 401 | Unauthorized |
| 403 | Not actor's turn or user does not control actor |
| 404 | Combat not found |



**URL Parameters**

| Parameter | Description |
|-----------|-------------|
| id | Combat ID |

**Query Parameters**

| Parameter | Description |
|-----------|-------------|
| token | JWT token |



The server will emit the following events:

| Event | Description | Data |
|-------|-------------|------|
| `combat_started` | Combat has begun | Combat object |
| `combat_updated` | Combat state has changed | Combat object |
| `turn_changed` | Turn has changed to a new actor | `{actor_id, actor_name, round_number}` |
| `action_performed` | An action was performed | Action result object |
| `combatant_updated` | A combatant's state changed | Combatant object |
| `combat_ended` | Combat has ended | `{id, winner_type}` |

**Client Messages**

Clients can send these messages to the server:

| Message | Description | Data |
|---------|-------------|------|
| `ready` | Indicates client is ready to receive updates | `{client_id}` |
| `ping` | Ping to keep connection alive | `{}` |

## Data Models

### Character

```json
{
  "id": "string",
  "user_id": "string",
  "name": "string",
  "race": "string",
  "class": "string",
  "level": "integer",
  "strength": "integer",
  "dexterity": "integer",
  "constitution": "integer",
  "intelligence": "integer",
  "wisdom": "integer",
  "charisma": "integer",
  "hit_points": "integer",
  "max_hit_points": "integer",
  "armor_class": "integer",
  "equipment": [
    {
      "name": "string",
      "type": "string",
      "damage": "string",
      "properties": ["string"]
    }
  ],
  "spells": [
    {
      "id": "string",
      "name": "string",
      "level": "integer"
    }
  ],
  "created_at": "string",
  "updated_at": "string"
}
```

### Monster

```json
{
  "index": "string",
  "name": "string",
  "size": "string",
  "type": "string",
  "alignment": "string",
  "armor_class": "integer",
  "hit_dice": "string",
  "speed": {
    "walk": "integer",
    "swim": "integer",
    "fly": "integer",
    "climb": "integer",
    "burrow": "integer"
  },
  "strength": "integer",
  "dexterity": "integer",
  "constitution": "integer",
  "intelligence": "integer",
  "wisdom": "integer",
  "charisma": "integer",
  "strength_mod": "integer",
  "dexterity_mod": "integer",
  "constitution_mod": "integer",
  "intelligence_mod": "integer",
  "wisdom_mod": "integer",
  "charisma_mod": "integer",
  "actions": [
    {
      "name": "string",
      "description": "string",
      "attack_bonus": "integer",
      "range": "integer",
      "damage": {
        "dice_count": "integer",
        "dice_value": "integer",
        "bonus": "integer",
        "type": "string"
      }
    }
  ],
  "challenge_rating": "number",
  "xp": "integer"
}
```

### Combat

```json
{
  "id": "string",
  "dm_user_id": "string",
  "current_turn_index": "integer",
  "round_number": "integer",
  "status": "string",
  "initiative": [
    {
      "id": "string",
      "name": "string",
      "initiative": "integer",
      "dexterity": "integer",
      "is_character": "boolean"
    }
  ],
  "participants": [
    {
      "id": "string",
      "user_id": "string",
      "character_id": "string",
      "monster_id": "string",
      "name": "string",
      "type": "string",
      "hp": "integer",
      "max_hp": "integer",
      "ac": "integer",
      "initiative": "integer",
      "position": [0, 0],
      "conditions": ["string"]
    }
  ],
  "battlefield": {
    "width": "integer",
    "height": "integer",
    "grid": {
      "x,y": "string"
    },
    "terrain": {
      "x,y": "string"
    },
    "obstacles": {
      "x,y": "boolean"
    }
  },
  "environment": "string",
  "created_at": "string",
  "updated_at": "string"
}
```

### Spell

```json
{
  "index": "string",
  "name": "string",
  "level": "integer",
  "school": "string",
  "casting_time": "string",
  "range": "string",
  "components": ["string"],
  "duration": "string",
  "description": "string",
  "higher_level": "string",
  "classes": ["string"]
}
```

## Status Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request - Invalid input |
| 401 | Unauthorized - Authentication failed |
| 403 | Forbidden - User doesn't have permission |
| 404 | Not Found - Resource not found |
| 409 | Conflict - Resource already exists or state conflict |
| 422 | Unprocessable Entity - Validation error |
| 500 | Server Error - Something went wrong on the server |

## Rate Limits

- 100 requests per minute per IP address
- 1000 requests per hour per user

Limits are specified in response headers:
- `X-RateLimit-Limit`: Total requests allowed in time window
- `X-RateLimit-Remaining`: Remaining requests in current window
- `X-RateLimit-Reset`: Time when the limit resets (Unix timestamp)

## API Versioning

The API is versioned in the URL path (`/api/v1`). When breaking changes are introduced, a new version will be released (e.g., `/api/v2`).

## Error Handling

All errors follow a consistent format:

```json
{
  "error": {
    "code": "string",
    "message": "string",
    "details": "string"
  }
}
```

## Pagination

List endpoints support pagination with the following query parameters:
- `limit`: Number of items per page (default: 20, max: 100)
- `offset`: Number of items to skip (default: 0)

Pagination metadata is included in the response:

```json
{
  "total": "integer",
  "limit": "integer",
  "offset": "integer",
  "items": []
}
```