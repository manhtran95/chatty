package main

import (
	"net/http"
)

// The routes() method returns a servemux containing our application routes.
func (app *application) routes() http.Handler {

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/user/signup", app.userSignup)
	mux.HandleFunc("/user/login", app.userLogin)
	mux.HandleFunc("/user/refresh", app.refreshToken)
	return app.logRequest(withCORS(mux))
}
