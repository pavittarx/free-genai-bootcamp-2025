package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	"github.com/pavittarx/lang-portal/backend/internal/config"
	"github.com/pavittarx/lang-portal/backend/pkg/handlers"
	"github.com/pavittarx/lang-portal/backend/pkg/repository"
	"github.com/pavittarx/lang-portal/backend/pkg/routes"
	"github.com/pavittarx/lang-portal/backend/pkg/services"

	// Swagger
	_ "github.com/pavittarx/lang-portal/backend/docs" // Swagger generated docs
	echoSwagger "github.com/swaggo/echo-swagger"
)

const (
	readTimeout     = 5 * time.Second
	writeTimeout    = 10 * time.Second
	shutdownTimeout = 10 * time.Second
)

// @title           Language Portal API
// @version         1.0
// @description     A backend API for managing language learning resources
// @contact.name    Language Portal Support
// @contact.email   support@languageportal.com
// @host            localhost:3000
// @BasePath        /
// @schemes         http
func main() {
	logger := initLogger()
	sugar := logger.Sugar()

	db, err := initDB()
	if err != nil {
		sugar.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	e := echo.New()
	setupMiddleware(e)
	setupRoutes(e, db, sugar)

	// Swagger route
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	server := createServer(e)
	startServer(server, e, sugar)
	gracefulShutdown(e, sugar)
}

func initDB() (*sql.DB, error) {
	return config.InitDatabase()
}

func initLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	return logger
}

func setupMiddleware(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
}

func setupRoutes(e *echo.Echo, db *sql.DB, sugar *zap.SugaredLogger) {
	// Initialize repositories
	wordRepo := repository.NewSQLiteWordRepository(db)
	groupRepo := repository.NewSQLiteGroupRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	studyActivityRepo := repository.NewStudyActivityRepository(db)
	sessionActivityRepo := repository.NewSessionActivityRepository(db)

	// Initialize services
	wordService := services.NewWordService(wordRepo)
	groupService := services.NewGroupService(groupRepo)
	sessionService := services.NewSessionService(sessionRepo)
	studyActivityService := services.NewStudyActivityService(studyActivityRepo)
	sessionActivityService := services.NewSessionActivityService(sessionActivityRepo, sessionRepo)

	// Initialize handlers
	wordHandler := handlers.NewWordHandler(wordService, wordRepo)
	groupHandler := handlers.NewGroupHandler(groupService, groupRepo)
	sessionHandler := handlers.NewSessionHandler(sessionService)
	studyActivityHandler := handlers.NewStudyActivityHandler(studyActivityService)
	sessionActivityHandler := handlers.NewSessionActivityHandler(sessionActivityService)

	// Register routes
	routes.RegisterRoutes(e,
		wordHandler,
		groupHandler,
		studyActivityHandler,
		sessionHandler,
		sessionActivityHandler)

	sugar.Info("Routes initialized successfully")
}

func createServer(e *echo.Echo) *http.Server {
	return &http.Server{
		Addr:         ":" + getPort(),
		Handler:      e,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	return port
}

func startServer(server *http.Server, e *echo.Echo, sugar *zap.SugaredLogger) {
	go func() {
		sugar.Infof("Starting server on port %s", server.Addr)
		if err := e.StartServer(server); err != nil && err != http.ErrServerClosed {
			sugar.Fatalf("Server startup error: %v", err)
		}
	}()
}

func gracefulShutdown(e *echo.Echo, sugar *zap.SugaredLogger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	sugar.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		sugar.Errorf("Server shutdown error: %v", err)
	}

	sugar.Info("Server gracefully stopped")
}
