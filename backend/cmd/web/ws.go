package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	upgrader websocket.Upgrader
}

func NewWebSocketHandler(clientOrigin string) *WebSocketHandler {
	return &WebSocketHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return r.Header.Get("Origin") == clientOrigin
			},
		},
	}
}

func (h *WebSocketHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// ✅ Extract token (e.g., from query param)
	token := r.URL.Query().Get("token")
	if !isValidToken(token) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// ✅ Upgrade to WebSocket
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// Handle incoming messages
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}
		fmt.Printf("Received: %s\n", msg)

		// Echo back
		conn.WriteMessage(websocket.TextMessage, msg)
	}
}

func isValidToken(token string) bool {
	// Decode and validate JWT here
	return token != "" // replace with real JWT verification
}
