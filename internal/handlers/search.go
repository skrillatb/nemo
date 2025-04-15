package handlers

import (
	"encoding/json"
	"net/http"
)

// Exemple de résultat de recherche pas encore stable
type SearchResult struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	PosterURL string `json:"poster_url"`
	StreamURL string `json:"stream_url"`
}

func (app *App) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Query manquant", http.StatusBadRequest)
		return
	}

	// TODO(long term): à remplacer par une vraie recherche plus tard
	results := []SearchResult{
		{
			ID:        1,
			Name:      "Test Result",
			PosterURL: "https://example.com/poster.jpg",
			StreamURL: "https://example.com/stream",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}
