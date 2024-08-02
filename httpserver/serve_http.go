package httpserver

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/a2y-d5l/serve/internal/envvar"
)

type Route struct {
	Pattern    string
	Handler    http.Handler
	Middleware []func(http.Handler) http.Handler
}

// Serve HTTP requests for the provided routes. The server will gracefully shut
// down when the provided context is done or a SIGINT or SIGTERM signal is
// received.
//
// The server is automatically configured with sensible defaults that can be 
// overridden by functional options. The options are applied in the order they
// are provided. Several options are provided:
//
// - Address: sets the TCP address for the server to listen on, in the form
//   "host:port". Default: the value of the HTTP_SERVER_ADDR environment 
//   variable or ":8080".
//
// - ErrorLog: sets an optional logger for errors accepting connections,
//   unexpected behavior from handlers, and underlying FileSystem errors. 
//   Default: the default slog.Logger.
//
// - ReadHeaderTimeout: sets the maximum duration for reading the entire request
//   header. Default: 10 seconds.
//
// - TLSConfig: provides a TLS configuration for use by ServeTLS and
//   ListenAndServeTLS. Default: nil.
//
// - MaxHeaderBytes: controls the maximum number of bytes the server will read
//   parsing the request header's keys and values, including the request line.
//   Default: 1 MiB.
//
// Defining Custom Options:
//
// Custom options can be created by defining a closure that
// accepts a pointer to an http.Server and returns a function that sets a field
// on the server. For example, to set the server's address:
//
//     func Address(addr string) func(*http.Server) {
//         return func(srv *http.Server) {
//             srv.Addr = addr
//         }
//     }
func Serve(ctx context.Context, routes []Route, options ...func(*http.Server)) error {
	server := newServer(ctx, routes, options...)
	errCh := make(chan error)
	go func() {
		if server.TLSConfig != nil {
			errCh <- server.ListenAndServeTLS("", "")
        } else {
			errCh <- server.ListenAndServe()
		}
	}()
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := strings.ToUpper((<-shutdownCh).String())
		slog.Debug(sig + " signal received, shutting down server")
		ctx, cancel := context.WithTimeout(ctx, 29*time.Second)
		defer cancel()
		errCh <- server.Shutdown(ctx)
	}()
	for {
		select {
		case <-ctx.Done():
			slog.Info("context done, shutting down server")
			shutdownCh <- syscall.SIGTERM
		case err := <-errCh:
			if err != http.ErrServerClosed {
				return fmt.Errorf("graceful shutdown failed: %w", err)
			}
			return nil
		}
	}
}

func newServer(ctx context.Context, routes []Route, options ...func(*http.Server)) *http.Server {
	mux := http.NewServeMux()
	for _, r := range routes {
		for _, mw := range r.Middleware {
			r.Handler = mw(r.Handler)
		}
	}
	srv := &http.Server{
		ReadHeaderTimeout: envvar.GetDuration("HTTP_SERVER_READ_HEADER_TIMEOUT", 10*time.Second),
		Handler:           handlePanic(mux),
		ErrorLog:          slog.NewLogLogger(slog.Default().Handler(), slog.LevelError),
        Addr:              envvar.Get("HTTP_SERVER_ADDR", ":8080"),
        MaxHeaderBytes:    envvar.GetInt("HTTP_SERVER_MAX_HEADER_BYTES", http.DefaultMaxHeaderBytes),
		BaseContext: func(net.Listener) context.Context {
			if ctx == nil {
				return context.Background()
			}
			return ctx
		},
	}
	for _, opt := range options {
		opt(srv)
	}
	return srv
}
