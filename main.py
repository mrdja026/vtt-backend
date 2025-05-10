from flask import Flask, jsonify, request
import os

app = Flask(__name__)

@app.route('/api/status', methods=['GET'])
def status():
    return jsonify({
        "status": "operational",
        "version": "1.0.0",
        "api": "D&D 5e Combat API",
        "endpoints": [
            {"path": "/api/status", "method": "GET", "description": "Get API status"},
            {"path": "/api/v1/auth/register", "method": "POST", "description": "Register a new user"},
            {"path": "/api/v1/auth/login", "method": "POST", "description": "Authenticate user"},
            {"path": "/api/v1/characters", "method": "GET", "description": "List characters"},
            {"path": "/api/v1/characters", "method": "POST", "description": "Create character"},
            {"path": "/api/v1/characters/{id}", "method": "GET", "description": "Get character details"},
            {"path": "/api/v1/combat", "method": "POST", "description": "Start combat encounter"},
            {"path": "/api/v1/combat/{id}", "method": "GET", "description": "Get combat state"},
            {"path": "/api/v1/combat/{id}/action", "method": "POST", "description": "Perform combat action"},
        ]
    })

# Auth endpoints
@app.route('/api/v1/auth/register', methods=['POST'])
def register():
    data = request.get_json()
    # In a real implementation, this would store the user in a database
    return jsonify({
        "message": "User registered successfully",
        "user": {
            "id": "user123",
            "username": data.get('username', ''),
            "email": data.get('email', '')
        }
    }), 201

@app.route('/api/v1/auth/login', methods=['POST'])
def login():
    data = request.get_json()
    # In a real implementation, this would validate credentials
    return jsonify({
        "token": "sample_jwt_token",
        "user": {
            "id": "user123",
            "username": data.get('username', ''),
            "email": "user@example.com"
        }
    })

# Character endpoints
@app.route('/api/v1/characters', methods=['GET'])
def list_characters():
    # Mock character data
    characters = [
        {
            "id": "char1",
            "name": "Aragorn",
            "race": "Human",
            "class": "Ranger",
            "level": 8,
            "hit_points": 75,
            "armor_class": 16
        },
        {
            "id": "char2",
            "name": "Gimli",
            "race": "Dwarf",
            "class": "Fighter",
            "level": 7,
            "hit_points": 85,
            "armor_class": 18
        }
    ]
    return jsonify({"characters": characters})

@app.route('/api/v1/characters', methods=['POST'])
def create_character():
    data = request.get_json()
    # In a real implementation, this would store the character in a database
    character = {
        "id": "char3",
        "name": data.get('name', ''),
        "race": data.get('race', ''),
        "class": data.get('class', ''),
        "level": data.get('level', 1),
        "hit_points": data.get('hit_points', 10),
        "armor_class": data.get('armor_class', 10)
    }
    return jsonify(character), 201

@app.route('/api/v1/characters/<character_id>', methods=['GET'])
def get_character(character_id):
    # Mock character data
    character = {
        "id": character_id,
        "name": "Aragorn",
        "race": "Human",
        "class": "Ranger",
        "level": 8,
        "strength": 16,
        "dexterity": 18,
        "constitution": 14,
        "intelligence": 12,
        "wisdom": 16,
        "charisma": 14,
        "hit_points": 75,
        "max_hit_points": 80,
        "armor_class": 16,
        "equipment": ["Longsword", "Bow", "Leather Armor"],
        "spells": []
    }
    return jsonify(character)

# Combat endpoints
@app.route('/api/v1/combat', methods=['POST'])
def start_combat():
    data = request.get_json()
    
    # Mock combat data
    combat = {
        "id": "combat123",
        "initiative": ["char1", "monster1", "char2", "monster2"],
        "participants": [
            {
                "id": "char1",
                "type": "character",
                "name": "Aragorn",
                "initiative": 18,
                "hp": 75,
                "max_hp": 80,
                "ac": 16,
                "position": [2, 3]
            },
            {
                "id": "char2",
                "type": "character",
                "name": "Gimli",
                "initiative": 12,
                "hp": 85,
                "max_hp": 85,
                "ac": 18,
                "position": [2, 4]
            },
            {
                "id": "monster1",
                "type": "monster",
                "name": "Goblin",
                "initiative": 15,
                "hp": 15,
                "max_hp": 15,
                "ac": 13,
                "position": [7, 3]
            },
            {
                "id": "monster2",
                "type": "monster",
                "name": "Orc",
                "initiative": 10,
                "hp": 30,
                "max_hp": 30,
                "ac": 14,
                "position": [7, 4]
            }
        ],
        "current_turn_index": 0,
        "round_number": 1,
        "status": "active",
        "environment": data.get('environment', 'forest'),
        "battlefield": {
            "width": 10,
            "height": 10,
            "terrain": [["normal" for _ in range(10)] for _ in range(10)],
            "objects": [["none" for _ in range(10)] for _ in range(10)]
        }
    }
    
    # Add some terrain features
    combat["battlefield"]["objects"][1][1] = "tree"
    combat["battlefield"]["objects"][3][4] = "tree"
    combat["battlefield"]["objects"][6][7] = "tree"
    combat["battlefield"]["terrain"][2][2] = "difficult"
    combat["battlefield"]["terrain"][2][3] = "difficult"
    
    return jsonify(combat), 201

@app.route('/api/v1/combat/<combat_id>', methods=['GET'])
def get_combat(combat_id):
    # Return the same mock combat data as the start_combat function
    combat = {
        "id": combat_id,
        "initiative": ["char1", "monster1", "char2", "monster2"],
        "participants": [
            {
                "id": "char1",
                "type": "character",
                "name": "Aragorn",
                "initiative": 18,
                "hp": 75,
                "max_hp": 80,
                "ac": 16,
                "position": [2, 3]
            },
            {
                "id": "char2",
                "type": "character",
                "name": "Gimli",
                "initiative": 12,
                "hp": 85,
                "max_hp": 85,
                "ac": 18,
                "position": [2, 4]
            },
            {
                "id": "monster1",
                "type": "monster",
                "name": "Goblin",
                "initiative": 15,
                "hp": 15,
                "max_hp": 15,
                "ac": 13,
                "position": [7, 3]
            },
            {
                "id": "monster2",
                "type": "monster",
                "name": "Orc",
                "initiative": 10,
                "hp": 30,
                "max_hp": 30,
                "ac": 14,
                "position": [7, 4]
            }
        ],
        "current_turn_index": 0,
        "round_number": 1,
        "status": "active",
        "environment": "forest",
        "battlefield": {
            "width": 10,
            "height": 10,
            "terrain": [["normal" for _ in range(10)] for _ in range(10)],
            "objects": [["none" for _ in range(10)] for _ in range(10)]
        }
    }
    
    # Add some terrain features
    combat["battlefield"]["objects"][1][1] = "tree"
    combat["battlefield"]["objects"][3][4] = "tree"
    combat["battlefield"]["objects"][6][7] = "tree"
    combat["battlefield"]["terrain"][2][2] = "difficult"
    combat["battlefield"]["terrain"][2][3] = "difficult"
    
    return jsonify(combat)

@app.route('/api/v1/combat/<combat_id>/action', methods=['POST'])
def perform_action(combat_id):
    data = request.get_json()
    
    action_type = data.get('action_type', '')
    actor_id = data.get('actor_id', '')
    target_ids = data.get('target_ids', [])
    
    # Process different action types
    result = {}
    
    if action_type == "attack":
        result = {
            "success": True,
            "damage": 12,
            "description": "Aragorn hits Goblin with longsword for 12 damage!"
        }
    elif action_type == "cast_spell":
        result = {
            "success": True,
            "damage": 18,
            "effects": ["stunned"],
            "description": "Wizard casts Thunderwave at Orc for 18 damage! Orc is stunned for 1 round."
        }
    elif action_type == "move":
        result = {
            "success": True,
            "description": "Aragorn moves from [2,3] to [3,4]"
        }
    else:
        result = {
            "success": False,
            "description": "Unknown action type"
        }
    
    # Get updated combat state
    combat = {
        "id": combat_id,
        "initiative": ["char1", "monster1", "char2", "monster2"],
        "participants": [
            {
                "id": "char1",
                "type": "character",
                "name": "Aragorn",
                "initiative": 18,
                "hp": 75,
                "max_hp": 80,
                "ac": 16,
                "position": [3, 4] if action_type == "move" and actor_id == "char1" else [2, 3]
            },
            {
                "id": "monster1",
                "type": "monster",
                "name": "Goblin",
                "initiative": 15,
                "hp": 3 if action_type == "attack" and "monster1" in target_ids else 15,
                "max_hp": 15,
                "ac": 13,
                "position": [7, 3]
            }
        ],
        "current_turn_index": 0,
        "round_number": 1,
        "status": "active"
    }
    
    return jsonify({
        "action_result": result,
        "combat": combat
    })

if __name__ == '__main__':
    port = int(os.environ.get('PORT', 5000))
    app.run(host='0.0.0.0', port=port)