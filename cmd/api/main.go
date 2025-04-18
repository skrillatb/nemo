package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"

	"github.com/joho/godotenv"
	"github.com/skrillatb/nemo/internal/db"
	"github.com/skrillatb/nemo/internal/handlers"
	"github.com/skrillatb/nemo/internal/middlewares"
)

var startTime = time.Now()

func main() {
	_ = godotenv.Load()

	// Init Sentry
	env := os.Getenv("APP_ENV")
	isProd := env == "prod"
	if isProd {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:              os.Getenv("SENTRY_DSN"),
			TracesSampleRate: 1.0,
			Environment:      env,
			BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
				if hint.Context != nil {
					if req, ok := hint.Context.Value(sentry.RequestContextKey).(*http.Request); ok {
						event.User.IPAddress = req.RemoteAddr
					}
				}
				return event
			},
		}); err != nil {
			panic("Sentry init failed: " + err.Error())
		}

		defer sentry.Flush(2 * time.Second)
	}

	// Init Zap logger
	logger, err := zap.NewProduction()
	if err != nil {
		panic("Zap init failed: " + err.Error())
	}
	defer logger.Sync()

	// DB
	database, err := db.Connect()
	if err != nil {
		logger.Fatal("DB error", zap.Error(err))
	}
	defer database.Close()

	app := &handlers.App{DB: database}

	// Init sentry middleware
	sentryHandler := sentryhttp.New(sentryhttp.Options{
		Repanic: true,
	})

	// Setup Chi router
	r := chi.NewRouter()
	r.Use(sentryHandler.Handle)
	r.Use(middlewares.ZapLoggerMiddleware(logger))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middlewares.SentryRecover(isProd))

	// Headers de sécurité
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Del("Server")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			w.Header().Set("Referrer-Policy", "no-referrer")
			next.ServeHTTP(w, r)
		})
	})

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", os.Getenv("PRODUCTION_URL")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Routes
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Health check", zap.String("remote", r.RemoteAddr))

		uptime := time.Since(startTime).Round(time.Second)
		start := startTime.Format(time.RFC3339)

		resp := map[string]interface{}{
			"status":  "🧣 Everything is fine — API is up!",
			"mood":    "Taylor Swift - All Too Well (10 Minute Version)",
			"link":    "https://open.spotify.com/track/5enxwA8aAbwZbf5qCHORXi",
			"uptime":  uptime.String(),
			"started": start,
			"version": "1.0.0",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	})

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/sites", func(r chi.Router) {
				r.Get("/", app.ListSites)
				r.Get("/{id}", app.GetSite)
				r.With(middlewares.RequireAuth(os.Getenv("API_TOKEN"))).Post("/", app.CreateSite)
				r.With(middlewares.RequireAuth(os.Getenv("API_TOKEN"))).Put("/{id}", app.UpdateSite)
				r.With(middlewares.RequireAuth(os.Getenv("API_TOKEN"))).Delete("/{id}", app.DeleteSite)
			})
			r.Get("/search", app.Search)
		})
	})

	// Serve
	logger.Info("Serveur démarré", zap.String("port", ":8080"))
	if err := http.ListenAndServe(":8080", r); err != nil {
		logger.Fatal("Erreur serveur", zap.Error(err))
	}
}
