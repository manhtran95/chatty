package main

// TODO: JWT validation is not implemented yet

import (
	"fmt"
	"net/http"

	"chatty.mtran.io/internal/auth"
	ws "chatty.mtran.io/internal/websocket"
	gorilla "github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	upgrader gorilla.Upgrader
	app      *application
}

func NewWebSocketHandler(clientOrigin string, app *application) *WebSocketHandler {
	return &WebSocketHandler{
		upgrader: gorilla.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return r.Header.Get("Origin") == clientOrigin
			},
		},
		app: app,
	}
}

func (h *WebSocketHandler) Handle(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)

	// Upgrade to WebSocket
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}

	// Create new client
	client := ws.NewClient(conn, userID, h.app.hub)

	// Start the client
	client.Start()
}
