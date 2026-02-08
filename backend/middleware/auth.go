package middleware

import (
	"context"
	"net/http"
	"rewardpage/service"
	"rewardpage/utils"

	// "rewardpage/utils"
	"strings"
)

type contextKey string

const UserContextKey contextKey = "user"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		// Validate the JWT token
		// Added to properly validate the token and extract claims
		claims, err := utils.ValidateToken(parts[1])
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Check if token is blacklisted
		// Added for logout functionality
		ctx := context.Background()
		isBlacklisted, err := service.BlacklistServiceInstance.IsTokenBlacklisted(ctx, parts[1])
		if err != nil {
			http.Error(w, "Error checking token", http.StatusInternalServerError)
			return
		}
		if isBlacklisted {
			http.Error(w, "Token is blacklisted", http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Added RequireRole middleware for role-based authorization
func RequireRole(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := r.Context().Value(UserContextKey).(*utils.Claims)
			if claims.Role != requiredRole {
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
