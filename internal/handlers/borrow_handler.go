package handlers

import (
	"book-api/internal/middlewares"
	"book-api/internal/services"
	"book-api/internal/utils"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type BorrowHandler struct {
	borrowService services.BorrowService
}

func NewBorrowHandler(borrowservice services.BorrowService) *BorrowHandler {
	return & BorrowHandler{borrowService: borrowservice}
}

type BorrowBookRequest struct {
	BookID uint `json:"book_id" validate:"required,gte=0"`
}

type ReturnBookRequest struct {
	BorrowID uint `json:"borrow_id" validate:"required,gte=0"`
}

// BorrowBook godoc
// @Summary Borrow a book
// @Description Borrow a book (requires authentication, decreases stock) 
// @Tags Borrows
// @Accept json
// @Produce json
// @Security BeareAuth
// @Param request body BorrowBookRequest true "Book ID to borrow"
// @Success 201 {object} utils.Response{data=models.Borrow}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /borrows [post]
func (h *BorrowHandler) BorrowBook(w http.ResponseWriter, r *http.Request) {
	// Get user dari context
	claims := middlewares.GetUserFromContext(r)
	if claims == nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req BorrowBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Borrow book
	borrow, err := h.borrowService.BorrowBook(claims.UserID, req.BookID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusCreated, "Book borrowed successfully", borrow)
}

// ReturnBook godoc
// @Summary Return a borrowed book
// @Description Return a borrowe book (requires authentication, increases stock) 
// @Tags Borrows
// @Accept json
// @Produce json
// @Security BeareAuth
// @Param request body ReturnBookRequest true "Book ID to return"
// @Success 200 {object} utils.Response{data=models.Borrow}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /borrows/return [post]
func (h *BorrowHandler) ReturnBook(w http.ResponseWriter, r *http.Request) {
	claims := middlewares.GetUserFromContext(r)
	if claims == nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req ReturnBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Return borrowed
	borrow, err := h.borrowService.ReturnBook(req.BorrowID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Book returned successfully", borrow)
}

// GetMyBorrows godoc
// @Summary Get my borrow history
// @Description Get current user`s borrow history with pagination 
// @Tags Borrows
// @Accept json
// @Produce json
// @Security BeareAuth
// @Param request body ReturnBookRequest true "Book ID to return"
// @Success 200 {object} utils.Response{data=utils.PaginatedResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /borrows/me [get]
func (h *BorrowHandler) GetMyBorrows(w http.ResponseWriter, r *http.Request) {
	claims := middlewares.GetUserFromContext(r)
	if claims == nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Unatuhorized")
		return
	}

	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")

	page := 1
	pageSize := 10

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

	// 'total' it contain all count borrowed
	borrows, total, err := h.borrowService.GetUserBorrows(claims.UserID, page, pageSize)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Count pages for paginate (total borrowed / pageSize)  
	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		// tambah satu halaman paginate jika ada sisa
		totalPages++
	}

	response := utils.PaginatedResponse{
		Data 			: borrows,
		Page 			: page,
		PageSize 		: pageSize,
		TotalItems 		: int(total),
		TotalPages		: totalPages,
	}

	utils.SuccessResponse(w, http.StatusOK, "Borrow retrieved successfully", response)
}

// GetBorrowByID godoc
// @Summary Get a borrow history
// @Description Get current user`s a borrow history 
// @Tags Borrows
// @Accept json
// @Produce json
// @Security BeareAuth
// @Param id query int true "Borrow ID to get"
// @Success 200 {object} utils.Response{data=models.Borrow}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /borrows/me [get]
func (h *BorrowHandler) GetBorrowByID(w http.ResponseWriter, r * http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid borrow ID")
		return
	}
	
	borrow, err := h.borrowService.GetBorrowByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Borrow retrieved successfully", borrow)
}