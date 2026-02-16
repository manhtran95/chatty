package messageprocessor

import (
	"encoding/json"
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

// functions
func (mp *MessageProcessor) ProcessMessage(senderId uuid.UUID, message []byte) {
	// request, err := unmarshalRawMessage(rawMessage)
	var request Request
	err := json.Unmarshal(message, &request)
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
	case CREATE_CHAT_REQUEST:
		var reqData CreateChatRequest
		if err := json.Unmarshal(request.Data, &reqData); err != nil {
			log.Printf("Error unmarshalling create chat data: %v", err)
			return
		}
		mp.handleCreateChatRequest(senderId, reqData)
	case SEND_MESSAGE_REQUEST:
		var reqData SendMessageRequest
		if err := json.Unmarshal(request.Data, &reqData); err != nil {
			log.Printf("Error unmarshalling send message data: %v", err)
			return
		}
		mp.handleSendMessageRequest(senderId, reqData)
	case GET_CHAT_HISTORY_REQUEST:
		var reqData GetChatHistoryRequest
		if err := json.Unmarshal(request.Data, &reqData); err != nil {
			log.Printf("Error unmarshalling get chat history data: %v", err)
			return
		}
		mp.handleGetChatHistoryRequest(senderId, reqData)
	case GET_CHATS_REQUEST:
		var reqData GetChatsRequest
		if err := json.Unmarshal(request.Data, &reqData); err != nil {
			log.Printf("Error unmarshalling get chats data: %v", err)
			return
		}
		mp.handleGetChatsRequest(senderId, reqData)
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

func (mp *MessageProcessor) handleCreateChatRequest(senderId uuid.UUID, reqData CreateChatRequest) {
	// check if participantIDs are valid
	if len(reqData.ParticipantEmails) < 2 {
		log.Printf("Cannot create a chat with less than 2 participants")
		responseMessage := &Response{
			Type:  CREATE_CHAT_RESPONSE,
			Data:  nil,
			Error: ErrCannotCreateChatWithLessThan2Participants.Error(),
		}
		mp.MessageSender.SendToUser(senderId, responseMessage)
		return
	}

	userInfos, err := mp.UserModel.UserInfosByEmails(reqData.ParticipantEmails)
	if err != nil {
		log.Printf("Error checking if user emails exist: %v", err)
		responseMessage := &Response{
			Type:  CREATE_CHAT_RESPONSE,
			Data:  nil,
			Error: ErrUserDoesNotExist.Error(),
		}
		mp.MessageSender.SendToUser(senderId, responseMessage)
		return
	}

	// TO DO: check if set of users already in a chat

	// Create chat in database
	chat, err := mp.ChatModel.InsertChat(reqData.Name)
	if err != nil {
		log.Printf("Error creating chat: %v", err)
		responseMessage := &Response{
			Type:  CREATE_CHAT_RESPONSE,
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
			Type:  CREATE_CHAT_RESPONSE,
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
	responseData := CreateChatResponse{
		Id:        chat.Id.String(),
		Name:      chat.Name,
		UpdatedAt: chat.UpdatedAt,
		UserInfos: wsUserInfos,
	}

	responseMessage := &Response{
		Type:  CREATE_CHAT_RESPONSE,
		Data:  getJsonRawMessage(responseData),
		Error: "",
	}

	for _, userID := range userIDs {
		mp.MessageSender.SendToUser(userID, responseMessage)
	}
}

func getJsonRawMessage(data any) json.RawMessage {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling data: %v", err)
		return nil
	}
	return json.RawMessage(jsonData)
}

func (mp *MessageProcessor) handleSendMessageRequest(senderId uuid.UUID, reqData SendMessageRequest) {
	chatId := uuid.MustParse(reqData.ChatID)
	content := reqData.Content

	// get chat by id
	_, err := mp.ChatModel.GetChat(chatId)
	if err != nil {
		log.Printf("Error getting chat: %v", err)
		responseMessage := &Response{
			Type:  SEND_MESSAGE_RESPONSE,
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
			Type:  SEND_MESSAGE_RESPONSE,
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
	responseData := SendMessageResponse{
		ChatID:    chatId.String(),
		SenderID:  message.SenderID.String(),
		Content:   message.Content,
		Timestamp: message.CreatedAt.String(),
		MessageID: message.ID.String(),
	}
	responseMessage := &Response{
		Type:  SEND_MESSAGE_RESPONSE,
		Data:  getJsonRawMessage(responseData),
		Error: "",
	}

	for _, memberId := range memberIds {
		mp.MessageSender.SendToUser(memberId, responseMessage)
	}
}

func (mp *MessageProcessor) handleGetChatHistoryRequest(senderId uuid.UUID, reqData GetChatHistoryRequest) {
	// messageData, ok := message.Data.(*ClientGetChatHistoryRequest)
	// if !ok {
	// 	log.Printf("Invalid message data")
	// 	return
	// }
}

func (mp *MessageProcessor) handleGetChatsRequest(senderId uuid.UUID, reqData GetChatsRequest) {
	modelChats, err := mp.ChatModel.GetChatsByUserID(senderId, reqData.CursorUpdatedAt, reqData.ChatID, reqData.Limit)
	if err != nil {
		log.Printf("Error getting chats: %v", err)
		return
	}
	chats := make([]ChatInfo, 0, len(modelChats))
	for _, modelChat := range modelChats {
		chats = append(chats, chatConvert(*modelChat))
	}

	responseData := GetChatsResponse{
		Chats: chats,
	}
	responseMessage := &Response{
		Type:  GET_CHATS_RESPONSE,
		Data:  getJsonRawMessage(responseData),
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
