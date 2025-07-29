package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 24 * time.Hour // 24 hours instead of 1 minute

	// Send pings to peer with this period. Must be less than pongWait
	pingPeriod = 30 * time.Second // Send ping every 30 seconds

	// Maximum message size allowed from peer
	maxMessageSize = 512
)

// Client represents a connected WebSocket client
type Client struct {
	UserID uuid.UUID
	Conn   *websocket.Conn
	Send   chan []byte
	Hub    *Hub
	mu     sync.Mutex
}

// NewClient creates a new WebSocket client
func NewClient(conn *websocket.Conn, userID uuid.UUID, hub *Hub) *Client {
	return &Client{
		UserID: userID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		Hub:    hub,
	}
}

type rawMessage struct {
	Type     string          `json:"type"`
	Data     json.RawMessage `json:"data"`
	SenderID string          `json:"senderId"`
}

func unmarshalRawMessage(p []byte) (message *Message, err error) {
	var raw rawMessage
	if err := json.Unmarshal(p, &raw); err != nil {
		return nil, err
	}

	var data MessageData
	switch raw.Type {
	case CLIENT_CREATE_CHAT:
		var createChatData *ClientCreateChatData
		if err := json.Unmarshal(raw.Data, &createChatData); err != nil {
			return nil, err
		}
		data = createChatData

	case CLIENT_SEND_MESSAGE:
		var sendMessageData *ClientSendMessageData
		if err := json.Unmarshal(raw.Data, &sendMessageData); err != nil {
			return nil, err
		}
		data = sendMessageData

	default:
		log.Printf("unknown message type: %s", raw.Type)
		return nil, nil
	}

	return &Message{
		Type:     raw.Type,
		Data:     data,
		SenderID: raw.SenderID,
	}, nil
}

// readPump pumps messages from the WebSocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		log.Printf("readPump: closing connection")
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			log.Printf("readPump: error: %v", err)
			break
		}

		message, err := unmarshalRawMessage(p)
		log.Printf("readPump: %v", message)
		if err != nil {
			log.Printf("unmarshalRawMessage: error: %v", err)
			break
		}

		// Send message to hub for processing
		c.Hub.Broadcast <- message
	}
}

// writePump pumps messages from the hub to the WebSocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// SendMessage sends a message to this client
func (c *Client) SendMessage(message *Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	select {
	case c.Send <- data:
		return nil
	default:
		return ErrClientBufferFull
	}
}

// Start starts the client's read and write pumps
func (c *Client) Start() {
	// Register client with hub
	c.Hub.Register <- c

	// Start read and write pumps
	go c.readPump()
	go c.writePump()
}

// Custom errors
var (
	ErrClientBufferFull = &websocket.CloseError{
		Code: websocket.CloseInternalServerErr,
		Text: "client send buffer is full",
	}
)
