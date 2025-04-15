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
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	site, fileHeader, err := BindAndUploadSite(r, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var imageURL string

	if fileHeader != nil {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Impossible d'ouvrir l'image", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		filename := strings.ReplaceAll(site.Name, " ", "_") + filepath.Ext(fileHeader.Filename)

		err = storage.UploadToBunny(file, filename)
		if err != nil {
			http.Error(w, "Échec upload image : "+err.Error(), http.StatusInternalServerError)
			return
		}

		pullZoneURL := os.Getenv("BUNNY_PULL_ZONE_URL")
		imageURL = strings.TrimRight(pullZoneURL, "/") + "/" + filename
	} else {
		err = app.DB.QueryRow(`SELECT image_url FROM sites WHERE id = ?`, siteID).Scan(&imageURL)
		if err != nil {
			http.Error(w, "Erreur récupération image existante", http.StatusInternalServerError)
			return
		}
	}

	_, err = app.DB.Exec(`
		UPDATE sites
		SET name = ?, site_url = ?, image_url = ?, language = ?, ads = ?, type = ?, updated_at = ?
		WHERE id = ?`,
		site.Name, site.SiteURL, imageURL, site.Language, site.Ads, site.Type, time.Now(), siteID)

	if err != nil {
		http.Error(w, "Erreur mise à jour site : "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "site mis à jour",
	})
}
