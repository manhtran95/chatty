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

		var createChatData *CreateChatRequest
		if err := json.Unmarshal(raw.Data, &createChatData); err != nil {
			return nil, err
		}
		data = createChatData

	case USER_SEND_MESSAGE_REQUEST:
		var sendMessageData *ClientSendMessageRequest
		if err := json.Unmarshal(raw.Data, &sendMessageData); err != nil {
			return nil, err
		}
		data = sendMessageData

	case USER_GET_CHATS_REQUEST:
		var getChatsData *GetChatsRequest
		if err := json.Unmarshal(raw.Data, &getChatsData); err != nil {
			return nil, err
		}
		data = getChatsData

	case USER_GET_CHAT_HISTORY_REQUEST:
		var getChatHistoryData *GetChatHistoryRequest
		if err := json.Unmarshal(raw.Data, &getChatHistoryData); err != nil {
			return nil, err
		}
		data = getChatHistoryData

	default:
		log.Printf("unknown message type: %s", raw.Type)
		return &Request{
			Type: raw.Type,
			Data: nil,
		}, errors.New("unknown message type")
	}

	return &Request{
		Type: raw.Type,
		Data: data,
	}, nil
}

// functions
func (mp *MessageProcessor) ProcessMessage(senderId uuid.UUID, rawMessage []byte) {
	request, err := unmarshalRawMessage(rawMessage)
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

	switch request.Type {
	case USER_CREATE_CHAT_REQUEST:
		mp.handleCreateChatRequest(senderId, request)
	case USER_SEND_MESSAGE_REQUEST:
		mp.handleSendMessageRequest(senderId, request)
	case USER_GET_CHAT_HISTORY_REQUEST:
		mp.handleGetChatHistoryRequest(senderId, request)
	case USER_GET_CHATS_REQUEST:
		mp.handleGetChatsRequest(senderId, request)
	default:
		log.Printf("unknown message type: %s", request.Type)
		responseMessage := &Response{
			Type:  "",
			Data:  nil,
			Error: ErrUnknownRequestType.Error(),
		}
		mp.MessageSender.SendToUser(senderId, responseMessage)
		return
	}
}

func (mp *MessageProcessor) handleCreateChatRequest(senderId uuid.UUID, message *Request) {
	chatData, ok := message.Data.(*CreateChatRequest)
	if !ok {
		log.Printf("Invalid chat creation data")
		return
	}

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
	err = mp.ChatModel.AddUsersToChat(chat.Id, userIDs)
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

	log.Printf("Created chat: %s with ID: %s", chat.Name, chat.Id)

	// Create response message
	wsUserInfos := make([]UserInfo, 0, len(userInfos))
	for _, userInfo := range userInfos {
		wsUserInfos = append(wsUserInfos, UserInfo{
			ID:    userInfo.ID.String(),
			Email: userInfo.Email,
			Name:  userInfo.Name,
		})
	}
	responseData := &CreateChatResponse{
		Id:        chat.Id.String(),
		Name:      chat.Name,
		UpdatedAt: chat.UpdatedAt,
		UserInfos: wsUserInfos,
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

func (mp *MessageProcessor) handleSendMessageRequest(senderId uuid.UUID, request *Request) {
	messageData, ok := request.Data.(*ClientSendMessageRequest)
	if !ok {
		log.Printf("Invalid message data")
		return
	}
	chatId := uuid.MustParse(messageData.ChatID)
	content := messageData.Content

	// get chat by id
	_, err := mp.ChatModel.GetChat(chatId)
	if err != nil {
		log.Printf("Error getting chat: %v", err)
		responseMessage := &Response{
			Type:  USER_SEND_MESSAGE_RESPONSE,
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
			Type:  USER_SEND_MESSAGE_RESPONSE,
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
		Type: USER_SEND_MESSAGE_RESPONSE,
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

func (mp *MessageProcessor) handleGetChatHistoryRequest(senderId uuid.UUID, message *Request) {
	// messageData, ok := message.Data.(*ClientGetChatHistoryRequest)
	// if !ok {
	// 	log.Printf("Invalid message data")
	// 	return
	// }
}

func (mp *MessageProcessor) handleGetChatsRequest(senderId uuid.UUID, message *Request) {
	request, ok := message.Data.(*GetChatsRequest)
	if !ok {
		log.Printf("Invalid message data")
		responseMessage := &Response{
			Type:  USER_GET_CHATS_RESPONSE,
			Data:  nil,
			Error: ErrUnknownRequestType.Error(),
		}
		mp.MessageSender.SendToUser(senderId, responseMessage)
		return
	}
	modelChats, err := mp.ChatModel.GetChatsByUserID(senderId, request.CursorUpdatedAt, request.ChatID, request.Limit)
	if err != nil {
		log.Printf("Error getting chats: %v", err)
		return
	}
	chats := make([]ChatInfo, 0, len(modelChats))
	for _, modelChat := range modelChats {
		chats = append(chats, chatConvert(*modelChat))
	}

	responseMessage := &Response{
		Type:  USER_GET_CHATS_RESPONSE,
		Data:  GetChatsResponse{Chats: chats},
		Error: "",
	}
	mp.MessageSender.SendToUser(senderId, responseMessage)
}

func chatConvert(chat models.Chat) ChatInfo {
	userInfos := make([]UserInfo, 0, len(chat.UserInfos))
	for _, userInfo := range chat.UserInfos {
		userInfos = append(userInfos, UserInfo{
			ID:    userInfo.ID.String(),
			Email: userInfo.Email,
			Name:  userInfo.Name,
		})
	}
	return ChatInfo{
		Id:        chat.Id.String(),
		Name:      chat.Name,
		UpdatedAt: chat.UpdatedAt,
		UserInfos: userInfos,
	}
}
