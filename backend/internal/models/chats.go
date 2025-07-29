package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Chat represents a chat room
type Chat struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// ChatModel wraps a database connection pool for chat operations
type ChatModel struct {
	DB *sql.DB
}

// InsertChat creates a new chat with the given name
func (m *ChatModel) InsertChat(name string) (*Chat, error) {
	chat := &Chat{
		Name: name,
	}

	stmt := `INSERT INTO chats (name) VALUES ($1) RETURNING id, name, created_at`

	err := m.DB.QueryRow(stmt, name).Scan(&chat.ID, &chat.Name, &chat.CreatedAt)
	if err != nil {
		return nil, err
	}

	return chat, nil
}

// GetChat retrieves a chat by ID
func (m *ChatModel) GetChat(id uuid.UUID) (*Chat, error) {
	chat := &Chat{}

	stmt := `SELECT id, name, created_at FROM chats WHERE id = $1`

	err := m.DB.QueryRow(stmt, id).Scan(&chat.ID, &chat.Name, &chat.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return chat, nil
}

// GetAllChatsByUserID retrieves all chats that a user is a member of
func (m *ChatModel) GetAllChatsByUserID(userID uuid.UUID) ([]*Chat, error) {
	stmt := `SELECT c.id, c.name, c.created_at 
	         FROM chats c 
	         INNER JOIN chat_users cu ON c.id = cu.chat_id 
	         WHERE cu.user_id = $1 
	         ORDER BY c.created_at DESC`

	rows, err := m.DB.Query(stmt, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []*Chat

	for rows.Next() {
		chat := &Chat{}
		err := rows.Scan(&chat.ID, &chat.Name, &chat.CreatedAt)
		if err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return chats, nil
}

// GetChatMembers retrieves all members of a specific chat
func (m *ChatModel) GetChatMembers(chatID uuid.UUID) ([]*UserInfo, error) {
	stmt := `SELECT u.id, u.name, u.email 
	         FROM users u 
	         INNER JOIN chat_users cu ON u.id = cu.user_id 
	         WHERE cu.chat_id = $1 
	         ORDER BY u.name`

	rows, err := m.DB.Query(stmt, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*UserInfo

	for rows.Next() {
		member := &UserInfo{}
		err := rows.Scan(&member.ID, &member.Name, &member.Email)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return members, nil
}

// AddUserToChat adds a user to a chat
func (m *ChatModel) AddUserToChat(chatID, userID uuid.UUID) error {
	stmt := `INSERT INTO chat_users (chat_id, user_id) VALUES ($1, $2) ON CONFLICT (chat_id, user_id) DO NOTHING`

	_, err := m.DB.Exec(stmt, chatID, userID)
	return err
}

// AddUsersToChat adds multiple users to a chat
func (m *ChatModel) AddUsersToChat(chatID uuid.UUID, userIDs []uuid.UUID) error {
	if len(userIDs) == 0 {
		return nil
	}

	// Build a single query with all users using VALUES clause
	stmt := `INSERT INTO chat_users (chat_id, user_id) VALUES `
	args := make([]any, 0, len(userIDs)+1)
	args = append(args, chatID)
	placeholders := make([]string, len(userIDs))
	for i := range userIDs {
		placeholders[i] = fmt.Sprintf("($1, $%d)", i+2)
		args = append(args, userIDs[i])
	}

	stmt += strings.Join(placeholders, ", ")
	stmt += ` ON CONFLICT (chat_id, user_id) DO NOTHING`

	_, err := m.DB.Exec(stmt, args...)
	return err
}
