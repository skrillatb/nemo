package middlewares

import (
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
)

func SentryRecover(sentryEnabled bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					if sentryEnabled {
						if hub := sentry.GetHubFromContext(r.Context()); hub != nil {
							hub.Recover(err)
							hub.Flush(time.Second * 2)
						}
					}

					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
