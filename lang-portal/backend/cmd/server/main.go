package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	// API Health Check Endpoint
	e.GET("/api", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "healthy",
			"message": "Language Portal Backend is running successfully!",
		})
	})

	// Routes (to be added later)
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "healthy",
		})
	})

	// Server configuration
	port := "3000"
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      e,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start server
	go func() {
		sugar.Infof("Starting server on port %s", port)
		if err := e.StartServer(server); err != nil && err != http.ErrServerClosed {
			sugar.Fatalf("Server shutdown: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	sugar.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		sugar.Fatalf("Server shutdown error: %v", err)
	}
	sugar.Info("Server stopped")
}
