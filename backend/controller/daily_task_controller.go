package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"rewardpage/service"
	"rewardpage/utils"
)

// CHANGE: GetDailyTasks handles GET /api/tasks/daily endpoint
//
//	Returns: {
//	  tasks: [{ id, number, completed, completedAt }],
//	  completedCount: number,
//	  lastCompletedAt: timestamp,
//	  nextResetAt: timestamp,
//	  cooldownUntil: timestamp
//	}
//
// Frontend uses this to:
// 1. Load 5 task buttons at page load
// 2. Show completion status (completed, disabled, available)
// 3. Show cooldown timer if active
func GetDailyTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// CHANGE: Extract user ID from JWT token (set by auth middleware)
	userID, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	ctx, cancel := utils.CreateContext()
	defer cancel()

	// CHANGE: Check if past midnight and auto-reset tasks
	// Lazy-loaded reset: runs on first request after day change
	// No cron job needed - resets when user opens app next day
	if err := service.DailyTaskServiceInstance.CheckAndResetDaily(ctx, userID); err != nil {
		log.Printf("Warning: CheckAndResetDaily failed: %v", err)
		// Don't return error - continue with existing tasks
	}

	// CHANGE: Get or create today's 5 tasks (auto-creates if first time today)
	tasks, err := service.DailyTaskServiceInstance.GetOrCreateDailyTasks(ctx, userID)
	if err != nil {
		http.Error(w, `{"error":"failed to load tasks"}`, http.StatusInternalServerError)
		return
	}

	// CHANGE: Get progress to return cooldown and completion count
	progress, err := service.DailyTaskServiceInstance.GetOrCreateProgress(ctx, userID)
	if err != nil {
		http.Error(w, `{"error":"failed to load progress"}`, http.StatusInternalServerError)
		return
	}

	// CHANGE: Return tasks with progress metadata to frontend
	response := map[string]interface{}{
		"tasks":           tasks,
		"completedCount":  progress.CompletedCount,
		"lastCompletedAt": progress.LastCompletedAt,
		"nextResetAt":     progress.NextResetAt,
		"cooldownUntil":   progress.LastCooldownEnd,
	}

	json.NewEncoder(w).Encode(response)
}

// CHANGE: CompleteTaskDaily handles POST /api/tasks/complete endpoint
// Request body: { taskId: string }
//
//	Response: {
//	  success: boolean,
//	  task: { id, number, completed, completedAt },
//	  completedCount: number,
//	  nextResetAt: timestamp,
//	  cooldownUntil: timestamp
//	}
//
// Validation performed server-side:
// - User not in cooldown (5 minutes since last task)
// - User hasn't completed 5 tasks today
// - Daily reset check (if past midnight)
func CompleteTaskDaily(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// CHANGE: Extract user ID from JWT token
	userID, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// CHANGE: Parse request body for taskId
	var input map[string]string
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	taskID := input["taskId"]
	if taskID == "" {
		http.Error(w, `{"error":"task_id required"}`, http.StatusBadRequest)
		return
	}

	ctx, cancel := utils.CreateContext()
	defer cancel()

	// CHANGE: Complete task with backend cooldown validation
	// This function performs all server-side checks before allowing completion
	result, err := service.DailyTaskServiceInstance.CompleteTask(ctx, userID, taskID)

	if err != nil {
		log.Printf("ERROR - CompleteTask failed: userID=%s, taskID=%s, error=%v, type=%T", userID, taskID, err, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
			"debug": err.Error(), // Return full error to frontend for debugging
		})
		return
	}

	// CHANGE: Award points to user for task completion (20 points per task)
	// This updates the leaderboard in real-time
	const POINTS_PER_TASK = 20
	if err := service.LeaderboardServiceInstance.AddPointsToUser(ctx, userID, POINTS_PER_TASK); err != nil {
		log.Printf("Warning: Failed to add points to user %s: %v", userID, err)
		// Non-blocking error - task still completed, just points not updated
	}

	// CHANGE: Success response must include success flag
	response := map[string]interface{}{
		"success":        true,
		"task":           result["task"],
		"completedCount": result["completed_count"],
		"nextResetAt":    result["next_reset_at"],
		"cooldownUntil":  result["cooldown_until"],
		"pointsAwarded":  POINTS_PER_TASK,
	}
	log.Printf("Task completed successfully: userID=%s, taskID=%s, pointsAwarded=%d", userID, taskID, POINTS_PER_TASK)
	json.NewEncoder(w).Encode(response)
}

// CHANGE: CheckCooldown handles GET /api/tasks/cooldown endpoint
//
//	Returns: {
//	  isCooldownActive: boolean,
//	  remainingSeconds: number,
//	  lastCompletedAt: timestamp,
//	  completedCount: number
//	}
//
// Frontend uses this to:
// 1. Determine if user can click next task
// 2. Calculate remaining countdown timer
// 3. Show disabled/enabled button state
func CheckCooldown(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := utils.ExtractUserIDFromToken(r)
	if err != nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	ctx, cancel := utils.CreateContext()
	defer cancel()

	// CHANGE: Get current progress to check cooldown status
	progress, err := service.DailyTaskServiceInstance.GetOrCreateProgress(ctx, userID)
	if err != nil {
		http.Error(w, `{"error":"failed to check cooldown"}`, http.StatusInternalServerError)
		return
	}

	// CHANGE: Calculate if user is currently in cooldown period
	isCooldownActive := false
	remainingSeconds := 0

	if progress.LastCompletedAt != nil {
		elapsed := time.Since(*progress.LastCompletedAt)
		if elapsed < 5*time.Minute {
			isCooldownActive = true
			remainingSeconds = int(((5 * time.Minute) - elapsed).Seconds())
		}
	}

	response := map[string]interface{}{
		"isCooldownActive": isCooldownActive,
		"remainingSeconds": remainingSeconds,
		"lastCompletedAt":  progress.LastCompletedAt,
		"completedCount":   progress.CompletedCount,
	}

	json.NewEncoder(w).Encode(response)
}

// CHANGE: Helper function to get map keys for logging
func getMapKeys(m map[string]interface{}) []string {
	if m == nil {
		return nil
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
