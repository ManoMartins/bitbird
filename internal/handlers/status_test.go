package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/ManoMartins/bitbird/test"
)

func TestStatusEndpoint(t *testing.T) {
	test.Setup()

	// Perform a GET request to the actual server
	resp, err := http.Get("http://localhost:8080/api/v1/status")
	if err != nil {
		t.Fatalf("Failed to perform GET request: %v", err)
	}
	defer resp.Body.Close()

	// Assert the response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Define the expected structure

	var actualResponse Response
	err = json.Unmarshal(body, &actualResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	if actualResponse.UpdatedAt == "" {
		t.Errorf("Expected non-empty 'updated_at' field")
	}

	if actualResponse.Dependencies.Database.Version != "16.0" {
		t.Errorf("Expected database 'version' 16.0, got %s", actualResponse.Dependencies.Database.Version)
	}

	if actualResponse.Dependencies.Database.MaxConnections != 100 {
		t.Errorf("Expected database 'max_connections' 100, got %d", actualResponse.Dependencies.Database.MaxConnections)
	}

	if actualResponse.Dependencies.Database.UsedConnections != 1 {
		t.Errorf("Expected database 'used_connections' 1, got %d", actualResponse.Dependencies.Database.UsedConnections)
	}

	parsedTime, err := time.Parse(time.RFC3339, actualResponse.UpdatedAt)
	if err != nil {
		t.Fatalf("Failed to parse time: %v", err)
	}

	if actualResponse.UpdatedAt != parsedTime.Format(time.RFC3339) {
		t.Errorf("Expected 'updated_at' %s, got %s", actualResponse.UpdatedAt, parsedTime.Format(time.RFC3339))
	}

}
