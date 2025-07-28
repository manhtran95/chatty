package websocket

import (
	"encoding/json"
	"log"
)

// Message represents a WebSocket message
type Message struct {
	Type      string `json:"type"`
	SenderID  string `json:"sender_id"`
	ChatID    string `json:"chat_id"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
	Data      any    `json:"data,omitempty"`
}

// Hub manages all WebSocket connections
type Hub struct {
	UserClients map[string]*Client // user ID -> client (one client per user)
	Register    chan *Client
	Unregister  chan *Client
	Broadcast   chan *Message
}

// NewHub creates a new WebSocket hub
func NewHub() *Hub {
	return &Hub{
		UserClients: make(map[string]*Client),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Broadcast:   make(chan *Message),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.UserClients[client.UserID] = client
			log.Printf("Client registered: User %s", client.UserID)

		case client := <-h.Unregister:
			if _, ok := h.UserClients[client.UserID]; ok {
				delete(h.UserClients, client.UserID)
				close(client.Send)
			}
			log.Printf("Client unregistered: User %s", client.UserID)

		case message := <-h.Broadcast:
			h.handleMessage(message)
		}
	}
}

// handleMessage processes incoming messages
func (h *Hub) handleMessage(message *Message) {
	switch message.Type {
	case "chat_message":
		h.broadcastToChat(message)
	case "join_chat":
		h.handleJoinChat(message)
	case "leave_chat":
		h.handleLeaveChat(message)
	default:
		log.Printf("Unknown message type: %s", message.Type)
	}
}

// broadcastToChat sends a message to all users in a specific chat
func (h *Hub) broadcastToChat(message *Message) {
	// For now, broadcast to all clients (you can implement chat-specific logic later)
	for _, client := range h.UserClients {
		select {
		case client.Send <- h.serializeMessage(message):
		default:
			close(client.Send)
			delete(h.UserClients, client.UserID)
		}
	}
}

// handleJoinChat handles when a user joins a chat
func (h *Hub) handleJoinChat(message *Message) {
	log.Printf("User %s joined chat %s", message.SenderID, message.ChatID)
	// You can implement chat-specific logic here
}

// handleLeaveChat handles when a user leaves a chat
func (h *Hub) handleLeaveChat(message *Message) {
	log.Printf("User %s left chat %s", message.SenderID, message.ChatID)
	// You can implement chat-specific logic here
}

// serializeMessage converts a Message to JSON bytes
func (h *Hub) serializeMessage(message *Message) []byte {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error serializing message: %v", err)
		return nil
	}
	return data
}

// GetClientByUserID returns a client by user ID
func (h *Hub) GetClientByUserID(userID string) (*Client, bool) {
	client, exists := h.UserClients[userID]
	return client, exists
}

// SendToUser sends a message to a specific user
func (h *Hub) SendToUser(userID string, message *Message) {
	if client, exists := h.GetClientByUserID(userID); exists {
		select {
		case client.Send <- h.serializeMessage(message):
		default:
			log.Printf("Failed to send message to user %s", userID)
		}
	}
}
