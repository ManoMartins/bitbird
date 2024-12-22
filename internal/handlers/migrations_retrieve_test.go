package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/ManoMartins/bitbird/configs"
	"github.com/ManoMartins/bitbird/test"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	if envErr := godotenv.Load("../../.env"); envErr != nil {
		log.Fatal(".env file missing")
	}

	if err := test.WaitForAllServices(); err != nil {
		panic(err)
	}

	configs.InitDatabase()

	_, err := configs.DB.Exec(context.Background(), `
		DROP SCHEMA public CASCADE;
		CREATE SCHEMA public;
	`)
	if err != nil {
		log.Fatalf("Failed to reset database: %v", err)
	}

	configs.CloseDatabase()

	code := m.Run()
	os.Exit(code)
}

func TestMigrationRetrieveEndpoint(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/api/v1/migrations")

	if err != nil {
		t.Fatalf("Failed to perform GET request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	var actualResponse MigrationsRetrieveOutput
	err = json.NewDecoder(resp.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	if len(actualResponse.PendingMigrations) == 0 {
		t.Error("Expected at least one pending migration")
	}
}
