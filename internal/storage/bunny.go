package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func UploadToBunny(file multipart.File, filename string) error {
	storageZone := os.Getenv("BUNNY_STORAGE_ZONE")
	region := os.Getenv("BUNNY_STORAGE_REGION")
	apiKey := os.Getenv("BUNNY_STORAGE_KEY")

	if storageZone == "" || apiKey == "" {
		return fmt.Errorf("clé Bunny manquante")
	}

	baseHost := "storage.bunnycdn.com"
	if region != "" {
		baseHost = region + "." + baseHost
	}

	url := fmt.Sprintf("https://%s/%s/%s", baseHost, storageZone, filename)

	req, err := http.NewRequest("PUT", url, file)
	if err != nil {
		return err
	}

	req.Header.Set("AccessKey", apiKey)
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("erreur upload Bunny: %s", string(body))
	}

	return nil
}
