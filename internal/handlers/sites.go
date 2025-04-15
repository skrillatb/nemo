package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type App struct {
	DB *sql.DB
}

type Site struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	SiteURL   string    `json:"site_url"`
	ImageURL  string    `json:"image_url"`
	Language  string    `json:"language"`
	Ads       bool      `json:"ads"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (app *App) ListRecentSites(w http.ResponseWriter, r *http.Request) {
	rows, err := app.DB.Query(`
		SELECT id, name, site_url, image_url, language, ads, type, created_at, updated_at
		FROM sites
		ORDER BY updated_at DESC
		LIMIT 50
	`)
	if err != nil {
		http.Error(w, "Erreur requête SQL", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var sites []Site

	for rows.Next() {
		var s Site
		err := rows.Scan(&s.ID, &s.Name, &s.SiteURL, &s.ImageURL, &s.Language, &s.Ads, &s.Type, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			http.Error(w, "Erreur lecture des données", http.StatusInternalServerError)
			return
		}
		sites = append(sites, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sites)
}
