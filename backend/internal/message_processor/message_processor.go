package messageprocessor

import (
	"encoding/json"
	"errors"
	"log"

	"chatty.mtran.io/internal/models"
	"github.com/google/uuid"
)

type MessageProcessor struct {
	ChatModel     *models.ChatModel
	UserModel     *models.UserModel
	MessageModel  *models.MessageModel
	MessageSender ResponseSender
}

type rawMessage struct {
	Type     string          `json:"type"`
	Data     json.RawMessage `json:"data"`
	SenderID string          `json:"senderId"`
}

type ResponseSender interface {
	SendToUser(userID uuid.UUID, response *Response)
}

// constructors and setters
func NewMessageProcessor(chatModel *models.ChatModel, userModel *models.UserModel, messageModel *models.MessageModel) *MessageProcessor {
	return &MessageProcessor{
		ChatModel:    chatModel,
		UserModel:    userModel,
		MessageModel: messageModel,
	}
}

func (mp *MessageProcessor) SetMessageSender(messageSender ResponseSender) {
	mp.MessageSender = messageSender
}

func unmarshalRawMessage(p []byte) (message *Request, err error) {
	var raw rawMessage
	if err := json.Unmarshal(p, &raw); err != nil {
		return nil, err
	}

	var data MessageData
	switch raw.Type {
	case USER_CREATE_CHAT_REQUEST:
		log.Printf("CLIENT_CREATE_CHAT: %v", raw.Data)

		var createChatData *UserCreateChatRequest
		if err := json.Unmarshal(raw.Data, &createChatData); err != nil {
			return nil, err
		}
		data = createChatData

	case CLIENT_SEND_MESSAGE_REQUEST:
		var sendMessageData *ClientSendMessageRequest
		if err := json.Unmarshal(raw.Data, &sendMessageData); err != nil {
			return nil, err
		}
		data = sendMessageData

	default:
		log.Printf("unknown message type: %s", raw.Type)
		return &Request{
			Type:     raw.Type,
			Data:     nil,
			SenderID: raw.SenderID,
		}, errors.New("unknown message type")
	}

	return &Request{
		Type:     raw.Type,
		Data:     data,
		SenderID: raw.SenderID,
	}, nil
}

// functions
func (mp *MessageProcessor) ProcessMessage(rawMessage []byte) {
	message, err := unmarshalRawMessage(rawMessage)
	senderId := uuid.MustParse(message.SenderID)
	if err != nil {
		log.Printf("Error unmarshalling message: %v", err)
		log.Printf("senderId: %s", senderId)
		responseMessage := &Response{
			Type:  "",
			Data:  nil,
			Error: ErrUnknownRequestType.Error(),
		}
		mp.MessageSender.SendToUser(senderId, responseMessage)
		return
	}

	switch message.Type {
	case USER_CREATE_CHAT_REQUEST:
		mp.handleUserCreateChatRequest(message)
	case CLIENT_SEND_MESSAGE_REQUEST:
		mp.handleClientSendMessageRequest(message)
	case CLIENT_GET_CHAT_HISTORY_REQUEST:
		mp.handleClientGetChatHistoryRequest(message)
	}
}

func (mp *MessageProcessor) handleUserCreateChatRequest(message *Request) {
	chatData, ok := message.Data.(*UserCreateChatRequest)
	if !ok {
		log.Printf("Invalid chat creation data")
		return
	}

	senderId := uuid.MustParse(message.SenderID)

	// check if participantIDs are valid
	if len(chatData.ParticipantEmails) < 2 {
		log.Printf("Cannot create a chat with less than 2 participants")
		responseMessage := &Response{
			Type:  USER_CREATE_CHAT_RESPONSE,
			Data:  nil,
			Error: ErrCannotCreateChatWithLessThan2Participants.Error(),
		}
		mp.MessageSender.SendToUser(senderId, responseMessage)
		return
	}

	userInfos, err := mp.UserModel.UserInfosByEmails(chatData.ParticipantEmails)
	if err != nil {
		log.Printf("Error checking if user emails exist: %v", err)
		responseMessage := &Response{
			Type:  USER_CREATE_CHAT_RESPONSE,
			Data:  nil,
			Error: ErrUserDoesNotExist.Error(),
		}
		mp.MessageSender.SendToUser(senderId, responseMessage)
		return
	}

	// TO DO: check if set of users already in a chat

	// Create chat in database
	chat, err := mp.ChatModel.InsertChat(chatData.Name)
	if err != nil {
		log.Printf("Error creating chat: %v", err)
		responseMessage := &Response{
			Type:  USER_CREATE_CHAT_RESPONSE,
			Data:  nil,
			Error: ErrCannotCreateChat.Error(),
		}
		mp.MessageSender.SendToUser(senderId, responseMessage)
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
		responseMessage := &Response{
			Type:  USER_CREATE_CHAT_RESPONSE,
			Data:  nil,
			Error: ErrCannotAddParticipantsToChat.Error(),
		}
		mp.MessageSender.SendToUser(senderId, responseMessage)
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

	responseMessage := &Response{
		Type:  USER_CREATE_CHAT_RESPONSE,
		Data:  responseData,
		Error: "",
	}

	for _, userID := range userIDs {
		mp.MessageSender.SendToUser(userID, responseMessage)
	}
}

func (mp *MessageProcessor) handleClientSendMessageRequest(request *Request) {
	messageData, ok := request.Data.(*ClientSendMessageRequest)
	if !ok {
		log.Printf("Invalid message data")
		return
	}
	chatId := uuid.MustParse(messageData.ChatID)
	senderId := uuid.MustParse(request.SenderID)
	content := messageData.Content

	// get chat by id
	_, err := mp.ChatModel.GetChat(chatId)
	if err != nil {
		log.Printf("Error getting chat: %v", err)
		responseMessage := &Response{
			Type:  CLIENT_SEND_MESSAGE_RESPONSE,
			Data:  nil,
			Error: ErrCannotGetChat.Error(),
		}
		mp.MessageSender.SendToUser(senderId, responseMessage)
		return
	}

	// Save message to database
	message, err := mp.MessageModel.InsertMessage(senderId, chatId, content)
	if err != nil {
		log.Printf("Error saving message: %v", err)
		responseMessage := &Response{
			Type:  CLIENT_SEND_MESSAGE_RESPONSE,
			Data:  nil,
			Error: ErrCannotSaveMessage.Error(),
		}
		mp.MessageSender.SendToUser(senderId, responseMessage)
		return
	}

	// response to chat members
	members, err := mp.ChatModel.GetChatMembers(chatId)
	if err != nil {
		log.Printf("Error getting chat members: %v", err)
		return
	}
	memberIds := make([]uuid.UUID, 0, len(members))
	for _, member := range members {
		memberIds = append(memberIds, member.ID)
	}
	responseMessage := &Response{
		Type: CLIENT_SEND_MESSAGE_RESPONSE,
		Data: ClientSendMessageResponse{
			ChatID:    chatId.String(),
			SenderID:  message.SenderID.String(),
			Content:   message.Content,
			Timestamp: message.CreatedAt.String(),
			MessageID: message.ID.String(),
		},
		Error: "",
	}

	for _, memberId := range memberIds {
		mp.MessageSender.SendToUser(memberId, responseMessage)
	}
}

func (mp *MessageProcessor) handleClientGetChatHistoryRequest(message *Request) {
	// messageData, ok := message.Data.(*ClientGetChatHistoryRequest)
	// if !ok {
	// 	log.Printf("Invalid message data")
	// 	return
	// }
}
