package httpserver

import (
    "net/http"
)

// Routes defines the routes and middleware for a server
type Routes struct {
	PathHandlers map[string]http.Handler
	Middleware   []func(http.Handler) http.Handler
}
