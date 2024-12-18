package main

import (
	"log"

	"github.com/ManoMartins/bitbird/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	api := app.Group("api/v1")

	api.Get("/status", handlers.Status)

	log.Fatal(app.Listen(":8080"))
}
