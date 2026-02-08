package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"rewardpage/middleware"
	"rewardpage/model"
	"rewardpage/service"
	"rewardpage/utils"
	"strconv"
	"time"
)

// GetLeaderboard retrieves top-ranked users sorted by points (descending)
// Backend: GET /api/leaderboard (authenticated)
// Query params: limit (optional, default 10, max 100)
// Response: array of LeaderboardUser { id, username, email, points, rank }
// Features:
// - Sorts by points descending (highest points first)
// - Assigns sequential rank numbers (1, 2, 3, ...)
// - Limits results to prevent large data transfers
func GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract userID from JWT claims (for future personalization)
	claims := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
	_ = claims

	// Parse optional limit query parameter (default 10, max 100)
	limitStr := r.URL.Query().Get("limit")
	limit := int64(10)
	if limitStr != "" {
		if parsed, err := strconv.ParseInt(limitStr, 10, 64); err == nil && parsed > 0 {
			if parsed > 100 {
				parsed = 100 // Cap at 100 to prevent large transfers
			}
			limit = parsed
		}
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	leaderboard, err := service.LeaderboardServiceInstance.GetLeaderboard(ctx, limit)
	if err != nil {
		http.Error(w, `{"error":"Error fetching leaderboard"}`, http.StatusInternalServerError)
		return
	}

	// Return empty array if no users
	if leaderboard == nil {
		leaderboard = []model.LeaderboardUser{}
	}

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	json.NewEncoder(w).Encode(leaderboard)
}

// GetUserRank retrieves a specific user's rank and points
// Frontend: GET /api/leaderboard/me (authenticated)
// Response: { id, username, email, points, rank }
// Could be used in future to display "Your Position" on the leaderboard UI
// Example response:
// { id: "...", username: "charlie", email: "charlie@example.com", points: 350, rank: 5 }
func GetUserRank(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract userID from JWT claims
	claims := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
	userID := claims.UserID

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	userRank, err := service.LeaderboardServiceInstance.GetUserRank(ctx, userID)
	if err != nil {
		http.Error(w, "Error fetching user rank", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(userRank)
}
