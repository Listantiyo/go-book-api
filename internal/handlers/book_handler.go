package handlers

import (
	"book-api/internal/services"
	"book-api/internal/utils"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type BookHandler struct {
	bookService services.BookService
}

func NewBookHandler(bookService services.BookService) *BookHandler {
	return &BookHandler{bookService: bookService}
}

type CreateBookRequest struct {
	Title string `json:"title" validate:"required,min=1,max=200"`
	Author string `json:"author" validate:"required,min=1,max=100"`
	ISBN string `json:"isbn" validate:"required,min=10,max=13"`
	Description string `json:"description" validate:"max=1000"`
	Stock int `json:"stock" validate:"gte=0"`
}

type UpdateBookRequest struct {
	Title string `json:"title" validate:"required,min=1,max=200"`
	Author string `json:"author" validate:"required,min=1,max=100"`
	ISBN string `json:"isbn" validate:"required,min=10,max=13"`
	Description string `json:"description" validate:"max=1000"`
	Stock int `json:"stock" validate:"gte=0"`
}

// CreateBook godoc
// @Summary Create a new book
// @Description Create a new book (requires authentication)
// @Tags Books
// @Accept json
// @Produce json
// @Security BearerAuth 
// @Param request body CreateBookRequest true "Book details"
// @Success 201 {object} utils.Response{data=models.Book}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /books [post]
func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var req CreateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validasi dengan validator
	if err := utils.ValidateStruct(req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	book, err := h.bookService.CreateBook(req.Title, req.Author, req.ISBN, req.Description, req.Stock)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(w, http.StatusCreated, "Book created successfully", book)

}

// GetAllBooks godoc
// @Summary Get all books
// @Description Get list of all books with pagination
// @Tags Books
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1) 
// @Param page_size query int false "Page size" default(10) 
// @Success 200 {object} utils.Response{data=utils.PaginatedResponse}
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /books [get]
func (h *BookHandler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	// Parse query parameter untuk pagination
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")

	page := 1 // default
	pageSize := 10 // default

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
		}
	}

	books, total, err := h.bookService.GetAllBooks(page, pageSize)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	//Hitung total page
	if pageSize == 0 {
		pageSize = 10
	}
	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	response := utils.PaginatedResponse{
		Data: books,
		Page: page,
		PageSize: pageSize,
		TotalItems: int(total),
		TotalPages: totalPages,
	}

	utils.SuccessResponse(w, http.StatusOK, "Books retrieved successfully", response)
}

// GetBookByID godoc
// @Summary Get book by ID
// @Description Get detailed information about a specific book 
// @Tags Books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} utils.Response{data=models.Book}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /books/{id} [get]
func (h *BookHandler) GetBookByID(w http.ResponseWriter, r *http.Request) {
	// Ambil url param id - getbook/:id
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	book, err := h.bookService.GetBookByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Book rretrieved successfully", book)
}

// UpdateBook godoc
// @Summary Update a book 
// @Description Update book information (requires authentication) 
// @Tags Books
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Book ID"
// @Param request body UpdateBookRequest true "Update book details"
// @Success 200 {object} utils.Response{data=models.Book}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /books/{id} [put]
func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	var req UpdateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validasi dengan validator
	if err := utils.ValidateStruct(req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	book, err := h.bookService.UpdateBook(uint(id), req.Title, req.Author, req.ISBN, req.Description, req.Stock)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Book updated successfully", book)
}

// DeleteBook godoc
// @Summary Delete a book
// @Description Delete a book (requires authentication) 
// @Tags Books
// @Accept json
// @Produce json
// @Security BeareAuth
// @Param id path int true "Book ID"
// @Success 200 {object} utils.Response{data=models.Book}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /books/{id} [delete]
func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	if err := h.bookService.DeleteBook(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Book deleted successfully", nil)
}