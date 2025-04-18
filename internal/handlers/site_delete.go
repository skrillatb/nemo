package handlers

import (
	"encoding/json"
	"net/http"
	"path"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/skrillatb/nemo/internal/storage"
)

func (app *App) DeleteSite(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	siteID, err := strconv.Atoi(idParam)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ID invalide"})
		return
	}

	var imageURL string
	err = app.DB.QueryRow(`SELECT image_url FROM sites WHERE id = ?`, siteID).Scan(&imageURL)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Site introuvable ou erreur DB"})
		return
	}

	if imageURL != "" {
		filename := path.Base(imageURL)
		err := storage.DeleteFromBunny(filename)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Erreur suppression image Bunny"})
			return
		}
	}

	result, err := app.DB.Exec("DELETE FROM sites WHERE id = ?", siteID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Erreur suppression BDD"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Site non trouvé"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
