package messageprocessor

import (
	"errors"
)

var (
	ErrUnknownRequestType = errors.New("messageprocessor: unknown request type")
	ErrCannotCreateChatWithLessThan2Participants = errors.New("messageprocessor: cannot create a chat with less than 2 participants")
	ErrUserDoesNotExist = errors.New("messageprocessor: user does not exist")
	ErrCannotCreateChat = errors.New("messageprocessor: cannot create chat")
	ErrCannotAddParticipantsToChat = errors.New("messageprocessor: cannot add participants to chat")
	ErrCannotGetChat = errors.New("messageprocessor: cannot get chat")
	ErrCannotSaveMessage = errors.New("messageprocessor: cannot save message")
)
