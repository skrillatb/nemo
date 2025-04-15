package main

import (
	"net/http"
	"os"

	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"github.com/skrillatb/nemo/internal/db"
	"github.com/skrillatb/nemo/internal/handlers"
	"github.com/skrillatb/nemo/internal/middlewares"
)

func main() {
	_ = godotenv.Load()

	// Initialisation du logger
	logger, err := zap.NewProduction()
	if err != nil {
		panic("Échec init zap logger: " + err.Error())
	}
	defer logger.Sync()

	apiToken := os.Getenv("API_TOKEN")
	if apiToken == "" {
		logger.Fatal("API_TOKEN manquant dans .env")
	}

	database, err := db.Connect()
	if err != nil {
		logger.Fatal("Erreur de connexion à la base de données", zap.Error(err))
	}
	defer database.Close()

	app := &handlers.App{DB: database}

	r := chi.NewRouter()

	// Middleware Zap Logger
	r.Use(middlewares.ZapLoggerMiddleware(logger))

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	// Headers de sécurité
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Supprime l'empreinte du serveur
			w.Header().Del("Server")

			// Headers de sécurité
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			w.Header().Set("Referrer-Policy", "no-referrer") 
			next.ServeHTTP(w, r)
		})
	})

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", os.Getenv("PRODUCTION_URL")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Health check hit", zap.String("remote", r.RemoteAddr))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"status": "🧣 Everything is fine — API is up!",
			"mood": "Taylor Swift - All Too Well (10 Minute Version)",
			"link": "https://open.spotify.com/track/5enxwA8aAbwZbf5qCHORXi"
		}`))
	})

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/sites", func(r chi.Router) {
				r.Get("/", app.ListSites)
				r.With(middlewares.RequireAuth(apiToken)).Post("/", app.CreateSite)
				r.With(middlewares.RequireAuth(apiToken)).Put("/{id}", app.UpdateSite)
				r.With(middlewares.RequireAuth(apiToken)).Delete("/{id}", app.DeleteSite)
			})
			r.Get("/search", app.Search)
		})
	})

	logger.Info("Serveur démarré", zap.String("port", ":8080"))
	if err := http.ListenAndServe(":8080", r); err != nil {
		logger.Fatal("Erreur au démarrage du serveur", zap.Error(err))
	}
}
