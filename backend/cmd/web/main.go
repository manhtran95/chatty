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
	formDecoder *form.Decoder
	wsHandler   *WebSocketHandler
}

func main() {
	err := godotenv.Load(".env.backend")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// db
	db, err := openDB()
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	migrateDB(errorLog)

	formDecoder := form.NewDecoder()

	// Initialize WebSocket handler with client origin from env
	clientOrigin := os.Getenv("CLIENT_ORIGIN")
	if clientOrigin == "" {
		clientOrigin = "http://localhost:5173" // default fallback
	}
	wsHandler := NewWebSocketHandler(clientOrigin)

	app := &application{
		errorLog:    errorLog,
		infoLog:     infoLog,
		users:       &models.UserModel{DB: db},
		formDecoder: formDecoder,
		wsHandler:   wsHandler,
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
