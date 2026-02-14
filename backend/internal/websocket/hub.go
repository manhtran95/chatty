package websocket

import (
	"encoding/json"
	"log"

	messageprocessor "chatty.mtran.io/internal/message_processor"
	"github.com/google/uuid"
)

// Hub manages all WebSocket connections
type Hub struct {
	UserClients      map[uuid.UUID]*Client // user ID -> client (one client per user)
	Register         chan *Client
	Unregister       chan *Client
	Broadcast        chan []byte
	MessageProcessor *messageprocessor.MessageProcessor
}

// NewHub creates a new WebSocket hub
func NewHub(messageProcessor *messageprocessor.MessageProcessor) *Hub {
	return &Hub{
		UserClients:      make(map[uuid.UUID]*Client),
		Register:         make(chan *Client),
		Unregister:       make(chan *Client),
		Broadcast:        make(chan []byte),
		MessageProcessor: messageProcessor,
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

		case bytes := <-h.Broadcast:
			log.Printf("h.Broadcast channel received message")
			h.MessageProcessor.ProcessMessage(bytes)
		}
	}
}

// serializeResponse converts a Message to JSON bytes
func (h *Hub) serializeResponse(response *messageprocessor.Response) []byte {
	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error serializing message: %v", err)
		return nil
	}
	return data
}

// GetClientByUserID returns a client by user ID
func (h *Hub) GetClientByUserID(userID uuid.UUID) (*Client, bool) {
	client, exists := h.UserClients[userID]
	return client, exists
}

// SendToUser sends a message to a specific user
func (h *Hub) SendToUser(userID uuid.UUID, response *messageprocessor.Response) {
	if client, exists := h.UserClients[userID]; exists {
		serialized := h.serializeResponse(response)
		select {
		case client.Send <- serialized:
		default:
			log.Printf("Failed to send message to user %s: channel full", userID)
		}
	}
}
