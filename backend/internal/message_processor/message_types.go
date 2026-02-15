package messageprocessor

import "time"

// Message types that are read from client (incoming messages)
const (
	USER_CREATE_CHAT_REQUEST      = "UserCreateChatRequest"
	USER_SEND_MESSAGE_REQUEST     = "UserSendMessageRequest"
	USER_GET_CHAT_HISTORY_REQUEST = "UserGetChatHistoryRequest"
	USER_GET_CHATS_REQUEST        = "UserGetChatsRequest"
)

// Message types that are written to client (outgoing messages)
const (
	USER_CREATE_CHAT_RESPONSE      = "UserCreateChatResponse"
	USER_SEND_MESSAGE_RESPONSE     = "UserSendMessageResponse"
	USER_GET_CHAT_HISTORY_RESPONSE = "UserGetChatHistoryResponse"
	USER_GET_CHATS_RESPONSE        = "UserGetChatsResponse"
)

// MessageData represents the data payload for any message type
type MessageData interface {
	GetType() string
}

// Base types
// Base request structure
type Request struct {
	Type string      `json:"type"`
	Data MessageData `json:"data"`
}

// Base response structure
type Response struct {
	Type  string      `json:"type"`
	Data  MessageData `json:"data"`
	Error string      `json:"error"`
}

type UserInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type ChatInfo struct {
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	UpdatedAt time.Time  `json:"updatedAt"`
	UserInfos []UserInfo `json:"userInfos"`
}

type ChatHistoryMessage struct {
	MessageID  string `json:"messageId"`
	SenderName string `json:"senderName"`
	Content    string `json:"content"`
	Timestamp  string `json:"timestamp"`
}

// Client APIs
//
// Send Message
type ClientSendMessageRequest struct {
	ChatID  string `json:"chatId"`
	Content string `json:"content"`
}

func (d ClientSendMessageRequest) GetType() string { return USER_SEND_MESSAGE_REQUEST }

type ClientSendMessageResponse struct {
	ChatID    string `json:"chatId"`
	SenderID  string `json:"senderId"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
	MessageID string `json:"messageId"`
}

func (d ClientSendMessageResponse) GetType() string { return USER_SEND_MESSAGE_RESPONSE }

// Create Chat
type CreateChatRequest struct {
	Name              string   `json:"name"`
	ParticipantEmails []string `json:"participantEmails"`
}

func (d CreateChatRequest) GetType() string { return USER_CREATE_CHAT_REQUEST }

type CreateChatResponse ChatInfo

func (d CreateChatResponse) GetType() string { return USER_CREATE_CHAT_RESPONSE }

// Get Chat History
type GetChatHistoryRequest struct {
	ChatID string `json:"chatID"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}

func (d GetChatHistoryRequest) GetType() string { return USER_GET_CHAT_HISTORY_REQUEST }

type GetChatHistoryResponse struct {
	ChatID   string               `json:"chatID"`
	Messages []ChatHistoryMessage `json:"messages"`
	HasMore  bool                 `json:"hasMore"`
}

func (d GetChatHistoryResponse) GetType() string { return USER_GET_CHAT_HISTORY_RESPONSE }

// Get Chats
type GetChatsRequest struct {
	CursorUpdatedAt *time.Time `json:"cursorUpdatedAt,omitempty"`
	ChatID          *string    `json:"chatID,omitempty"`
	Limit           int        `json:"limit"`
}

func (d GetChatsRequest) GetType() string { return USER_GET_CHATS_REQUEST }

type GetChatsResponse struct {
	Chats []ChatInfo `json:"chats"`
}

func (d GetChatsResponse) GetType() string { return USER_GET_CHATS_RESPONSE }

// NewMessage creates a new message with the given data
func NewMessage(data MessageData) *Request {
	return &Request{
		Type: data.GetType(),
		Data: data,
	}
}

// IsIncomingMessage checks if a message type is incoming (from client)
func IsIncomingMessage(messageType string) bool {
	switch messageType {
	case USER_SEND_MESSAGE_REQUEST, USER_CREATE_CHAT_REQUEST, USER_GET_CHAT_HISTORY_REQUEST:
		return true
	default:
		return false
	}
}

// IsOutgoingMessage checks if a message type is outgoing (to client)
func IsOutgoingMessage(messageType string) bool {
	switch messageType {
	case USER_SEND_MESSAGE_RESPONSE, USER_CREATE_CHAT_RESPONSE, USER_GET_CHAT_HISTORY_RESPONSE:
		return true
	default:
		return false
	}
}
