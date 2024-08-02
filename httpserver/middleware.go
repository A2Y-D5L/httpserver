package httpserver

import (
	"fmt"
	"log/slog"
	"net/http"
)

func handlePanic(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				var errMsg string
				if e, ok := err.(error); ok {
					errMsg = e.Error()
				} else {
					errMsg = fmt.Sprintf("%v", err)
				}
				slog.Error("panicked while handling " + r.Method + " " + r.URL.Path + ":" + errMsg)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)
	})
}
