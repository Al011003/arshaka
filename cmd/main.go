package main

import (
	"log"

	"backend/app"
)

func main() {
	application, err := app.NewApp()
	if err != nil {
		log.Fatalf("Gagal start app: %v", err)
	}

	// Start Server
	if err := application.Router.Run(":8080"); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
