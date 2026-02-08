package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"rewardpage/middleware"
	"rewardpage/model"
	"rewardpage/service"
	"rewardpage/utils"
	"time"
)

// GetTasks retrieves all tasks for the logged-in user
// Frontend: GET /api/tasks (authenticated)
// Response: array of { id, title, completed, userId, createdAt, updatedAt }
// Called by TaskDashboard.jsx on mount to populate NormalTasks component
func GetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract userID from JWT claims (set by AuthMiddleware)
	claims := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
	userID := claims.UserID

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	tasks, err := service.TaskServiceInstance.GetTasksByUserID(ctx, userID)
	if err != nil {
		http.Error(w, "Error fetching tasks", http.StatusInternalServerError)
		return
	}

	// Return empty array if no tasks (instead of null)
	if tasks == nil {
		tasks = []model.Task{}
	}

	json.NewEncoder(w).Encode(tasks)
}

// CompleteTask marks a task as completed for the logged-in user
// Frontend: POST /api/tasks/complete (authenticated)
// Request body: { taskId, box }
// Response: { message: "Task completed" }
// Called by NormalTasks.jsx when user checks a task checkbox
// Points: Adds 10 points to user when task is completed
func CompleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract userID from JWT claims
	claims := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
	userID := claims.UserID

	// Parse request body
	var req struct {
		TaskID string `json:"taskId"`
		Box    string `json:"box"` // "left", "right", "center"
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.TaskID == "" {
		http.Error(w, "taskId is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Mark task as completed
	err := service.TaskServiceInstance.CompleteTask(ctx, userID, req.TaskID)
	if err != nil {
		http.Error(w, "Error completing task", http.StatusInternalServerError)
		return
	}

	// Add 10 points to user for completing task
	_ = service.LeaderboardServiceInstance.AddPointsToUser(ctx, userID, 10)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Task completed successfully",
	})
}

// CreateTask allows creating a new task (admin or user-initiated)
// Frontend: POST /api/tasks (authenticated) - future use
// Request body: { title, box }
// Response: { id, title, completed, userId, createdAt }
func CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract userID from JWT claims
	claims := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
	userID := claims.UserID

	var task model.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	// Set userId and defaults
	var userObjID interface{}
	_, _ = userObjID, userID // Set userID from claims in real scenario

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	err := service.TaskServiceInstance.CreateTask(ctx, &task)
	if err != nil {
		http.Error(w, "Error creating task", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(task)
}
