package main

import (
	"log"

	"github.com/pavittarx/lang-portal/internal/config"
	"github.com/pavittarx/lang-portal/internal/server"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create server
	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Start server
	if err := srv.Start(); err != nil {
		log.Fatalf("Server shutdown: %v", err)
	}

	// Graceful shutdown
	defer func() {
		if err := srv.Shutdown(); err != nil {
			log.Printf("Error during server shutdown: %v", err)
		}
	}()
}
