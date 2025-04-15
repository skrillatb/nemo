package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"github.com/skrillatb/nemo/internal/db"
	"github.com/skrillatb/nemo/internal/handlers"
	mw "github.com/skrillatb/nemo/internal/middleware"
)

func main() {
	_ = godotenv.Load()

	apiToken := os.Getenv("API_TOKEN")
	if apiToken == "" {
		fmt.Println("API_TOKEN manquant dans .env")
		return
	}

	database, err := db.Connect()
	if err != nil {
		fmt.Println("Erreur DB:", err)
		return
	}
	defer database.Close()

	app := &handlers.App{DB: database}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", os.Getenv("PRODUCTION_URL")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
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
				r.With(mw.RequireAuth(apiToken)).Post("/", app.CreateSite)
				r.With(mw.RequireAuth(apiToken)).Put("/{id}", app.UpdateSite)
				r.With(mw.RequireAuth(apiToken)).Delete("/{id}", app.DeleteSite)
			})

			r.Get("/search", app.Search)

		})

	})

	http.ListenAndServe(":3000", r)
}
