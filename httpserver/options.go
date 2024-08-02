package httpserver

import (
	"crypto/tls"
	"net/http"
	"slog"
	"time"
)

// ErrorLog returns a closure that sets the value of http.Server.ErrorLog.
//
// http.Server.ErrorLog specifies an optional logger for errors accepting connections, unexpected behavior from handlers, and underlying FileSystem errors. If nil, logging is done via the log package's standard logger.
func ErrorLog(logger *slog.Logger) func(*http.Server) {
	return func(srv *http.Server) {
		srv.ErrorLog = slog.NewLogLogger(logger.Handler(), slog.LevelError)
	}
}

// Address sets http.Server.Addr to addr.
//
// http.Server.Addr optionally specifies the TCP address for the server to listen on, in the form "host:port". If empty, ":http" (port 80) is used. The service names are defined in RFC 6335 and assigned by IANA. See net.Dial for details of the address format.
func Address(addr string) func(*http.Server) {
	return func(srv *http.Server) {
		srv.Addr = addr
	}
}

// TLSConfig sets http.Server.TLSConfig.
//
// http.Server.TLSConfig optionally provides a TLS configuration for use by ServeTLS and ListenAndServeTLS. Note that this value is cloned by ServeTLS and ListenAndServeTLS, so it's not possible to modify the configuration with methods like tls.Config.SetSessionTicketKeys. To use SetSessionTicketKeys, use Server.Serve with a TLS Listener instead.
func TLSConfig(tlsConfig *tls.Config) func(*http.Server) {
	return func(srv *http.Server) {
		srv.TLSConfig = tlsConfig
	}
}

// MaxHeaderBytes sets http.Server.MaxHeaderBytes
//
// http.Server.MaxHeaderBytes controls the maximum number of bytes the server will read parsing the request header's keys and values, including the request line. It does not limit the size of the request body. If zero, DefaultMaxHeaderBytes is used.
func MaxHeaderBytes(maxHeaderBytes int) func(*http.Server) {
	return func(srv *http.Server) {
		srv.MaxHeaderBytes = maxHeaderBytes
	}
}

// ReadHeaderTimeout sets http.Server.ReadHeaderTimeout.
//
// http.Server.ReadHeaderTimeout is the amount of time allowed to read request headers. The connection's read deadline is reset after reading the headers and the Handler can decide what is considered too slow for the body. If ReadHeaderTimeout is zero, the value of ReadTimeout is used. If both are zero, there is no timeout.
func ReadHeaderTimeout(seconds int) func(*http.Server) {
	return func(srv *http.Server) {
		srv.ReadHeaderTimeout = time.Duration(seconds) * time.Second
	}
}
