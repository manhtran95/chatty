package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"chatty.mtran.io/internal/auth"
	messageprocessor "chatty.mtran.io/internal/message_processor"
	"github.com/google/uuid"
	gorilla "github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

// TestWebSocketConnection_Success tests successful WebSocket connection
func TestWebSocketConnection_Success(t *testing.T) {
	cleanDB(t, db)

	cleanup, server, conn, _ := setupTestAppFull(t)
	defer cleanup()
	defer server.Close()
	defer conn.Close()

	// Give some time for the connection to register with the hub
	time.Sleep(100 * time.Millisecond)

	t.Log("WebSocket connection established successfully")

	// Close connection before cleanup
	conn.Close()
	time.Sleep(50 * time.Millisecond)
}

// TestWebSocketConnection_InvalidOrigin tests connection with invalid origin
func TestWebSocketConnection_InvalidOrigin(t *testing.T) {
	cleanDB(t, db)

	app, cleanup := setupTestApp(t)
	defer cleanup()

	// Create WebSocket handler with specific origin
	handler := NewWebSocketHandler("http://localhost", app)

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testUserID := uuid.New().String()
		ctx := context.WithValue(r.Context(), auth.UserIDKey, testUserID)
		handler.Handle(w, r.WithContext(ctx))
	}))
	defer server.Close()

	// Convert http URL to ws URL
	wsURL := strings.Replace(server.URL, "http://", "ws://", 1)

	// Try to connect with wrong origin
	dialer := gorilla.Dialer{}
	conn, resp, err := dialer.Dial(wsURL, http.Header{
		"Origin": []string{"http://malicious-site.com"},
	})

	// Connection should fail due to origin check
	if err == nil {
		if conn != nil {
			conn.Close()
		}
		t.Fatal("Expected connection to fail with invalid origin")
	}

	if resp != nil && resp.StatusCode != http.StatusForbidden && resp.StatusCode != http.StatusUpgradeRequired {
		t.Logf("Expected 403 Forbidden or 426 Upgrade Required, got: %d", resp.StatusCode)
	}

	t.Log("Connection rejected as expected due to invalid origin")
}

func TestUnknownRequestType(t *testing.T) {
	cleanDB(t, db)

	cleanup, server, conn, testUserID := setupTestAppFull(t)
	defer cleanup()
	defer server.Close()
	defer conn.Close()

	// Give some time for the connection to register with the hub
	time.Sleep(100 * time.Millisecond)

	// Send a test message (use a random user ID since we don't track it from setupTestAppFull)
	testMessage := fmt.Sprintf(`{
		"type": "unknown_request_type",
		"senderId": "%s",
		"data": {
			"chatId": "%s",
			"content": "Hello, World!"
		}
	}`, testUserID, uuid.New().String())

	writeMessage(t, conn, testMessage)
	response := readMessage(t, conn)

	assert.Equal(t, messageprocessor.ErrUnknownRequestType.Error(), response.Error)
	assert.Equal(t, "", response.Type)
	assert.Equal(t, json.RawMessage("null"), response.Data)
}

func TestCreateChat(t *testing.T) {
	cleanDB(t, db)

	cleanup, server, _ := setupTestAppWithServerAndUsers(t)
	defer cleanup()
	defer server.Close()

	// Establish WebSocket connection as user 1
	accessToken := loginUserAndGetAccessToken(t, server, testUser1Email, testPassword)
	conn := getWebSocketConnectionWithAccessToken(t, server, accessToken)
	defer conn.Close()

	// Give some time for the connection to register with the hub
	time.Sleep(100 * time.Millisecond)

	chatID := createChatSuccess(t, conn)
	assert.NotEmpty(t, chatID)
}

// TODO
func TestCreateMultipleChatsAndGetChats(t *testing.T) {
	cleanDB(t, db)

	cleanup, server, _ := setupTestAppWithServerAndUsers(t)
	defer cleanup()
	defer server.Close()	
}

// TODO: Get Chats multiple times with different cursors

// TestWebSocketConnection_SendMessage tests sending a message through WebSocket
func TestSendMessageSuccess(t *testing.T) {
	cleanDB(t, db)

	cleanup, server, _ := setupTestAppWithServerAndUsers(t)
	defer cleanup()
	defer server.Close()

	// Establish WebSocket connection as user 1
	accessToken := loginUserAndGetAccessToken(t, server, testUser1Email, testPassword)
	conn := getWebSocketConnectionWithAccessToken(t, server, accessToken)
	defer conn.Close()

	// Give some time for the connection to register with the hub
	time.Sleep(100 * time.Millisecond)

	chatID := createChatSuccess(t, conn)

	message := "Hello, World!"

	testMessage := fmt.Sprintf(`{
		"type": "%s",
		"data": {
			"chatId": "%s",
			"content": "%s"
		}
	}`, messageprocessor.SEND_MESSAGE_REQUEST, chatID, message)

	writeMessage(t, conn, testMessage)
	response := readMessage(t, conn)

	assert.Equal(t, messageprocessor.SEND_MESSAGE_RESPONSE, response.Type)
	var responseData messageprocessor.SendMessageResponse
	err := json.Unmarshal(response.Data, &responseData)
	if err != nil {
		t.Fatalf("Failed to unmarshal response data: %v", err)
	}

	assert.Equal(t, chatID, responseData.ChatID)
	assert.Equal(t, message, responseData.Content)
}

// TestWebSocketConnection_MultipleCLients tests multiple concurrent connections
func TestWebSocketConnection_MultipleClients(t *testing.T) {
	cleanDB(t, db)

	app, cleanup := setupTestApp(t)
	defer cleanup()

	// Create WebSocket handler
	handler := NewWebSocketHandler("http://localhost", app)

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testUserID := uuid.New().String()
		ctx := context.WithValue(r.Context(), auth.UserIDKey, testUserID)
		handler.Handle(w, r.WithContext(ctx))
	}))
	defer server.Close()

	// Convert http URL to ws URL
	wsURL := strings.Replace(server.URL, "http://", "ws://", 1)

	// Create multiple connections
	numClients := 5
	connections := make([]*gorilla.Conn, numClients)

	for i := 0; i < numClients; i++ {
		dialer := gorilla.Dialer{}
		conn, _, err := dialer.Dial(wsURL, http.Header{
			"Origin": []string{"http://localhost"},
		})
		if err != nil {
			t.Fatalf("Failed to connect client %d: %v", i, err)
		}
		connections[i] = conn
		defer conn.Close()
	}

	// Give some time for all connections to register
	time.Sleep(200 * time.Millisecond)

	// Verify all connections are established
	for i, conn := range connections {
		if conn == nil {
			t.Errorf("Client %d connection is nil", i)
		}
	}

	t.Logf("Successfully established %d concurrent WebSocket connections", numClients)

	// Close all connections before cleanup
	for _, conn := range connections {
		if conn != nil {
			conn.Close()
		}
	}
	time.Sleep(100 * time.Millisecond)
}
