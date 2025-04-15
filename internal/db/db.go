package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func Connect() (*sql.DB, error) {
	err := godotenv.Load() 
	if err != nil {
		return nil, fmt.Errorf("erreur de chargement du .env : %w", err)
	}

	dbURL := os.Getenv("TURSO_DATABASE_URL")
	authToken := os.Getenv("TURSO_AUTH_TOKEN")

	if dbURL == "" || authToken == "" {
		return nil, fmt.Errorf("TURSO_DATABASE_URL ou TURSO_AUTH_TOKEN manquant")
	}

	dsn := fmt.Sprintf("%s?authToken=%s", dbURL, authToken)
	db, err := sql.Open("libsql", dsn)
	if err != nil {
		return nil, fmt.Errorf("connexion échouée : %w", err)
	}

	return db, nil
}

func Init(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS sites (
			id INTEGER PRIMARY KEY,
			name TEXT UNIQUE,
			site_url TEXT,
			image_url TEXT,
			language TEXT,
			ads BOOLEAN DEFAULT true,
			type TEXT DEFAULT 'null',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("erreur création table : %w", err)
	}
	return nil
}