package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"rewardpage/middleware"
	"rewardpage/model"
	"rewardpage/service"
	"rewardpage/utils"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var input model.UserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var user model.User
	err := service.UserServiceInstance.FindUserByEmail(ctx, input.Email, &user)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	// Updated to include role and generate refresh token for product-ready auth
	accessToken, err := utils.GenerateToken(user.ID.Hex(), user.Email, user.Role)
	if err != nil {
		http.Error(w, "Error generating access token", http.StatusInternalServerError)
		return
	}

	// Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(user.ID.Hex())
	if err != nil {
		http.Error(w, "Error generating refresh token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})

}

// Added Refresh endpoint for token refresh functionality
func Refresh(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Validate refresh token
	claims, err := utils.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Get user by ID
	var user model.User
	err = service.UserServiceInstance.FindUserByID(ctx, claims.UserID, &user)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Generate new access token
	accessToken, err := utils.GenerateToken(claims.UserID, user.Email, user.Role)
	if err != nil {
		http.Error(w, "Error generating access token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"access_token": accessToken,
	})
}

// Added Logout endpoint for token blacklisting
func Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Authorization header missing or invalid", http.StatusBadRequest)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Blacklist the token
	err := service.BlacklistServiceInstance.BlacklistToken(ctx, token)
	if err != nil {
		http.Error(w, "Error logging out", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logged out successfully",
	})
}

// Added Me endpoint to get logged-in user details
func Me(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	claims := r.Context().Value(middleware.UserContextKey).(*utils.Claims)

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var user model.User
	err := service.UserServiceInstance.FindUserByID(ctx, claims.UserID, &user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	userOutput := model.UserOutput{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}

	json.NewEncoder(w).Encode(userOutput)
}
