package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"

	"github.com/pavittarx/lang-portal/internal/config"
	"github.com/pavittarx/lang-portal/pkg/handlers"
	"github.com/pavittarx/lang-portal/pkg/repository"
	"github.com/pavittarx/lang-portal/pkg/services"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	db, err := sql.Open("sqlite3", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Ping the database to verify connection
	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Initialize repositories
	wordRepo := repository.NewWordRepository(db)
	groupRepo := repository.NewSQLiteGroupRepository(db)

	// Initialize services
	wordService := services.NewWordService(wordRepo, groupRepo)

	// Initialize handlers
	wordHandler := handlers.NewWordHandler(wordService)

	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	e.GET("/random", wordHandler.GetRandomWord)
	e.GET("/word/:id", wordHandler.GetWordDetails)
	e.GET("/group/:group", wordHandler.ListWordsByGroup)

	// Start server
	serverAddress := fmt.Sprintf(":%d", cfg.ServerPort)
	log.Printf("Starting server on %s", serverAddress)
	
	if err := e.Start(serverAddress); err != nil {
		log.Fatalf("Server shutdown: %v", err)
	}
}
