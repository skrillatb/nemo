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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Impossible d'ouvrir le fichier"})
		return
	}
	defer file.Close()

	filename := strings.ReplaceAll(site.Name, " ", "_") + filepath.Ext(fileHeader.Filename)

	err = storage.UploadToBunny(file, filename)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Upload Bunny échoué : " + err.Error()})
		return
	}

	pullZoneURL := os.Getenv("BUNNY_PULL_ZONE_URL")
	site.ImageURL = strings.TrimRight(pullZoneURL, "/") + "/" + filename

	site.CreatedAt = time.Now()
	site.UpdatedAt = time.Now()

	_, err = app.DB.Exec(`
		INSERT INTO sites (name, site_url, image_url, language, ads, type, created_at, updated_at, hidden)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		site.Name, site.SiteURL, site.ImageURL, site.Language, site.Ads, site.Type, site.CreatedAt, site.UpdatedAt, site.Hidden)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Erreur DB : " + err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"success": "true",
	})
}
