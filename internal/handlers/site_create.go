package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/skrillatb/nemo/internal/storage"
)

func (app *App) CreateSite(w http.ResponseWriter, r *http.Request) {
	site, fileHeader, err := BindAndUploadSite(r, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		http.Error(w, "Impossible d'ouvrir le fichier", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	filename := strings.ReplaceAll(site.Name, " ", "_") + filepath.Ext(fileHeader.Filename)

	err = storage.UploadToBunny(file, filename)
	if err != nil {
		http.Error(w, "Upload Bunny échoué : "+err.Error(), http.StatusInternalServerError)
		return
	}

	pullZoneURL := os.Getenv("BUNNY_PULL_ZONE_URL")
	site.ImageURL = strings.TrimRight(pullZoneURL, "/") + "/" + filename

	site.CreatedAt = time.Now()
	site.UpdatedAt = time.Now()

	_, err = app.DB.Exec(`
		INSERT INTO sites (name, site_url, image_url, language, ads, type, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		site.Name, site.SiteURL, site.ImageURL, site.Language, site.Ads, site.Type, site.CreatedAt, site.UpdatedAt)

	if err != nil {
		http.Error(w, "Erreur DB : "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(site)
}
