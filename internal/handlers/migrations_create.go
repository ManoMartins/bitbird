package handlers

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-migrate/migrate/v4"
)

func MigrationCreate(c *fiber.Ctx) error {
	m, err := migrate.New("file://internal/database/migrations", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Println("Error initializing migration:", err)
		return err
	}

	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			log.Println("Error closing migration source:", srcErr)
		}
		if dbErr != nil {
			log.Println("Error closing database connection:", dbErr)
		}
	}()

	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			return c.JSON(fiber.Map{
				"message": "No pending migrations",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error applying migrations",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Migrations applied successfully",
	})
}
