package main

import (
	"log"
	"os"

	"github.com/anomalyco/auta/internal/metadata"
)

func main() {
	// Load configuration
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://localhost/auta?sslmode=disable"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// Initialize metadata service
	service, err := metadata.NewService(dbURL)
	if err != nil {
		log.Fatalf("Failed to initialize metadata service: %v", err)
	}
	defer service.Close()

	// Start the service
	log.Printf("Starting metadata service on port %s", port)
	if err := service.Start(port); err != nil {
		log.Fatalf("Failed to start metadata service: %v", err)
	}
}
