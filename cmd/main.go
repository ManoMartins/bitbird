package main

import (
	"log"

	"github.com/ManoMartins/bitbird/configs"
	"github.com/ManoMartins/bitbird/internal/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func init() {
	if envErr := godotenv.Load(); envErr != nil {
		log.Fatal(".env file missing")
	}
}

func main() {
	configs.InitDatabase()
	defer configs.CloseDatabase()

	app := fiber.New()

	api := app.Group("api/v1")

	api.Get("/status", handlers.Status)
	api.Get("/migrations", handlers.MigrationsRetrieve)
	api.Post("/migrations", handlers.MigrationCreate)

	log.Fatal(app.Listen(":8080"))
}
