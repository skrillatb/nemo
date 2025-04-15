package middlewares

import (
	"net/http"
	"time"

	"go.uber.org/zap"
	"github.com/go-chi/chi/v5/middleware"
)

func ZapLoggerMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			duration := time.Since(start)

			logger.Info("HTTP Request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", ww.Status()),
				zap.String("remote", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()),
				zap.Duration("duration", duration),
			)
		})
	}
}
