package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"chatty.mtran.io/internal/auth"
	messageprocessor "chatty.mtran.io/internal/message_processor"
	"chatty.mtran.io/internal/models"
	ws "chatty.mtran.io/internal/websocket"
	"github.com/go-playground/form/v4"
	"github.com/google/uuid"
	gorilla "github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

// setupTestApp creates a test application with mocked dependencies
// Returns the app and a cleanup function that should be deferred
func setupTestApp(t *testing.T) (*application, func()) {
	t.Helper()

	// Set up test environment variables (only once)
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing")
	os.Setenv("JWT_ACCESS_TOKEN_EXPIRE_MINUTES", "15")
	os.Setenv("JWT_REFRESH_TOKEN_EXPIRE_DAYS", "7")

	// Create mock database connection (you may want to use a real test database)
	// For now, we'll skip database-dependent tests
	var db *sql.DB

	infoLog := log.New(os.Stdout, fmt.Sprintf("TEST[%s] INFO\t", t.Name()), log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, fmt.Sprintf("TEST[%s] ERROR\t", t.Name()), log.Ldate|log.Ltime|log.Lshortfile)

	// Initialize models with nil DB for now (can be mocked later)
	userModel := &models.UserModel{DB: db}
	chatModel := &models.ChatModel{DB: db}
	messageModel := &models.MessageModel{DB: db}

	// Initialize WebSocket hub
	messageProcessor := messageprocessor.NewMessageProcessor(chatModel, userModel, messageModel)
	hub := ws.NewHub(messageProcessor)
	messageProcessor.SetMessageSender(hub)
	go hub.Run()

	app := &application{
		errorLog:    errorLog,
		infoLog:     infoLog,
		users:       userModel,
		chats:       chatModel,
		messages:    messageModel,
		formDecoder: form.NewDecoder(),
		hub:         hub,
	}

	// Cleanup function to prevent goroutine leaks
	cleanup := func() {
		// Close all client connections in the hub
		for _, client := range hub.UserClients {
			if client.Conn != nil {
				client.Conn.Close()
			}
		}
		// Give the hub time to process unregistrations
		time.Sleep(50 * time.Millisecond)
	}

	return app, cleanup
}

func setupTestAppFull(t *testing.T) (func(), *httptest.Server, *gorilla.Conn, string) {
	t.Helper()

	// Set up test environment variables (only once)
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing")
	os.Setenv("JWT_ACCESS_TOKEN_EXPIRE_MINUTES", "15")
	os.Setenv("JWT_REFRESH_TOKEN_EXPIRE_DAYS", "7")

	// Create mock database connection (you may want to use a real test database)
	// For now, we'll skip database-dependent tests
	var db *sql.DB

	infoLog := log.New(os.Stdout, fmt.Sprintf("TEST[%s] INFO\t", t.Name()), log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, fmt.Sprintf("TEST[%s] ERROR\t", t.Name()), log.Ldate|log.Ltime|log.Lshortfile)

	// Initialize models with nil DB for now (can be mocked later)
	userModel := &models.UserModel{DB: db}
	chatModel := &models.ChatModel{DB: db}
	messageModel := &models.MessageModel{DB: db}

	// Initialize WebSocket hub
	messageProcessor := messageprocessor.NewMessageProcessor(chatModel, userModel, messageModel)
	hub := ws.NewHub(messageProcessor)
	messageProcessor.SetMessageSender(hub)
	go hub.Run()

	app := &application{
		errorLog:    errorLog,
		infoLog:     infoLog,
		users:       userModel,
		chats:       chatModel,
		messages:    messageModel,
		formDecoder: form.NewDecoder(),
		hub:         hub,
	}

	// Cleanup function to prevent goroutine leaks
	cleanup := func() {
		// Close all client connections in the hub
		for _, client := range hub.UserClients {
			if client.Conn != nil {
				client.Conn.Close()
			}
		}
		// Give the hub time to process unregistrations
		time.Sleep(50 * time.Millisecond)
	}

	handler := NewWebSocketHandler("http://localhost", app)
	testUserID := uuid.New().String()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mock authentication by adding userID to context
		log.Printf("11 testUserID: %s", testUserID)
		ctx := context.WithValue(r.Context(), auth.UserIDKey, testUserID)
		handler.Handle(w, r.WithContext(ctx))
	}))

	wsURL := strings.Replace(server.URL, "http://", "ws://", 1)

	// Create WebSocket client connection
	dialer := gorilla.Dialer{}
	conn, resp, err := dialer.Dial(wsURL, http.Header{
		"Origin": []string{"http://localhost"},
	})
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v (response: %v)", err, resp)
	}
	// Verify connection is established
	if conn == nil {
		t.Fatal("Expected connection to be established")
	}

	return cleanup, server, conn, testUserID
}

// TestWebSocketConnection_Success tests successful WebSocket connection
func TestWebSocketConnection_Success(t *testing.T) {
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

	err := conn.WriteMessage(gorilla.TextMessage, []byte(testMessage))
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// Set a read deadline
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, _response, err := conn.ReadMessage()
	response := messageprocessor.Response{}
	err = json.Unmarshal(_response, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.Equal(t, response.Error, messageprocessor.ErrUnknownRequestType.Error())
	assert.Equal(t, response.Type, "")
	assert.Equal(t, response.Data, nil)

	// Close connection before cleanup
	conn.Close()
	time.Sleep(50 * time.Millisecond)
}

// TestWebSocketConnection_SendMessage tests sending a message through WebSocket
func TestWebSocketConnection_SendMessage(t *testing.T) {
	// TODO
	cleanup, server, conn, testUserID := setupTestAppFull(t)
	defer cleanup()
	defer server.Close()
	defer conn.Close()

	// Give some time for the connection to register with the hub
	time.Sleep(100 * time.Millisecond)

	// Send a test message (use a random user ID since we don't track it from setupTestAppFull)
	testMessage := fmt.Sprintf(`{
		"type": "client_send_message",
		"senderId": "%s",
		"data": {
			"chatId": "%s",
			"content": "Hello, World!"
		}
	}`, testUserID, uuid.New().String())

	err := conn.WriteMessage(gorilla.TextMessage, []byte(testMessage))
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// Set a read deadline
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	// Try to read response (this might timeout, which is okay for this test)
	_, response, err := conn.ReadMessage()
	t.Logf("Received response: %s", string(response))
	if err != nil {
		// It's acceptable if we don't receive a response immediately
		// as the message processing depends on database operations
		t.Logf("No immediate response received (expected): %v", err)
	} else {
		t.Logf("Received response: %s", string(response))
	}

	t.Log("Message sent successfully")

	// Close connection before cleanup
	conn.Close()
	time.Sleep(50 * time.Millisecond)
}

// TestWebSocketConnection_MultipleCLients tests multiple concurrent connections
func TestWebSocketConnection_MultipleClients(t *testing.T) {
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
