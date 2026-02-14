package messageprocessor

// Message types that are read from client (incoming messages)
const (
	USER_CREATE_CHAT_REQUEST        = "UserCreateChatRequest"
	CLIENT_SEND_MESSAGE_REQUEST     = "ClientSendMessageRequest"
	CLIENT_GET_CHAT_HISTORY_REQUEST = "ClientGetChatHistoryRequest"
	CLIENT_GET_ALL_CHATS_REQUEST    = "ClientGetAllChatsRequest"
)

// Message types that are written to client (outgoing messages)
const (
	USER_CREATE_CHAT_RESPONSE        = "UserCreateChatResponse"
	CLIENT_SEND_MESSAGE_RESPONSE     = "ClientSendMessageResponse"
	CLIENT_GET_CHAT_HISTORY_RESPONSE = "ClientGetChatHistoryResponse"
	CLIENT_GET_ALL_CHATS_RESPONSE    = "ClientGetAllChatsResponse"
)

// MessageData represents the data payload for any message type
type MessageData interface {
	GetType() string
}

// Base request structure
type Request struct {
	Type     string      `json:"type"`
	Data     MessageData `json:"data"`
	SenderID string      `json:"senderId"`
}

// Base response structure
type Response struct {
	Type  string      `json:"type"`
	Data  MessageData `json:"data"`
	Error string      `json:"error"`
}

// ClientSendMessageData represents data for sending a message
type ClientSendMessageRequest struct {
	ChatID   string `json:"chatId"`
	Content  string `json:"content"`
}

func (d ClientSendMessageRequest) GetType() string { return CLIENT_SEND_MESSAGE_REQUEST }

// ClientReceiveMessageData represents data for receiving a message
type ClientSendMessageResponse struct {
	ChatID    string `json:"chatId"`
	SenderID  string `json:"senderId"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
	MessageID string `json:"messageId"`
}

func (d ClientSendMessageResponse) GetType() string { return CLIENT_SEND_MESSAGE_RESPONSE }

// UserCreateChatRequest represents data for creating a chat
type UserCreateChatRequest struct {
	Name              string   `json:"name"`
	ParticipantEmails []string `json:"participantEmails"`
}

func (d UserCreateChatRequest) GetType() string { return USER_CREATE_CHAT_REQUEST }

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

func (d UserCreateChatResponse) GetType() string { return USER_CREATE_CHAT_RESPONSE }

// ClientRequestChatHistoryData represents data for requesting chat history
type ClientGetChatHistoryRequest struct {
	ChatID string `json:"chatID"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}

func (d ClientGetChatHistoryRequest) GetType() string { return CLIENT_GET_CHAT_HISTORY_REQUEST }

// ChatHistoryMessage represents a single message in chat history
type ChatHistoryMessage struct {
	MessageID  string `json:"messageId"`
	SenderName string `json:"senderName"`
	Content    string `json:"content"`
	Timestamp  string `json:"timestamp"`
}

// ClientReceiveChatHistoryData represents data for receiving chat history
type ClientGetChatHistoryResponse struct {
	ChatID   string               `json:"chatID"`
	Messages []ChatHistoryMessage `json:"messages"`
	HasMore  bool                 `json:"hasMore"`
}

func (d ClientGetChatHistoryResponse) GetType() string { return CLIENT_GET_CHAT_HISTORY_RESPONSE }

type ClientGetAllChatsRequest struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

func (d ClientGetAllChatsRequest) GetType() string { return CLIENT_GET_ALL_CHATS_REQUEST }

type ClientGetAllChatsResponse struct {
	Chats []UserCreateChatResponse `json:"chats"`
}

func (d ClientGetAllChatsResponse) GetType() string { return CLIENT_GET_ALL_CHATS_RESPONSE }

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
	case CLIENT_SEND_MESSAGE_REQUEST, USER_CREATE_CHAT_REQUEST, CLIENT_GET_CHAT_HISTORY_REQUEST:
		return true
	default:
		return false
	}
}

// IsOutgoingMessage checks if a message type is outgoing (to client)
func IsOutgoingMessage(messageType string) bool {
	switch messageType {
	case CLIENT_SEND_MESSAGE_RESPONSE, USER_CREATE_CHAT_RESPONSE, CLIENT_GET_CHAT_HISTORY_RESPONSE:
		return true
	default:
		return false
	}
}
