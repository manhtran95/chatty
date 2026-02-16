package messageprocessor

import (
	"encoding/json"
	"time"
)

// Message types that are read from client (incoming messages)
const (
	CREATE_CHAT_REQUEST      = "CreateChatRequest"
	SEND_MESSAGE_REQUEST     = "SendMessageRequest"
	GET_CHAT_HISTORY_REQUEST = "GetChatHistoryRequest"
	GET_CHATS_REQUEST        = "GetChatsRequest"
)

// Message types that are written to client (outgoing messages)
const (
	CREATE_CHAT_RESPONSE      = "CreateChatResponse"
	SEND_MESSAGE_RESPONSE     = "SendMessageResponse"
	GET_CHAT_HISTORY_RESPONSE = "GetChatHistoryResponse"
	GET_CHATS_RESPONSE        = "GetChatsResponse"
)

// Base types
// Base request structure
type Request struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// Base response structure
type Response struct {
	Type  string          `json:"type"`
	Data  json.RawMessage `json:"data"`
	Error string          `json:"error"`
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
type SendMessageRequest struct {
	ChatID  string `json:"chatId"`
	Content string `json:"content"`
}

type SendMessageResponse struct {
	ChatID    string `json:"chatId"`
	SenderID  string `json:"senderId"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
	MessageID string `json:"messageId"`
}

// Create Chat
type CreateChatRequest struct {
	Name              string   `json:"name"`
	ParticipantEmails []string `json:"participantEmails"`
}

type CreateChatResponse ChatInfo

// Get Chat History
type GetChatHistoryRequest struct {
	ChatID string `json:"chatID"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}

type GetChatHistoryResponse struct {
	ChatID   string               `json:"chatID"`
	Messages []ChatHistoryMessage `json:"messages"`
	HasMore  bool                 `json:"hasMore"`
}

// Get Chats
type GetChatsRequest struct {
	CursorUpdatedAt *time.Time `json:"cursorUpdatedAt,omitempty"`
	ChatID          *string    `json:"chatID,omitempty"`
	Limit           int        `json:"limit"`
}

type GetChatsResponse struct {
	Chats []ChatInfo `json:"chats"`
}
