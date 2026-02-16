package main

import (
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
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	gorilla "github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	logrus "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var db *sql.DB

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		logrus.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		logrus.Fatalf("Could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "16-alpine",
		Env: []string{
			"POSTGRES_PASSWORD=chatty_password",
			"POSTGRES_USER=chatty",
			"POSTGRES_DB=chatty_db",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		logrus.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://chatty:chatty_password@%s/chatty_db?sslmode=disable", hostAndPort)

	logrus.Println("Connecting to database on url: ", databaseUrl)

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		db, err = sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	log.Println("Connected to database")

	// Run migrations
	os.Setenv("DATABASE_URL", databaseUrl)
	if err := runMigrations(databaseUrl); err != nil {
		log.Fatalf("Could not run migrations: %s", err)
	}
	log.Println("Migrations completed")

	defer func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}()

	// run tests
	m.Run()
}

// runMigrations runs database migrations
func runMigrations(databaseURL string) error {
	// The databaseURL already has sslmode=disable
	m, err := migrate.New("file://../../migrations", databaseURL)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func cleanDB(t *testing.T, db *sql.DB) {
	tables := []string{"users", "chats", "messages"}
	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			t.Fatalf("Failed to clean table %s: %v", table, err)
		}
	}
}

// setupTestApp creates a test application with mocked dependencies
// Returns the app and a cleanup function that should be deferred
func setupTestApp(t *testing.T) (*application, func()) {
	t.Helper()

	// Set up test environment variables (only once)
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing")
	os.Setenv("JWT_ACCESS_TOKEN_EXPIRE_MINUTES", "15")
	os.Setenv("JWT_REFRESH_TOKEN_EXPIRE_DAYS", "7")
	os.Setenv("PASSWORD_MIN_LENGTH", "8")

	// Create mock database connection (you may want to use a real test database)
	// For now, we'll skip database-dependent tests

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

const (
	testUser1Email = "test1@example.com"
	testUser1Name  = "Test User 1"
	testUser2Email = "test2@example.com"
	testUser2Name  = "Test User 2"
	testPassword   = "password123"
)

// setupTestAppWithServerAndUsers creates a test application with HTTP server (no WebSocket connection)
// Returns cleanup function, HTTP server, and the app
func setupTestAppWithServerAndUsers(t *testing.T) (func(), *httptest.Server, *application) {
	t.Helper()

	app, cleanup := setupTestApp(t)

	// Create HTTP test server with the app's routes
	server := httptest.NewServer(app.routes())

	signupUserSuccess(t, server, testUser1Name, testUser1Email, testPassword)
	signupUserSuccess(t, server, testUser2Name, testUser2Email, testPassword)

	return cleanup, server, app
}

func setupTestAppWithServer(t *testing.T) (func(), *httptest.Server, *application) {
	t.Helper()

	app, cleanup := setupTestApp(t)

	// Create HTTP test server with the app's routes
	server := httptest.NewServer(app.routes())

	return cleanup, server, app
}

// createChatSuccess creates a chat and returns the chat ID
func createChatSuccess(t *testing.T, conn *gorilla.Conn) string {
	t.Helper()

	testChatName := "Test Chat"
	createChatRequest := fmt.Sprintf(`{
		"type": "%s",
		"data": {
			"name": "%s",
			"participantEmails": ["%s", "%s"]
		}
	}`, messageprocessor.CREATE_CHAT_REQUEST, testChatName, testUser1Email, testUser2Email)
	writeMessage(t, conn, createChatRequest)

	response := readMessage(t, conn)

	var responseData messageprocessor.CreateChatResponse
	err := json.Unmarshal(response.Data, &responseData)
	if err != nil {
		t.Fatalf("Failed to unmarshal response data: %v", err)
	}

	assert.Equal(t, messageprocessor.CREATE_CHAT_RESPONSE, response.Type)
	assert.Equal(t, testChatName, responseData.Name)
	assert.Equal(t, len(responseData.UserInfos), 2)
	assert.Equal(t, responseData.UserInfos[0].Email, testUser1Email)
	assert.Equal(t, responseData.UserInfos[0].Name, testUser1Name)
	assert.Equal(t, responseData.UserInfos[1].Email, testUser2Email)
	assert.Equal(t, responseData.UserInfos[1].Name, testUser2Name)
	return responseData.Id
}

func setupTestAppFull(t *testing.T) (func(), *httptest.Server, *gorilla.Conn, string) {
	t.Helper()

	// Use setupTestApp to get the app with all setup
	app, cleanup := setupTestApp(t)

	// Create HTTP test server with the app's full routes (including /user/signup, /ws, etc.)
	server := httptest.NewServer(app.routes())

	// Generate a test user ID and access token for WebSocket authentication
	testUserID := uuid.New().String()
	accessToken := auth.GenerateAccessToken(testUserID)

	// Convert http URL to ws URL for WebSocket connection
	wsURL := strings.Replace(server.URL, "http://", "ws://", 1) + "/ws?access_token=" + accessToken

	// Create WebSocket client connection
	dialer := gorilla.Dialer{}
	conn, resp, err := dialer.Dial(wsURL, http.Header{
		"Origin": []string{"http://localhost:5173"}, // Use the configured client origin
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

func TestRealbob(t *testing.T) {
	// all tests
}
