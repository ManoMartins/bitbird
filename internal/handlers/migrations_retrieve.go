package handlers

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-migrate/migrate/v4"

	// Required for migration drivers
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type MigrationsRetrieveOutput struct {
	PendingMigrations []string `json:"pending_migrations"`
}

func MigrationsRetrieve(c *fiber.Ctx) error {
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

	currentVersion, _, _ := m.Version()

	files, err := os.ReadDir("internal/database/migrations")
	if err != nil {
		log.Println("Error reading migration files:", err)
		return err
	}

	var pendingFiles []string
	re := regexp.MustCompile(`^(\d{6})_.+\.up\.sql$`)
	currentVersionStr := fmt.Sprintf("%06d", currentVersion)

	for _, file := range files {
		if !file.IsDir() && re.MatchString(file.Name()) {
			match := re.FindStringSubmatch(file.Name())
			if len(match) > 1 {
				fileVersion := match[1]
				if strings.Compare(fileVersion, currentVersionStr) > 0 {
					pendingFiles = append(pendingFiles, file.Name())
				}
			}
		}
	}

	return c.JSON(MigrationsRetrieveOutput{
		PendingMigrations: pendingFiles,
	})
}
