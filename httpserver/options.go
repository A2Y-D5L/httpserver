package httpserver




type Option func(*http.Server)

// WithAddr sets the address the server will listen on.
func WithAddr(addr string) Option {
	return func(srv *http.Server) {
		srv.Addr = addr
	}
}

// WithTLSConfig sets the server's TLS configuration.
func WithTLSConfig(tlsConfig *tls.Config) Option {
	return func(srv *http.Server) {
		srv.TLSConfig = tlsConfig
	}
}

// WithMaxHeaderBytes sets the server's max header bytes.
//
// MaxHeaderBytes controls the maximum number of bytes the server will read
// parsing the request header's keys and values, including the request line. It
// does not limit the size of the request body.
func WithMaxHeaderBytes(maxHeaderBytes int) Option {
	return func(srv *http.Server) {
		srv.MaxHeaderBytes = maxHeaderBytes
	}
}
