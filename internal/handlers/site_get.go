package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *App) GetSite(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	siteID, err := strconv.Atoi(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ID invalide"})
		return
	}

	var site Site
	err = app.DB.QueryRow(`
		SELECT id, name, site_url, image_url, language, ads, type, created_at, updated_at
		FROM sites
		WHERE id = ?`, siteID).Scan(&site.ID, &site.Name, &site.SiteURL, &site.ImageURL, &site.Language, &site.Ads, &site.Type, &site.CreatedAt, &site.UpdatedAt)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Site non trouvé"})
		return
	}

	json.NewEncoder(w).Encode(site)
}
