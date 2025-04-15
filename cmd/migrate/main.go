package main

import (
	"fmt"

	"github.com/skrillatb/nemo/internal/db"
)

func main() {
	database, err := db.Connect()
	if err != nil {
		fmt.Println("Erreur de connexion à la DB:", err)
		return
	}
	defer database.Close()

	if err := db.Init(database); err != nil {
		fmt.Println("Erreur lors des migrations:", err)
		return
	}

	fmt.Println("Base initialisée avec succès.")
}
