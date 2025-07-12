package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Snippetbox"))
}

func main() {
	if err := godotenv.Load(".env.backend"); err != nil {
		log.Println("No .env file found")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", home)

	log.Printf("Server starting on port %s", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
	log.Fatal(err)
}
