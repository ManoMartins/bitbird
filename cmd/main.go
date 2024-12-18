package main

import (
	"log"

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
	app := fiber.New()

	api := app.Group("api/v1")

	api.Get("/status", handlers.Status)

	log.Fatal(app.Listen(":8080"))
}
