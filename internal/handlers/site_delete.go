package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *App) DeleteSite(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	siteID, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	result, err := app.DB.Exec("DELETE FROM sites WHERE id = ?", siteID)
	if err != nil {
		http.Error(w, "Erreur suppression site", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Site introuvable", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
