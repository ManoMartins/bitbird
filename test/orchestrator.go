package test

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

func WaitForAllServices() error {
	const maxRetries = 100
	const retryDelay = 100 * time.Millisecond

	for i := 0; i < maxRetries; i++ {
		err := fetchStatusPage()

		if err == nil {
			return nil
		}

		log.Printf("service not ready, retrying in %s\n", retryDelay)
		time.Sleep(retryDelay)
	}

	return errors.New("service not ready after max retries")
}

func fetchStatusPage() error {
	resp, err := http.Get("http://localhost:8080/api/v1/status")

	if err != nil {
		return fmt.Errorf("failed to perform GET request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("service not ready")
	}

	return nil
}
