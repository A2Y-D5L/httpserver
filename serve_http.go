package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Routes struct {
	PathHandlers map[string]http.Handler
	Middleware   []func(http.Handler) http.Handler
}

type ServerOption func(*http.Server)

// WithAddr sets the address the server will listen on.
func WithAddr(addr string) ServerOption {
	return func(srv *http.Server) {
		srv.Addr = addr
	}
}

// WithTLSConfig sets the server's TLS configuration.
func WithTLSConfig(tlsConfig *tls.Config) ServerOption {
	return func(srv *http.Server) {
		srv.TLSConfig = tlsConfig
	}
}

// WithMaxHeaderBytes sets the server's max header bytes.
//
// MaxHeaderBytes controls the maximum number of bytes the server will read
// parsing the request header's keys and values, including the request line. It
// does not limit the size of the request body.
func WithMaxHeaderBytes(maxHeaderBytes int) ServerOption {
	return func(srv *http.Server) {
		srv.MaxHeaderBytes = maxHeaderBytes
	}
}

func ServeHTTP(ctx context.Context, routes Routes, logger *slog.Logger, serverOptions ...ServerOption) error {
	if logger == nil {
		logger = slog.Default()
	}
	mux := setupMux(ctx, routes, logger)
	srv := setupServer(ctx, mux, logger, serverOptions...)
	return startServer(ctx, srv, logger)
}

func setupMux(ctx context.Context, routes Routes, logger *slog.Logger) *http.ServeMux {
	mux := http.NewServeMux()
	for path, handler := range routes.PathHandlers {
		for _, mw := range routes.Middleware {
			handler = mw(handler)
		}
		mux.Handle(path, handler)
	}
	mux.Handle(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				var errMsg string
				if e, ok := err.(error); ok {
					errMsg = e.Error()
				} else {
					errMsg = fmt.Sprintf("%v", err)
				}
				logger.ErrorContext(ctx, "panicked while handling "+r.Method+" "+r.URL.Path+":"+errMsg)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
	}))
	return mux.ServeHTTP()
}

func handlePanic(ctx context.Context, logger *slog.Logger, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				var errMsg string
				if e, ok := err.(error); ok {
					errMsg = e.Error()
				} else {
					errMsg = fmt.Sprintf("%v", err)
				}
				logger.ErrorContext(ctx, "panicked while handling "+r.Method+" "+r.URL.Path+":"+errMsg)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		handler.ServeHTTP(w, r)
	})
}

func setupServer(ctx context.Context, mux *http.ServeMux, logger *slog.Logger, serverOptions ...ServerOption) *http.Server {
	srv := &http.Server{
		Addr:              getEnvVar("SERVER_ADDR", ":8080"),
		ReadHeaderTimeout: getDurationEnvVar("READ_HEADER_TIMEOUT", 10*time.Second),
		Handler:           mux,
		ErrorLog:          slog.NewLogLogger(logger.Handler(), slog.LevelError),
		BaseContext: func(net.Listener) context.Context {
			if ctx == nil {
				return context.Background()
			}
			return ctx
		},
	}
	for _, opt := range serverOptions {
		opt(srv)
	}
	return srv
}

func getEnvVar(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getDurationEnvVar(key string, fallback time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return fallback
}

func startServer(ctx context.Context, srv *http.Server, logger *slog.Logger) error {
	errCh := make(chan error)
	go func() {
		slog.Info("server listening on " + srv.Addr)
		switch {
		case srv.TLSConfig != nil:
			errCh <- srv.ListenAndServeTLS("", "")
		default:
			errCh <- srv.ListenAndServe()
		}
	}()
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := strings.ToUpper((<-shutdownCh).String())
		logger.Debug(sig + " signal received, shutting down server")
		ctx, cancel := context.WithTimeout(ctx, 29*time.Second)
		defer cancel()
		errCh <- srv.Shutdown(ctx)
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
