package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ManoMartins/bitbird/test"
)

func TestMigrationRetrieveEndpoint(t *testing.T) {
	test.Setup()

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
