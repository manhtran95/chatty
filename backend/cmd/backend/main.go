package main

import (
	"log"
	"net/http"
	"os"

	"chatty.mtran.io/internal/api"
)

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	mux := http.NewServeMux()
	mux.HandleFunc("/", api.Home)
	mux.HandleFunc("/snippet/view", api.SnippetView)
	mux.HandleFunc("/snippet/create", api.SnippetCreate)
	infoLog.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	errorLog.Fatal(err)
}
