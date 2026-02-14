package websocket

import (
	"encoding/json"
	"log"

	"chatty.mtran.io/internal/models"
	"github.com/google/uuid"
)

// Hub manages all WebSocket connections
type Hub struct {
	UserClients  map[uuid.UUID]*Client // user ID -> client (one client per user)
	Register     chan *Client
	Unregister   chan *Client
	Broadcast    chan *Message
	ChatModel    *models.ChatModel    // Add ChatModel for database operations
	UserModel    *models.UserModel    // Add UserModel for user validation
	MessageModel *models.MessageModel // Add MessageModel for message operations
}

// NewHub creates a new WebSocket hub
func NewHub(chatModel *models.ChatModel, userModel *models.UserModel, messageModel *models.MessageModel) *Hub {
	return &Hub{
		UserClients:  make(map[uuid.UUID]*Client),
		Register:     make(chan *Client),
		Unregister:   make(chan *Client),
		Broadcast:    make(chan *Message),
		ChatModel:    chatModel,
		UserModel:    userModel,
		MessageModel: messageModel,
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
func (h *Hub) handleMessage(message *Message) {
	switch message.Type {
	case CLIENT_SEND_MESSAGE:
		h.handleSendMessage(message)
	case CLIENT_CREATE_CHAT:
		h.handleCreateChat(message)
	case CLIENT_REQUEST_CHAT_HISTORY:
		h.handleRequestChatHistory(message)
	default:
		log.Printf("Unknown message type: %s", message.Type)
	}
}

// handleSendMessage processes a message sent by a client
func (h *Hub) handleSendMessage(message *Message) {
	log.Printf("handleSendMessage: %v", message)
	// TODO: Save message to database
	// TODO: Broadcast to all users in the chat

	// For now, just broadcast to all clients
	h.broadcastToChat(message)
}

// handleCreateChat processes a chat creation request
func (h *Hub) handleCreateChat(message *Message) {
	log.Printf("handleCreateChat")
	// Type assert to get the chat data
	chatData, ok := message.Data.(*UserCreateChatRequest)
	if !ok {
		log.Printf("Invalid chat creation data")
		return
	}

	// check if participantIDs are valid
	if len(chatData.ParticipantEmails) < 2 {
		log.Printf("Cannot create a chat with less than 2 participants")
		return
	}

	userInfos, err := h.UserModel.UserInfosByEmails(chatData.ParticipantEmails)
	if err != nil {
		log.Printf("Error checking if user emails exist: %v", err)
		return
	}

	// Create chat in database
	chat, err := h.ChatModel.InsertChat(chatData.Name)
	if err != nil {
		log.Printf("Error creating chat: %v", err)
		// TODO: Send error response to client
		return
	}

	// Add participants to chat
	userIDs := make([]uuid.UUID, 0, len(userInfos))
	for _, userInfo := range userInfos {
		userIDs = append(userIDs, userInfo.ID)
	}
	err = h.ChatModel.AddUsersToChat(chat.ID, userIDs)
	if err != nil {
		log.Printf("Error adding participants to chat: %v", err)
		return
	}
	log.Printf("Created chat: %s with ID: %s", chat.Name, chat.ID)

	// Create response message
	wsUserInfos := make([]UserInfo, 0, len(userInfos))
	for _, userInfo := range userInfos {
		wsUserInfos = append(wsUserInfos, UserInfo{
			ID:    userInfo.ID.String(),
			Email: userInfo.Email,
			Name:  userInfo.Name,
		})
	}
	responseData := &UserCreateChatResponse{
		ChatID:           chat.ID.String(),
		Name:             chat.Name,
		ParticipantInfos: wsUserInfos,
	}

	responseMessage := NewMessage(responseData)

	for _, userID := range userIDs {
		h.SendToUser(userID, responseMessage)
	}
}

// handleRequestChatHistory processes a request for chat history
func (h *Hub) handleRequestChatHistory(message *Message) {
	// TODO: Fetch chat history from database
	// TODO: Send messages to requesting user

	log.Printf("User requested chat history")
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
func (h *Hub) GetClientByUserID(userID uuid.UUID) (*Client, bool) {
	client, exists := h.UserClients[userID]
	return client, exists
}

// SendToUser sends a message to a specific user
func (h *Hub) SendToUser(userID uuid.UUID, message *Message) {
	if client, exists := h.UserClients[userID]; exists {
		serialized := h.serializeMessage(message)
		select {
		case client.Send <- serialized:
		default:
			log.Printf("Failed to send message to user %s: channel full", userID)
		}
	}
}
