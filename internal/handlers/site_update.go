package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/skrillatb/nemo/internal/storage"
)

func (app *App) UpdateSite(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	siteID, err := strconv.Atoi(idParam)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ID invalide"})
		return
	}

	site, fileHeader, err := BindAndUploadSite(r, false)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	var imageURL string

	if fileHeader != nil {
		file, err := fileHeader.Open()
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Impossible d'ouvrir l'image"})
			return
		}
		defer file.Close()

		filename := strings.ReplaceAll(site.Name, " ", "_") + filepath.Ext(fileHeader.Filename)

		err = storage.UploadToBunny(file, filename)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Échec upload image : " + err.Error()})
			return
		}

		pullZoneURL := os.Getenv("BUNNY_PULL_ZONE_URL")
		imageURL = strings.TrimRight(pullZoneURL, "/") + "/" + filename
	} else {
		err = app.DB.QueryRow(`SELECT image_url FROM sites WHERE id = ?`, siteID).Scan(&imageURL)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Erreur récupération image existante"})
			return
		}
	}

	_, err = app.DB.Exec(`
		UPDATE sites
		SET name = ?, site_url = ?, image_url = ?, language = ?, ads = ?, type = ?, updated_at = ?
		WHERE id = ?`,
		site.Name, site.SiteURL, imageURL, site.Language, site.Ads, site.Type, time.Now(), siteID)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Erreur mise à jour site : " + err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"success": "true",
	})
}
