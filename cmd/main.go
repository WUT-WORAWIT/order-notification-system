package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"order-notification-system/internal/config"
	"order-notification-system/internal/middleware" // Added import for middleware
	"order-notification-system/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.Init()
	if db == nil {
		log.Fatalf("Failed to initialize database connection")
	}

	gin.SetMode(gin.ReleaseMode)
	// Initialize Gin router with Logger and Recovery middleware
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.Use(middleware.CORSMiddleware()) // Use CORSMiddleware from middleware package

	routes.SetupRouter(r, db)

	serverAddr := ":8080" // TODO: Make server address configurable
	server := &http.Server{
		Addr:    serverAddr,
		Handler: r,
		// Consider adding ReadTimeout, WriteTimeout, IdleTimeout for better production hardening
	}

	// Start server in a goroutine so that it doesn't block.
	go func() {
		log.Printf("Starting server on %s", serverAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %s\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
