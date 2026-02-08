package router

import (
	"rewardpage/controller"
	"rewardpage/middleware"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	// ========== PUBLIC ENDPOINTS (NO AUTH REQUIRED) ==========
	router.HandleFunc("/api/users/register", controller.Create1user).Methods("POST")
	router.HandleFunc("/api/users/login", controller.Login).Methods("POST")
	router.HandleFunc("/api/users/refresh", controller.Refresh).Methods("POST")

	// ========== SECURED ENDPOINTS (AUTH REQUIRED) ==========
	secured := router.PathPrefix("/api").Subrouter()
	secured.Use(middleware.AuthMiddleware)

	// User endpoints
	secured.HandleFunc("/users/update", controller.Update1user).Methods("PUT")
	secured.HandleFunc("/users/profile", controller.Get1user).Methods("GET")
	secured.HandleFunc("/users/delete", controller.Delete1user).Methods("DELETE")
	secured.HandleFunc("/users/getall", controller.GetAlluser).Methods("GET")
	secured.HandleFunc("/users/deleteAll", controller.DeleteAlluser).Methods("DELETE")
	secured.HandleFunc("/users/logout", controller.Logout).Methods("POST")
	secured.HandleFunc("/users/me", controller.Me).Methods("GET")

	// Task endpoints - for legacy task management
	// Frontend: Task creation endpoints (if used)
	secured.HandleFunc("/tasks", controller.GetTasks).Methods("GET")           // Fetch all tasks for user
	secured.HandleFunc("/tasks/create", controller.CreateTask).Methods("POST") // Create new task
	// REMOVED: secured.HandleFunc("/tasks/complete", controller.CompleteTask).Methods("POST") // Old endpoint removed to avoid route conflict

	// CHANGE: Daily task checklist endpoints - 5 tasks with 5-minute cooldown
	// Frontend: NormalTasks component calls these for daily checklist system
	secured.HandleFunc("/tasks/daily", controller.GetDailyTasks).Methods("GET")         // Fetch 5 daily tasks (auto-creates if needed)
	secured.HandleFunc("/tasks/complete", controller.CompleteTaskDaily).Methods("POST") // Complete task with cooldown validation
	secured.HandleFunc("/tasks/cooldown", controller.CheckCooldown).Methods("GET")      // Check cooldown status

	// Streak endpoints - for weekly check-in grid
	// Frontend: DailyStreak component calls these
	secured.HandleFunc("/streak", controller.GetStreak).Methods("GET")            // Fetch current streak data
	secured.HandleFunc("/streak/update", controller.UpdateStreak).Methods("POST") // Check in for today

	// Leaderboard endpoints - for ranking display
	// Frontend: LeaderboardBox component calls these
	secured.HandleFunc("/leaderboard", controller.GetLeaderboard).Methods("GET") // Fetch top users by points
	secured.HandleFunc("/leaderboard/me", controller.GetUserRank).Methods("GET") // Fetch logged-in user's rank

	// Legacy endpoints (kept for backward compatibility)
	router.HandleFunc("/users", controller.GetAlluser).Methods("GET")
	router.HandleFunc("/users/{id}", controller.Get1user).Methods("GET")
	router.HandleFunc("/users", controller.Create1user).Methods("POST")
	router.HandleFunc("/users/{id}", controller.Update1user).Methods("PUT")
	router.HandleFunc("/users/{id}", controller.Delete1user).Methods("DELETE")
	router.HandleFunc("/users", controller.DeleteAlluser).Methods("DELETE")

	return router
}
