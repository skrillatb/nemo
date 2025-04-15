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

func (app *App) ListSites(w http.ResponseWriter, r *http.Request) {
	allowedFilters := map[string]string{
		"language": "language = ?",
		"type":     "type = ?",
		"ads":      "ads = ?",
	}

	query := `SELECT id, name, site_url, image_url, language, ads, type, created_at, updated_at FROM sites WHERE 1=1`
	var args []interface{}

	for key, condition := range allowedFilters {
		value := r.URL.Query().Get(key)
		if value == "" {
			continue
		}

		if key == "ads" {
			if value == "true" {
				value = "1"
			} else if value == "false" {
				value = "0"
			} else {
				http.Error(w, "Paramètre ads invalide (true/false)", http.StatusBadRequest)
				return
			}
		}

		query += " AND " + condition
		args = append(args, value)
	}

	query += " ORDER BY updated_at DESC LIMIT 50"

	rows, err := app.DB.Query(query, args...)
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
			http.Error(w, "Erreur lecture", http.StatusInternalServerError)
			return
		}
		sites = append(sites, s)
	}

	if len(sites) == 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Site{})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sites)
}


