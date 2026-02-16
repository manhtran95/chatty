package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	messageprocessor "chatty.mtran.io/internal/message_processor"
	gorilla "github.com/gorilla/websocket"
)

func getWebSocketConnectionWithAccessToken(t *testing.T, server *httptest.Server, accessToken string) *gorilla.Conn {
	t.Helper()

	wsURL := strings.Replace(server.URL, "http://", "ws://", 1) + "/ws?access_token=" + accessToken
	dialer := gorilla.Dialer{}
	conn, resp, err := dialer.Dial(wsURL, http.Header{
		"Origin": []string{"http://localhost:5173"},
	})
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v (response: %v)", err, resp)
	}
	return conn
}

func writeMessage(t *testing.T, conn *gorilla.Conn, message string) {
	err := conn.WriteMessage(gorilla.TextMessage, []byte(message))
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}
}

func readMessage(t *testing.T, conn *gorilla.Conn) messageprocessor.Response {
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, _response, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}
	response := messageprocessor.Response{}
	err = json.Unmarshal(_response, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	return response
}
