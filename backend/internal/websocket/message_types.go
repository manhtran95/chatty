package websocket

// Message types that are read from client (incoming messages)
const (
	CLIENT_SEND_MESSAGE         = "ClientSendMessage"
	CLIENT_CREATE_CHAT          = "ClientCreateChat"
	CLIENT_REQUEST_CHAT_HISTORY = "ClientRequestChatHistory"
	CLIENT_REQUEST_ALL_CHATS    = "ClientRequestAllChats"
)

// Message types that are written to client (outgoing messages)
const (
	CLIENT_RECEIVE_MESSAGE      = "ClientReceiveMessage"
	CLIENT_RECEIVE_CHAT         = "ClientReceiveChat"
	CLIENT_RECEIVE_CHAT_HISTORY = "ClientReceiveChatHistory"
	CLIENT_RECEIVE_ALL_CHATS    = "ClientReceiveAllChats"
)

// MessageData represents the data payload for any message type
type MessageData interface {
	GetType() string
}

// Base message structure
type Message struct {
	Type     string      `json:"type"`
	Data     MessageData `json:"data"`
	SenderID string      `json:"senderId"`
}

// ClientSendMessageData represents data for sending a message
type ClientSendMessageData struct {
	ChatID   string `json:"chatId"`
	SenderID string `json:"senderId"`
	Content  string `json:"content"`
}

func (d ClientSendMessageData) GetType() string { return CLIENT_SEND_MESSAGE }

// ClientReceiveMessageData represents data for receiving a message
type ClientReceiveMessageData struct {
	ChatID     string `json:"chatId"`
	SenderName string `json:"senderName"`
	Content    string `json:"content"`
	Timestamp  string `json:"timestamp"`
	MessageID  string `json:"messageId"`
}

func (d ClientReceiveMessageData) GetType() string { return CLIENT_RECEIVE_MESSAGE }

// UserCreateChatRequest represents data for creating a chat
type UserCreateChatRequest struct {
	Name              string   `json:"name"`
	ParticipantEmails []string `json:"participantEmails"`
}

func (d UserCreateChatRequest) GetType() string { return CLIENT_CREATE_CHAT }

type UserInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// UserCreateChatResponse represents data for receiving a chat
type UserCreateChatResponse struct {
	ChatID           string     `json:"chatID"`
	Name             string     `json:"name"`
	ParticipantInfos []UserInfo `json:"participantInfos"`
}

func (d UserCreateChatResponse) GetType() string { return CLIENT_RECEIVE_CHAT }

// ClientRequestChatHistoryData represents data for requesting chat history
type ClientRequestChatHistoryData struct {
	ChatID string `json:"chatID"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}

func (d ClientRequestChatHistoryData) GetType() string { return CLIENT_REQUEST_CHAT_HISTORY }

// ChatHistoryMessage represents a single message in chat history
type ChatHistoryMessage struct {
	MessageID  string `json:"messageId"`
	SenderName string `json:"senderName"`
	Content    string `json:"content"`
	Timestamp  string `json:"timestamp"`
}

// ClientReceiveChatHistoryData represents data for receiving chat history
type ClientReceiveChatHistoryData struct {
	ChatID   string               `json:"chatID"`
	Messages []ChatHistoryMessage `json:"messages"`
	HasMore  bool                 `json:"hasMore"`
}

func (d ClientReceiveChatHistoryData) GetType() string { return CLIENT_RECEIVE_CHAT_HISTORY }

type ClientRequestAllChatsData struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

func (d ClientRequestAllChatsData) GetType() string { return CLIENT_REQUEST_ALL_CHATS }

type ClientReceiveAllChatsData struct {
	Chats []UserCreateChatResponse `json:"chats"`
}

func (d ClientReceiveAllChatsData) GetType() string { return CLIENT_RECEIVE_ALL_CHATS }

// NewMessage creates a new message with the given data
func NewMessage(data MessageData) *Message {
	return &Message{
		Type: data.GetType(),
		Data: data,
	}
}

// IsIncomingMessage checks if a message type is incoming (from client)
func IsIncomingMessage(messageType string) bool {
	switch messageType {
	case CLIENT_SEND_MESSAGE, CLIENT_CREATE_CHAT, CLIENT_REQUEST_CHAT_HISTORY:
		return true
	default:
		return false
	}
}

// IsOutgoingMessage checks if a message type is outgoing (to client)
func IsOutgoingMessage(messageType string) bool {
	switch messageType {
	case CLIENT_RECEIVE_MESSAGE, CLIENT_RECEIVE_CHAT, CLIENT_RECEIVE_CHAT_HISTORY:
		return true
	default:
		return false
	}
}
