package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"rewardpage/middleware"
	"rewardpage/service"
	"rewardpage/utils"
	"time"
)

// GetStreak retrieves the current streak data for the logged-in user
// Frontend: GET /api/streak (authenticated)
// Response: { id, userId, mon, tue, wed, thu, fri, sat, sun, lastCheckIn, updatedAt }
// Called by TaskDashboard.jsx on mount to populate DailyStreak component
func GetStreak(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract userID from JWT claims (set by AuthMiddleware)
	claims := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
	userID := claims.UserID

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	streak, err := service.StreakServiceInstance.GetStreakByUserID(ctx, userID)
	if err != nil {
		http.Error(w, "Error fetching streak", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(streak)
}

// UpdateStreak performs a daily check-in and updates the streak
// Frontend: POST /api/streak/update (authenticated)
// Request body: {} (empty - today's day is determined server-side)
// Response: { id, userId, mon, tue, wed, thu, fri, sat, sun, lastCheckIn, updatedAt }
// Called by DailyStreak.jsx when user clicks "Check in Today" button
// Points: Adds 5 points to user for daily check-in
// Logic:
// 1. Determine current day of week (Mon-Sun)
// 2. Set that day to true in the streak
// 3. Update lastCheckIn timestamp
// 4. Award 5 points
// 5. Return updated streak object for DailyStreak component to display
func UpdateStreak(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract userID from JWT claims
	claims := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
	userID := claims.UserID

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Update streak with today's check-in
	updatedStreak, err := service.StreakServiceInstance.UpdateStreak(ctx, userID)
	if err != nil {
		http.Error(w, "Error updating streak", http.StatusInternalServerError)
		return
	}

	// Add 5 points to user for daily check-in
	_ = service.LeaderboardServiceInstance.AddPointsToUser(ctx, userID, 5)

	json.NewEncoder(w).Encode(updatedStreak)
}

// GetStreakCount returns the number of consecutive days checked in
// Frontend: GET /api/streak/count (authenticated) - optional, for displaying total
// Response: { count: 5 }
// Could be used for achievements or special display on DailyStreak component
func GetStreakCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	claims := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
	userID := claims.UserID

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	count, err := service.StreakServiceInstance.GetStreakCount(ctx, userID)
	if err != nil {
		http.Error(w, "Error fetching streak count", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]int{
		"count": count,
	})
}
