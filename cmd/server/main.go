package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abiewardani/go-messaging-system/internal/app"
	"github.com/abiewardani/go-messaging-system/internal/config"
	"github.com/abiewardani/go-messaging-system/internal/consumer"
	"github.com/abiewardani/go-messaging-system/internal/repository"
	"github.com/abiewardani/go-messaging-system/internal/service"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
	// Load configuration
	_, err := config.LoadConfig("config/config.json")
	if err != nil {
		log.Fatalf("Could not load config: %s\n", err)
		return
	}

	// Initialize RabbitMQ connection
	amqpConn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Could not connect to RabbitMQ: %s\n", err)
		return
	}
	defer amqpConn.Close()

	tenantManager, err := consumer.NewTenantManager(amqpConn)
	if err != nil {
		log.Fatalf("Could not create tenant manager: %s\n", err)
		return
	}

	// Initialize database connection
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/dbname")
	if err != nil {
		log.Fatalf("Could not connect to database: %s\n", err)
		return
	}
	defer db.Close()

	messageRepo := repository.NewMessageRepository(db)
	messageService := service.NewMessageService(*messageRepo)
	server := app.NewServer(tenantManager, messageService)

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: server.Router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v\n", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server is shutting down...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown steps
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v\n", err)
	}

	// Close tenant manager (this will close all consumer connections)
	if err := tenantManager.Close(); err != nil {
		log.Printf("Tenant manager shutdown error: %v\n", err)
	}

	// Close database connections
	if err := db.Close(); err != nil {
		log.Printf("Database shutdown error: %v\n", err)
	}

	log.Println("Server gracefully stopped")
}
