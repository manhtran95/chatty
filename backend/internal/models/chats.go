package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Chat represents a chat room
type Chat struct {
	Id        uuid.UUID   `json:"id"`
	Name      string      `json:"name"`
	UpdatedAt time.Time   `json:"updated_at"`
	UserInfos []*UserInfo `json:"user_infos"`
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

	stmt := `INSERT INTO chats (name) VALUES ($1) RETURNING id, name, updated_at`

	err := m.DB.QueryRow(stmt, name).Scan(&chat.Id, &chat.Name, &chat.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return chat, nil
}

// GetChat retrieves a chat by ID
func (m *ChatModel) GetChat(id uuid.UUID) (*Chat, error) {
	chat := &Chat{}

	stmt := `SELECT id, name, updated_at FROM chats WHERE id = $1`

	err := m.DB.QueryRow(stmt, id).Scan(&chat.Id, &chat.Name, &chat.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return chat, nil
}

// GetChatsByUserID retrieves all chats that a user is a member of
func (m *ChatModel) GetChatsByUserID(userID uuid.UUID, cursorUpdatedAt *time.Time, chatId *string, limit int) ([]*Chat, error) {
	var rows *sql.Rows
	var err error

	if cursorUpdatedAt == nil {
		stmt := `SELECT c.id, c.name, c.updated_at 
		FROM chats c 
		INNER JOIN chat_users cu ON c.id = cu.chat_id 
		WHERE cu.user_id = $1 
		ORDER BY c.updated_at DESC
		LIMIT $2`
		rows, err = m.DB.Query(stmt, userID, limit)
	} else {
		chatIdUUID := uuid.MustParse(*chatId)
		stmt := `SELECT c.id, c.name, c.updated_at 
		FROM chats c 
		INNER JOIN chat_users cu ON c.id = cu.chat_id 
		WHERE cu.user_id = $1 
		AND (c.updated_at, c.id) < ($2, $3)
		ORDER BY c.updated_at DESC
		LIMIT $4`
		rows, err = m.DB.Query(stmt, userID, *cursorUpdatedAt, chatIdUUID, limit)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []*Chat
	var chatIds []uuid.UUID

	for rows.Next() {
		chat := &Chat{}
		err := rows.Scan(&chat.Id, &chat.Name, &chat.UpdatedAt)
		if err != nil {
			return nil, err
		}
		chats = append(chats, chat)
		chatIds = append(chatIds, chat.Id)
	}

	// get user infos
	chatToUserMap, err := m.getChatUserInfos(chatIds)
	if err != nil {
		return nil, err
	}

	for _, chat := range chats {
		chat.UserInfos = chatToUserMap[chat.Id]
	}

	return chats, nil
}

func (m *ChatModel) getChatUserInfos(chatIds []uuid.UUID) (map[uuid.UUID][]*UserInfo, error) {
	chatToUserMap := make(map[uuid.UUID][]*UserInfo)
	stmt := `SELECT users.id, users.name, users.email, chat_users.chat_id
	         FROM users
	         INNER JOIN chat_users ON users.id = chat_users.user_id 
	         WHERE chat_users.chat_id = ANY($1)
	        `
	rows, err := m.DB.Query(stmt, pq.Array(chatIds))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		userInfo := &UserInfo{}
		chatId := uuid.UUID{}
		err = rows.Scan(&userInfo.ID, &userInfo.Name, &userInfo.Email, &chatId)
		if err != nil {
			return nil, err
		}
		if _, ok := chatToUserMap[chatId]; !ok {
			chatToUserMap[chatId] = make([]*UserInfo, 0)
		}
		chatToUserMap[chatId] = append(chatToUserMap[chatId], userInfo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return chatToUserMap, nil
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
