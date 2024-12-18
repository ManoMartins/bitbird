package handlers

import (
	"context"
	"log"

	"github.com/ManoMartins/bitbird/configs"
	"github.com/gofiber/fiber/v2"
)

func Status(c *fiber.Ctx) error {
	configs.InitDatabase()
	defer configs.CloseDatabase()

	var sum int
	err := configs.DB.QueryRow(context.Background(), "SELECT 1 + 1;").Scan(&sum)

	if err != nil {
		log.Fatalf("Erro ao realizar consulta: %v", err)
	}

	log.Print(sum)

	return c.JSON(fiber.Map{
		"message": "oi",
	})
}
