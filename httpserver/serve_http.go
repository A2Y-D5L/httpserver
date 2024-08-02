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
// are provided.
func Serve(ctx context.Context, options ...Option) error {
	server := New(ctx, options...)
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

func New(ctx context.Context, options ...func(*http.Server)) *http.Server {
	srv := &http.Server{
		ReadHeaderTimeout: envvar.GetDuration("HTTP_SERVER_READ_HEADER_TIMEOUT", 10*time.Second),
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
