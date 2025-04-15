package storage

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func DeleteFromBunny(filename string) error {
	storageZone := os.Getenv("BUNNY_STORAGE_ZONE")
	region := os.Getenv("BUNNY_STORAGE_REGION")
	apiKey := os.Getenv("BUNNY_STORAGE_KEY")

	if storageZone == "" || apiKey == "" {
		return fmt.Errorf("clé Bunny ou storage zone manquante")
	}

	baseHost := "storage.bunnycdn.com"
	if region != "" {
		baseHost = region + "." + baseHost
	}

	filename = strings.TrimPrefix(filename, "/")
	url := fmt.Sprintf("https://%s/%s/%s", baseHost, storageZone, filename)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("AccessKey", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("échec suppression Bunny : code %d", resp.StatusCode)
	}

	return nil
}
