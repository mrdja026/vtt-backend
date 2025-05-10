package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period (must be less than pongWait)
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512 * 1024
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins for now (in production, you'd want to restrict this)
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Message represents a message sent through the websocket
type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// Client represents a connected websocket client
type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	roomID   string
	userID   string
	mu       sync.Mutex
	isClosed bool
}

// Hub maintains the set of active clients and broadcasts messages to clients
type Hub struct {
	// Registered clients by room ID
	rooms      map[string]map[*Client]bool
	roomsMu    sync.RWMutex
	
	// Map of user IDs to clients
	users      map[string]*Client
	usersMu    sync.RWMutex
	
	// Register requests from clients
	register   chan *Client
	
	// Unregister requests from clients
	unregister chan *Client
}

// NewHub creates a new hub for websocket connections
func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]map[*Client]bool),
		users:      make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			// Add client to the room
			h.roomsMu.Lock()
			if _, ok := h.rooms[client.roomID]; !ok {
				h.rooms[client.roomID] = make(map[*Client]bool)
			}
			h.rooms[client.roomID][client] = true
			h.roomsMu.Unlock()
			
			// Add client to users map
			h.usersMu.Lock()
			h.users[client.userID] = client
			h.usersMu.Unlock()
			
			// Send a welcome message
			client.send <- []byte(fmt.Sprintf(`{"type":"connect","data":{"message":"Connected to combat session %s"}}`, client.roomID))
			
		case client := <-h.unregister:
			// Remove client from room
			h.roomsMu.Lock()
			if _, ok := h.rooms[client.roomID]; ok {
				delete(h.rooms[client.roomID], client)
				// If room is empty, remove it
				if len(h.rooms[client.roomID]) == 0 {
					delete(h.rooms, client.roomID)
				}
			}
			h.roomsMu.Unlock()
			
			// Remove client from users map
			h.usersMu.Lock()
			if h.users[client.userID] == client {
				delete(h.users, client.userID)
			}
			h.usersMu.Unlock()
			
			// Close client's send channel
			client.mu.Lock()
			if !client.isClosed {
				close(client.send)
				client.isClosed = true
			}
			client.mu.Unlock()
		}
	}
}

// BroadcastToRoom sends a message to all clients in a room
func (h *Hub) BroadcastToRoom(roomID string, message Message) {
	// Marshal the message to JSON
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}
	
	// Get all clients in the room
	h.roomsMu.RLock()
	clients, ok := h.rooms[roomID]
	h.roomsMu.RUnlock()
	
	if !ok {
		// Room doesn't exist
		return
	}
	
	// Send message to all clients in the room
	for client := range clients {
		client.mu.Lock()
		if !client.isClosed {
			select {
			case client.send <- data:
				// Message sent
			default:
				// Buffer full, remove client
				h.unregister <- client
			}
		}
		client.mu.Unlock()
	}
}

// SendToUser sends a message to a specific user
func (h *Hub) SendToUser(userID string, message Message) {
	// Marshal the message to JSON
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}
	
	// Get the client for this user
	h.usersMu.RLock()
	client, ok := h.users[userID]
	h.usersMu.RUnlock()
	
	if !ok {
		// User not connected
		return
	}
	
	// Send message to the client
	client.mu.Lock()
	if !client.isClosed {
		select {
		case client.send <- data:
			// Message sent
		default:
			// Buffer full, remove client
			h.unregister <- client
		}
	}
	client.mu.Unlock()
}

// readPump reads messages from the websocket connection and handles them
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}
		
		// For now, we don't process incoming messages - the client just receives updates
		// You could add message handling here if needed
	}
}

// writePump sends messages from the hub to the websocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			
			// Add queued messages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}
			
			if err := w.Close(); err != nil {
				return
			}
			
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWs handles websocket connections from clients
func (h *Hub) ServeWs(w http.ResponseWriter, r *http.Request, roomID string, userID string) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}
	
	// Create new client
	client := &Client{
		hub:      h,
		conn:     conn,
		send:     make(chan []byte, 256),
		roomID:   roomID,
		userID:   userID,
		isClosed: false,
	}
	
	// Register client
	h.register <- client
	
	// Start goroutines for reading and writing
	go client.readPump()
	go client.writePump()
}
