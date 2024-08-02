package httpserver

import (
    "context"
    "fmt"
    "net/http"
    "slog"
)

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
