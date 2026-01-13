package main

import (
	"log"

	"backend/app"
	_ "backend/docs" // Import docs yang bakal di-generate
)

// @title           Backend API
// @version         1.0
// @description     API Documentation untuk Frontend
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @BasePath  /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description JWT token tanpa prefix Bearer

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