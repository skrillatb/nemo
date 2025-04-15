package handlers

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
)

func BindAndUploadSite(r *http.Request, requireImage bool) (Site, *multipart.FileHeader, error) {
	var site Site

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return site, nil, fmt.Errorf("fichier trop gros ou mauvais format")
	}

	if err := json.Unmarshal([]byte(r.FormValue("data")), &site); err != nil {
		return site, nil, fmt.Errorf("erreur JSON : %v", err)
	}

	validate := validator.New()
	if err := validate.Struct(site); err != nil {
		var errors []string
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, err.Field()+": "+err.Tag())
		}
		return site, nil, fmt.Errorf("champs invalides : %s", strings.Join(errors, ", "))
	}

	_, fileHeader, err := r.FormFile("image")
	if err != nil {
		if requireImage {
			return site, nil, fmt.Errorf("fichier manquant")
		}
		return site, nil, nil
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		return site, nil, fmt.Errorf("extension non autorisée")
	}

	return site, fileHeader, nil
}
