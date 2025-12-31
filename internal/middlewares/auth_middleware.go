package middlewares

import (
	"book-api/internal/utils"
	"context"
	"net/http"
	"strings"
)

type contextKey string

const UserContextKey contextKey = "user"

// AuthMiddleware - middleware untuk validasi JWT
func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Ambil token dari header Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.ErrorResponse(w, http.StatusUnauthorized, "Missing authorization header")
				return
			}

			// Format: "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid authorization header format")
				return
			}

			token := parts[1]

			// Validasi token
			claims, err := utils.ValidateToken(token, jwtSecret)
			if err != nil {
				utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid or expired token")
			}

			// Simpan user info ke context
			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserFromContext - helper untuk ambil user info dari context
func GetUserFromContext(r *http.Request) *utils.JWTClaim {
	claims, ok := r.Context().Value(UserContextKey).(*utils.JWTClaim)
	if !ok {
		return nil
	}
	return claims
}