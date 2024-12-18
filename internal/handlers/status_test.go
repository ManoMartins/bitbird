package handlers

import (
	"net/http"
	"testing"
)

func TestStatusEndpoint(t *testing.T) {
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
}
