package messageprocessor

import (
	"log"

	"chatty.mtran.io/internal/models"
	"github.com/google/uuid"
)

type MessageProcessor struct {
	ChatModel     *models.ChatModel
	UserModel     *models.UserModel
	MessageModel  *models.MessageModel
	MessageSender MessageSender
}

type MessageSender interface {
	SendToUser(userID uuid.UUID, message *Message)
}

func NewMessageProcessor(chatModel *models.ChatModel, userModel *models.UserModel, messageModel *models.MessageModel) *MessageProcessor {
	return &MessageProcessor{
		ChatModel:    chatModel,
		UserModel:    userModel,
		MessageModel: messageModel,
	}
}

func (mp *MessageProcessor) SetMessageSender(messageSender MessageSender) {
	mp.MessageSender = messageSender
}

func (mp *MessageProcessor) ProcessMessage(message *Message) {
	switch message.Type {
	case USER_CREATE_CHAT_REQUEST:
		mp.handleUserCreateChatRequest(message)
	case CLIENT_SEND_MESSAGE_REQUEST:
		mp.handleClientSendMessageRequest(message)
	case CLIENT_GET_CHAT_HISTORY_REQUEST:
		mp.handleClientGetChatHistoryRequest(message)
	}
}

func (mp *MessageProcessor) handleUserCreateChatRequest(message *Message) {
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

	userInfos, err := mp.UserModel.UserInfosByEmails(chatData.ParticipantEmails)
	if err != nil {
		log.Printf("Error checking if user emails exist: %v", err)
		return
	}

	// TO DO: check if set of users already in a chat

	// Create chat in database
	chat, err := mp.ChatModel.InsertChat(chatData.Name)
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
	err = mp.ChatModel.AddUsersToChat(chat.ID, userIDs)
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
		mp.MessageSender.SendToUser(userID, responseMessage)
	}
}

func (mp *MessageProcessor) handleClientSendMessageRequest(message *Message) {
	// messageData, ok := message.Data.(*ClientSendMessageRequest)
	// if !ok {
	// 	log.Printf("Invalid message data")
	// 	return
	// }
}

func (mp *MessageProcessor) handleClientGetChatHistoryRequest(message *Message) {
	// messageData, ok := message.Data.(*ClientGetChatHistoryRequest)
	// if !ok {
	// 	log.Printf("Invalid message data")
	// 	return
	// }
}
