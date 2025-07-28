package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Message represents a chat message
type Message struct {
	ID        uuid.UUID `json:"id"`
	SenderID  uuid.UUID `json:"sender_id"`
	ChatID    uuid.UUID `json:"chat_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// MessageModel wraps a database connection pool for message operations
type MessageModel struct {
	DB *sql.DB
}

// InsertMessage creates a new message
func (m *MessageModel) InsertMessage(senderID, chatID uuid.UUID, content string) (*Message, error) {
	message := &Message{
		SenderID: senderID,
		ChatID:   chatID,
		Content:  content,
	}

	stmt := `INSERT INTO messages (sender_id, chat_id, content) VALUES ($1, $2, $3) RETURNING id, sender_id, chat_id, content, created_at`

	err := m.DB.QueryRow(stmt, senderID, chatID, content).Scan(
		&message.ID, &message.SenderID, &message.ChatID, &message.Content, &message.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return message, nil
}

// GetMessagesByChat retrieves messages for a specific chat with pagination
func (m *MessageModel) GetMessagesByChat(chatID uuid.UUID, limit, offset int) ([]*Message, error) {
	stmt := `SELECT id, sender_id, chat_id, content, created_at FROM messages WHERE chat_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	rows, err := m.DB.Query(stmt, chatID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*Message

	for rows.Next() {
		message := &Message{}
		err := rows.Scan(&message.ID, &message.SenderID, &message.ChatID, &message.Content, &message.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
