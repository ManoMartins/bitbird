package handlers

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ManoMartins/bitbird/configs"
	"github.com/gofiber/fiber/v2"
)

type Database struct {
	Version         string `json:"version"`
	MaxConnections  int    `json:"max_connections"`
	UsedConnections int    `json:"used_connections"`
}

type Dependencies struct {
	Database Database `json:"database"`
}

type Response struct {
	UpdatedAt    string       `json:"updated_at"`
	Dependencies Dependencies `json:"dependencies"`
}

func Status(c *fiber.Ctx) error {
	configs.InitDatabase()
	defer configs.CloseDatabase()

	var databaseVersionValue string
	if err := configs.DB.QueryRow(context.Background(), "SHOW server_version;").Scan(&databaseVersionValue); err != nil {
		log.Fatalf("Failed to get database version: %v", err)
	}

	var maxConnections string
	if err := configs.DB.QueryRow(context.Background(), "SHOW max_connections;").Scan(&maxConnections); err != nil {
		log.Fatalf("Failed to get database max connections: %v", err)
	}

	var usedConnections int
	const query = "SELECT COUNT(*)::int FROM pg_stat_activity WHERE datname = $1;"
	err := configs.DB.QueryRow(context.Background(), query, os.Getenv("POSTGRES_DB")).Scan(&usedConnections)

	if err != nil {
		log.Fatalf("Failed to get database used connections: %v", err)
	}

	maxConnectionsInt, err := strconv.Atoi(maxConnections)
	if err != nil {
		log.Fatalf("Failed to convert max connections to int: %v", err)
	}

	response := Response{
		UpdatedAt: time.Now().Format(time.RFC3339),
		Dependencies: Dependencies{
			Database: Database{
				Version:         databaseVersionValue,
				MaxConnections:  maxConnectionsInt,
				UsedConnections: usedConnections,
			},
		},
	}

	return c.JSON(response)
}
