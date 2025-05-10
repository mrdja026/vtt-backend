package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"

	"dnd-combat/api/v1"
	"dnd-combat/config"
	"dnd-combat/pkg/database"
	"dnd-combat/pkg/dnd5e"
	"dnd-combat/pkg/middleware"
	"dnd-combat/pkg/websocket"
)

func main() {
	// Initialize configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Setup database connection
	db, err := database.New(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize HTTP client for SRD API
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create SRD client
	srdClient := dnd5e.NewSRDClient(httpClient, dnd5e.NewInMemoryCache())

	// Create websocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Setup Gin router
	router := gin.Default()

	// Apply middleware
	router.Use(middleware.Logger())
	router.Use(middleware.Cors())

	// Setup routes
	v1.SetupRoutes(router, db, srdClient, hub, cfg)

	// Create HTTP server
	srv := &http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	// Give the server a grace period to finish handling requests
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
