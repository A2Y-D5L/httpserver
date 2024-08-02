package examples

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/a2y-d5l/serve/httpserver"
)

// handleGETRoot is an HTTP handler that writes "Hello, World!" to the response.
func handleGETRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

// EXAMPLE: demonstrates how to run a simple HTTP server.
func basicServerExample() {
	if err := httpserver.Serve(
		context.Background(),
		httpserver.Routes(
			httpserver.Route{
				Pattern: "GET /",
				Handler: http.HandlerFunc(handleGETRoot),
			},
		),
	); err != nil {
		slog.Error("httpserver.Serve error:" + err.Error())
	}
}
