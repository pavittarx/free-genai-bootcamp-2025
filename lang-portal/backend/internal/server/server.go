package server

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/pavittarx/lang-portal/internal/config"
	"github.com/pavittarx/lang-portal/pkg/handlers"
	"github.com/pavittarx/lang-portal/pkg/repository"
	"github.com/pavittarx/lang-portal/pkg/services"
)

// Server represents the application server configuration
type Server struct {
	Echo *echo.Echo
	DB   *sql.DB
	Cfg  *config.Config
}

// NewServer creates and configures a new server instance
func NewServer(cfg *config.Config) (*Server, error) {
	// Initialize database connection
	db, err := sql.Open("sqlite3", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Ping the database to verify connection
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	return &Server{
		Echo: e,
		DB:   db,
		Cfg:  cfg,
	}, nil
}

// SetupRepositories initializes all repository dependencies
func (s *Server) SetupRepositories() (*repository.SQLiteWordRepository, *repository.SQLiteGroupRepository) {
	wordRepo := repository.NewWordRepository(s.DB)
	groupRepo := repository.NewSQLiteGroupRepository(s.DB)

	return wordRepo, groupRepo
}

// SetupServices initializes all service dependencies
func (s *Server) SetupServices(wordRepo *repository.SQLiteWordRepository, groupRepo *repository.SQLiteGroupRepository) (*services.WordService, *services.GroupService) {
	wordService := services.NewWordService(wordRepo, groupRepo)
	groupService := services.NewGroupService(groupRepo, wordRepo)

	return wordService, groupService
}

// SetupHandlers initializes all handler dependencies
func (s *Server) SetupHandlers(wordService *services.WordService, groupService *services.GroupService) (*handlers.WordHandler, *handlers.GroupHandler) {
	wordHandler := handlers.NewWordHandler(wordService)
	groupHandler := handlers.NewGroupHandler(groupService)

	return wordHandler, groupHandler
}

// SetupRoutes configures all API routes
func (s *Server) SetupRoutes(wordHandler *handlers.WordHandler, groupHandler *handlers.GroupHandler) {
	// Word routes
	s.Echo.GET("/random", wordHandler.GetRandomWord)
	s.Echo.GET("/word/:id", wordHandler.GetWordDetails)
	s.Echo.GET("/group/:group", wordHandler.ListWordsByGroup)

	// Group routes
	s.Echo.GET("/groups", groupHandler.ListGroups)
	s.Echo.GET("/groups/:id", groupHandler.GetGroupDetails)
	s.Echo.GET("/groups/:id/words", groupHandler.ListGroupWords)
}

// Start begins the HTTP server
func (s *Server) Start() error {
	// Repositories
	wordRepo, groupRepo := s.SetupRepositories()

	// Services
	wordService, groupService := s.SetupServices(wordRepo, groupRepo)

	// Handlers
	wordHandler, groupHandler := s.SetupHandlers(wordService, groupService)

	// Routes
	s.SetupRoutes(wordHandler, groupHandler)

	// Start server
	serverAddress := fmt.Sprintf(":%d", s.Cfg.ServerPort)
	log.Printf("Starting server on %s", serverAddress)
	
	return s.Echo.Start(serverAddress)
}

// Shutdown gracefully closes the server and database connection
func (s *Server) Shutdown() error {
	// Close database connection
	if err := s.DB.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}

	return nil
}
