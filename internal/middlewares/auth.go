package middlewares

import (
	"net/http"
	"strings"
)

func RequireAuth(expectedToken string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Token manquant ou invalide", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token != expectedToken {
				http.Error(w, "Token invalide", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}