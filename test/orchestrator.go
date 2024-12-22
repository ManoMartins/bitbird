package test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ManoMartins/bitbird/configs"
	"github.com/joho/godotenv"
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

func CleanDatabase() error {
	if envErr := godotenv.Load("../../.env"); envErr != nil {
		return errors.New(".env file missing")
	}

	configs.InitDatabase()

	_, err := configs.DB.Exec(context.Background(), `
		DROP SCHEMA public CASCADE;
		CREATE SCHEMA public;
	`)
	if err != nil {
		return fmt.Errorf("failed to reset database: %v", err)
	}

	configs.CloseDatabase()

	return nil
}

func Setup() {
	if err := WaitForAllServices(); err != nil {
		panic(err)
	}

	if err := CleanDatabase(); err != nil {
		panic(err)
	}
}
