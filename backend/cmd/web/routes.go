package main

import (
	"net/http"
	"os"

	"chatty.mtran.io/internal/auth"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// The routes() method returns a servemux containing our application routes.
func (app *application) routes() http.Handler {
	// Get client origin from environment or use default
	clientOrigin := os.Getenv("CLIENT_ORIGIN")
	if clientOrigin == "" {
		clientOrigin = "http://localhost:5173" // default fallback
	}
	wsHandler := NewWebSocketHandler(clientOrigin, app)

	router := httprouter.New()
	dynamic := alice.New(app.logRequest, app.logResponse)

	// fileServer := http.FileServer(http.Dir("./ui/static/"))
	// mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/refresh", dynamic.ThenFunc(app.refreshToken))

	protected := dynamic.Append(auth.AuthMiddleware)
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogout))
	// Create WebSocket handler
	router.Handler(http.MethodGet, "/ws", protected.ThenFunc(wsHandler.Handle))

	return withCORS(clientOrigin)(router)
}
