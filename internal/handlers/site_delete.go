package handlers

import (
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
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	var imageURL string
	err = app.DB.QueryRow(`SELECT image_url FROM sites WHERE id = ?`, siteID).Scan(&imageURL)
	if err != nil {
		http.Error(w, "Site introuvable ou erreur DB", http.StatusNotFound)
		return
	}

	if imageURL != "" {
		filename := path.Base(imageURL)
		err := storage.DeleteFromBunny(filename)
		if err != nil {
			http.Error(w, "Erreur suppression image Bunny", http.StatusInternalServerError)
			return
		}
	}

	result, err := app.DB.Exec("DELETE FROM sites WHERE id = ?", siteID)
	if err != nil {
		http.Error(w, "Erreur suppression BDD", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Site non trouvé", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
