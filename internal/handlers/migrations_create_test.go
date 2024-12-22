package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ManoMartins/bitbird/test"
)

func TestMigrationCreateEndpoint(t *testing.T) {
	test.Setup()

	resp, err := http.Post("http://localhost:8080/api/v1/migrations", "application/json", nil)

	if err != nil {
		t.Fatalf("Failed to perform POST request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	var actualResponse map[string]string

	err = json.NewDecoder(resp.Body).Decode(&actualResponse)

	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	if actualResponse["message"] != "Migrations applied successfully" {
		t.Errorf("Expected message 'Migrations applied successfully', got %s", actualResponse["message"])
	}

	resp, err = http.Post("http://localhost:8080/api/v1/migrations", "application/json", nil)

	if err != nil {
		t.Fatalf("Failed to perform POST request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&actualResponse)

	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	if actualResponse["message"] != "No pending migrations" {
		t.Errorf("Expected message 'No pending migrations', got %s", actualResponse["message"])
	}

}
