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
	Broadcast        chan *messageprocessor.Message
	MessageProcessor *messageprocessor.MessageProcessor
}

// NewHub creates a new WebSocket hub
func NewHub(messageProcessor *messageprocessor.MessageProcessor) *Hub {
	return &Hub{
		UserClients:      make(map[uuid.UUID]*Client),
		Register:         make(chan *Client),
		Unregister:       make(chan *Client),
		Broadcast:        make(chan *messageprocessor.Message),
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

		case message := <-h.Broadcast:
			log.Printf("h.Broadcast channel received message: %v", message)
			h.handleMessage(message)
		}
	}
}

// handleMessage processes incoming messages
func (h *Hub) handleMessage(message *messageprocessor.Message) {
	switch message.Type {
	case messageprocessor.CLIENT_SEND_MESSAGE_REQUEST:
		h.handleSendMessage(message)
	case messageprocessor.USER_CREATE_CHAT_REQUEST:
		h.handleCreateChat(message)
	case messageprocessor.CLIENT_GET_CHAT_HISTORY_REQUEST:
		h.handleRequestChatHistory(message)
	default:
		log.Printf("Unknown message type: %s", message.Type)
	}
}

// handleSendMessage processes a message sent by a client
func (h *Hub) handleSendMessage(message *messageprocessor.Message) {
	log.Printf("handleSendMessage: %v", message)
	// TODO: Save message to database
	// TODO: Broadcast to all users in the chat

	// For now, just broadcast to all clients
	h.broadcastToChat(message)
}

// handleCreateChat processes a chat creation request
func (h *Hub) handleCreateChat(message *messageprocessor.Message) {
	log.Printf("handleCreateChat")
	// Type assert to get the chat data
	h.MessageProcessor.ProcessMessage(message)
}

// handleRequestChatHistory processes a request for chat history
func (h *Hub) handleRequestChatHistory(message *messageprocessor.Message) {
	// TODO: Fetch chat history from database
	// TODO: Send messages to requesting user

	log.Printf("User requested chat history")
}

// broadcastToChat sends a message to all users in a specific chat
func (h *Hub) broadcastToChat(message *messageprocessor.Message) {
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

// serializeMessage converts a Message to JSON bytes
func (h *Hub) serializeMessage(message *messageprocessor.Message) []byte {
	data, err := json.Marshal(message)
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
func (h *Hub) SendToUser(userID uuid.UUID, message *messageprocessor.Message) {
	if client, exists := h.UserClients[userID]; exists {
		serialized := h.serializeMessage(message)
		select {
		case client.Send <- serialized:
		default:
			log.Printf("Failed to send message to user %s: channel full", userID)
		}
	}
}
