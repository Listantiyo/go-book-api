package routes

import (
	"net/http"

	"book-api/internal/handlers"
	"book-api/internal/middlewares"

	_ "book-api/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRoutes(authHandler *handlers.AuthHandler, bookHandler *handlers.BookHandler, borrowHandler *handlers.BorrowHandler, jwtSecret string) *chi.Mux {
	r := chi.NewRouter()

	//Middleware global
	r.Use(middleware.Logger)		// Log semua request
	r.Use(middleware.Recoverer)		// Recover dari semua panic
	r.Use(middleware.RequestID)		// Add request ID untuk memberikan id pada log.
	r.Use(middleware.RealIP)		// Get real IP
	r.Use(middleware.AllowContentType("application/json","application/json; charset=utf-8")) // Only accept JSON

	// Healt check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request){
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	// API Routes
	r.Route("/api/v1", func(r chi.Router) {
		// Auth routes (public)
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)

		// Book routes (akan ditambahkan auth middleware nantinya)
		r.Route("/books", func(r chi.Router) {
			
			// Public endpoints - siapa aja bisa akses
			r.Get("/", bookHandler.GetAllBooks)			// GET /api/v1/books
			r.Get("/{id}", bookHandler.GetBookByID)		// GET /api/v1/books/1
			
			// Protected endpoints - harus login dulu
			r.Group(func(r chi.Router){
				r.Use(middlewares.AuthMiddleware(jwtSecret))
				r.Post("/", bookHandler.CreateBook)			// POST /api/v1/books
				r.Put("/{id}", bookHandler.UpdateBook)		// PUT /api/v1/books/1
				r.Delete("/{id}", bookHandler.DeleteBook)	// DELETE /api/v1/books/1
			})
		})

		r.Route("/borrow", func(r chi.Router) {
			r.Use(middlewares.AuthMiddleware(jwtSecret))
			r.Post("/", borrowHandler.BorrowBook)
			r.Post("/return", borrowHandler.ReturnBook)
			r.Get("/me", borrowHandler.GetMyBorrows)
			r.Get("/{id}", borrowHandler.GetBorrowByID)
		})
	})

	return r
}