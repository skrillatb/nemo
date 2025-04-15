package main

import (
	"fmt"
	"net/http"

	"github.com/skrillatb/nemo/internal/db"
	"github.com/skrillatb/nemo/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
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

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Get("/sites", app.ListRecentSites)

	http.ListenAndServe(":3000", r)
}
