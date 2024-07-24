package main

import (
	"hotel_management_system/database"
	"hotel_management_system/routes"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

// @title Hotel Management System API
// @version 1.0
// @description This is a sample server for a hotel management system.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.Connect()
	database.Migrate()

	r := routes.InitRouter()

	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
