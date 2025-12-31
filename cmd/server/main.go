package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"book-api/internal/config"
	"book-api/internal/database"
	"book-api/internal/handlers"
	"book-api/internal/models"
	"book-api/internal/repository"
	"book-api/internal/routes"
	"book-api/internal/services"
)

// @title Book API
// @version 1.0
// @description A product-ready REST API for managing books and book borrowing system
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@bookapi.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinition.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by space and JWT token.

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Connect to database
	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate models
	if err := db.AutoMigrate(&models.User{}, &models.Book{}, &models.Borrow{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("✅ Database migration completed")

	// Initialize repository
	userRepo 	:= repository.NewUserRepository(db)
	bookRepo 	:= repository.NewBookRepository(db)
	borrowRepo 	:= repository.NewBorrowRepository(db)

	// Initialize transaction manager
	txManager	:= database.NewTransactionManager(db)

	// Initialize services
	authService 	:= services.NewAuthService(userRepo)
	bookService 	:= services.NewBookService(bookRepo)
	borrowService 	:= services.NewBorrowService(borrowRepo, bookRepo, txManager)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, cfg.JWTSecret)
	bookHandler := handlers.NewBookHandler(bookService)
	borrowHandler := handlers.NewBorrowHandler(borrowService)

	// Setup routes
	router := routes.SetupRoutes(authHandler, bookHandler, borrowHandler, cfg.JWTSecret)

	// Create HTTP server
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	log.Println("🔧 Server configured, starting goroutine...")

	// Start server in goroutine
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("🚀 Server running on http://localhost%s", addr)
		log.Printf("📚 API Documentation: http://localhost%s/api/v1", addr)
		log.Println("👉 Press Ctrl+C to stop")
		
		serverErrors <- srv.ListenAndServe()
	}()

	log.Println("⏸️  Main goroutine waiting for signal...")

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	
	log.Println("✅ Signal handler registered")
	
	// Wait for either error from server or interrupt signal
	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server error: %v", err)
		}
	case sig := <-quit:
		log.Printf("\n🛑 Signal received: %v", sig)
		log.Println("⏳ Shutting down server gracefully...")
	}

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("❌ Error during shutdown: %v\n", err)
		os.Exit(1)
	}

	log.Println("✅ Server stopped gracefully")
	log.Println("👋 Goodbye!")
}
