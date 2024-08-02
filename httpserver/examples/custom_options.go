package examples

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/a2y-d5l/serve/httpserver"
)

// DisableGeneralOptionsHandler defines a custom functional option for
// configuring the server.
func DisableGeneralOptionsHandler() httpserver.Option {
	return func(s *http.Server) {
		s.DisableGeneralOptionsHandler = true
	}
}

// EXAMPLE: demonstrates how to run an HTTP server with a custom option.
func customOptionsServer() {
	if err := httpserver.Serve(
		context.Background(),
		httpserver.Address(":8080"),    // provided by the httpserver package
		DisableGeneralOptionsHandler(), // custom option
	); err != nil {
		slog.Error("httpserver.Serve error:" + err.Error())
	}
}
