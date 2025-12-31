package main

import (
	"fmt"
	"log"
	"net/http"

	"book-api/internal/config"
	"book-api/internal/database"
	"book-api/internal/handlers"
	"book-api/internal/models"
	"book-api/internal/repository"
	"book-api/internal/routes"
	"book-api/internal/services"
)

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

	// Strat server
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server running on http://localhost%s", addr)
	log.Printf("API Documentatio: http://localhost%s/api/v1", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
