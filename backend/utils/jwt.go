package utils

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	UserID string `json:"userID"`
	Email  string `json:"email"`
	Role   string `json:"role"` // Added role for role-based authorization
	jwt.RegisteredClaims
}

// Added RefreshClaims for refresh tokens
type RefreshClaims struct {
	UserID string `json:"userID"`
	jwt.RegisteredClaims
}

func GenerateToken(userID, email, role string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role, // Added role to claims
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), // Short-lived access token
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// Added GenerateRefreshToken for refresh functionality
func GenerateRefreshToken(userID string) (string, error) {
	claims := &RefreshClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // Long-lived refresh token
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// Added ValidateRefreshToken for refresh endpoint
func ValidateRefreshToken(tokenStr string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&RefreshClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		},
	)
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*RefreshClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

func ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		},
	)
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok {
		return claims, nil
	}
	return nil, err
}

// CHANGE: ExtractUserIDFromToken extracts user ID from JWT token in request header
// Used in controllers to get the current user's ID from the Authorization header
// Returns: userID (string) from token claims
func ExtractUserIDFromToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", context.DeadlineExceeded
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", context.DeadlineExceeded
	}

	claims, err := ValidateToken(parts[1])
	if err != nil {
		return "", err
	}

	return claims.UserID, nil
}

// CHANGE: CreateContext creates a context with timeout for MongoDB operations
// Used in controllers to ensure database queries don't hang indefinitely
// Returns: context with 10-second timeout and cancel function
func CreateContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
