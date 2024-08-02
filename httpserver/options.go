package httpserver

import (
	"crypto/tls"
	"net/http"
	"slog"
	"time"
)

// Option is a functional option for configuring an http.Server.
//
// The Option type is just an alias for a function that takes an *http.Server.
// It is only purpose is to provice clarity that a function is intended to be
// used as a functional option. Using the Option type is not required to define
// custom option functions. For example, the following two function signatures
// are equivalent:
//
//	// Using the Option type
//	func MyCustomOption() Option {
//		// Return a closure to modify the server
//		return func(s *http.Server) {
//			// do something with srv
//		}
//	}
//
//	// Not using the Option type
//	func MyCustomOption() func(*http.Server) {
//		// Return a closure to modify the server
//		return func(s *http.Server) {
//			// do something with srv
//		}
//	}
type Option = func(*http.Server)

// ErrorLog specifies an optional logger for errors accepting connections,
// unexpected behavior from handlers, and underlying FileSystem errors.
func ErrorLog(logger *slog.Logger) func(*http.Server) {
	return func(s *http.Server) {
		s.ErrorLog = slog.NewLogLogger(logger.Handler(), slog.LevelError)
	}
}

// Address specifies the TCP address for the server to listen on, in the form
// "host:port".
func Address(addr string) func(*http.Server) {
	return func(s *http.Server) {
		s.Addr = addr
	}
}

// Routes sets http.Server.Handler to an http.ServeMux configured to handle the
// provided routes. If http.Server.Handler is already set, the routes are added
// to the existing http.ServeMux. If no routes are specified, the server will
// use http.DefaultServeMux.
func Routes(routes ...Route) func(*http.Server) {
	return func(s *http.Server) {
		var mux *http.ServeMux
		if s.Handler == nil {
			mux = http.NewServeMux()
		} else {
			mux = s.Handler.(*http.ServeMux)
		}
		for _, route := range routes {
			handler := route.Handler
			for i := len(route.Middleware) - 1; i >= 0; i-- {
				handler = route.Middleware[i](handler)
			}
			mux.Handle(route.Pattern, handler)
		}
		s.Handler = mux
	}
}

// TLSConfig provides a TLS configuration for use by ServeTLS and
// ListenAndServeTLS.
func TLSConfig(tlsConfig *tls.Config) func(*http.Server) {
	return func(s *http.Server) {
		s.TLSConfig = tlsConfig
	}
}

// MaxHeaderBytes controls the maximum number of bytes the server will read
// parsing the request header's keys and values, including the request line.
//
// If not specified, or set to zero, a default of 1 MiB is used.
func MaxHeaderBytes(maxHeaderBytes int) func(*http.Server) {
	return func(s *http.Server) {
		s.MaxHeaderBytes = maxHeaderBytes
	}
}

// ReadHeaderTimeout is the amount of time allowed to read request headers.
func ReadHeaderTimeout(seconds int) func(*http.Server) {
	return func(s *http.Server) {
		s.ReadHeaderTimeout = time.Duration(seconds) * time.Second
	}
}
