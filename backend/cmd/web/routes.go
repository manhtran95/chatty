package main

import (
	"net/http"
	"os"
)

// The routes() method returns a servemux containing our application routes.
func (app *application) routes() http.Handler {
	clientOrigin := os.Getenv("CLIENT_ORIGIN")

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/user/signup", app.userSignup)
	mux.HandleFunc("/user/login", app.userLogin)
	mux.HandleFunc("/user/refresh", app.refreshToken)
	mux.HandleFunc("/user/logout", app.userLogout)
	mux.HandleFunc("/ws", app.wsHandler.Handle)

	return app.logRequest(app.logResponse(withCORS(clientOrigin)(mux)))
}
