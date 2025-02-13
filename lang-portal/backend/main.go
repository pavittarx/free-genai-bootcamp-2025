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
)

const (
	serverPort     = "3000"
	readTimeout    = 5 * time.Second
	writeTimeout   = 10 * time.Second
	shutdownTimeout = 10 * time.Second
)

func main() {
	logger := initLogger()
	sugar := logger.Sugar()
	defer logger.Sync()

	db := initDatabase(sugar)
	defer db.Close()

	e := initServer()
	setupMiddleware(e)
	setupRoutes(e, db, sugar)

	server := createServer(e)
	startServer(server, e, sugar)
	gracefulShutdown(e, sugar)
}

func initLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	return logger
}

func initDatabase(sugar *zap.SugaredLogger) *sql.DB {
	db, err := config.InitDatabase()
	if err != nil {
		sugar.Fatalf("Failed to initialize database: %v", err)
	}
	sugar.Info("Database initialized successfully")
	return db
}

func initServer() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	return e
}

func setupMiddleware(e *echo.Echo) {
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
}

func setupRoutes(e *echo.Echo, db *sql.DB, sugar *zap.SugaredLogger) {
	// Initialize repositories
	wordRepo := repository.NewSQLiteWordRepository(db)
	groupRepo := repository.NewSQLiteGroupRepository(db)

	// Initialize services
	wordService := services.NewWordService(wordRepo)
	groupService := services.NewGroupService(groupRepo)

	// Initialize handlers
	wordHandler := handlers.NewWordHandler(wordService, wordRepo)
	groupHandler := handlers.NewGroupHandler(groupService, groupRepo)

	// Register all routes
	routes.RegisterRoutes(e, wordHandler, groupHandler)

	sugar.Info("Routes initialized successfully")
}

func createServer(e *echo.Echo) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%s", serverPort),
		Handler:      e,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
}

func startServer(server *http.Server, e *echo.Echo, sugar *zap.SugaredLogger) {
	go func() {
		sugar.Infof("Starting server on port %s", serverPort)
		if err := e.StartServer(server); err != nil && err != http.ErrServerClosed {
			sugar.Fatalf("Server startup failed: %v", err)
		}
	}()
}

func gracefulShutdown(e *echo.Echo, sugar *zap.SugaredLogger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	sugar.Info("Initiating graceful shutdown...")
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		sugar.Fatalf("Server shutdown error: %v", err)
	}
	sugar.Info("Server stopped gracefully")
}
