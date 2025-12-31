package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginatedResponse struct {
	Data interface{} `json:"data"`
	Page int `json:"page"`
	PageSize int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

// WriteJSON - helper untuk kirim response JSON
func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// SuccessResponse - respon sukses
func SuccessResponse(w http.ResponseWriter, status int, message string, data interface{}) {
	WriteJSON(w, status, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse - respon error
func ErrorResponse(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, Response{
		Success: false,
		Error: 	 message,
	})
}