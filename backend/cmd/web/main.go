package main

import (
	"database/sql"

	_ "github.com/lib/pq"

	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"chatty.mtran.io/internal/models"
	"chatty.mtran.io/internal/websocket"
	"github.com/go-playground/form/v4"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"

	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/joho/godotenv"
)

type application struct {
	errorLog    *log.Logger
	infoLog     *log.Logger
	users       *models.UserModel
	chats       *models.ChatModel
	messages    *models.MessageModel
	formDecoder *form.Decoder
	hub         *websocket.Hub
}

func main() {
	err := godotenv.Load(".env.backend")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "Chatty INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "Chatty ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// db
	db, err := openDB()
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	migrateDB(errorLog)

	formDecoder := form.NewDecoder()

	// Initialize models
	userModel := &models.UserModel{DB: db}
	chatModel := &models.ChatModel{DB: db}
	messageModel := &models.MessageModel{DB: db}

	// Initialize WebSocket hub with all models
	hub := websocket.NewHub(chatModel, userModel, messageModel)
	go hub.Run() // Start the hub in a goroutine

	app := &application{
		errorLog:    errorLog,
		infoLog:     infoLog,
		users:       userModel,
		chats:       chatModel,
		messages:    messageModel,
		formDecoder: formDecoder,
		hub:         hub,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	infoLog.Printf("Starting server on :%s", *addr)

	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func migrateDB(errorLog *log.Logger) {
	m, err := migrate.New("file://migrations", fmt.Sprintf("%s?sslmode=disable", os.Getenv("DATABASE_URL")))
	if err != nil {
		errorLog.Fatal(err)
	}
	// Migrate all the way up ...
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		errorLog.Fatal(err)
	}
}

func openDB() (*sql.DB, error) {
	dsn := fmt.Sprintf("%s?sslmode=disable", os.Getenv("DATABASE_URL"))

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
