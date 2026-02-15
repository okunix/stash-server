package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func Logger() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()
			writer := &wrappedWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(writer, r)
			latency := time.Since(now).Microseconds()
			slog.Info(
				"incoming request",
				"path", r.URL.Path,
				"method", r.Method,
				"statusCode", writer.statusCode,
				"latency", latency,
				"remoteAddr", r.RemoteAddr,
			)
		})
	}
}
