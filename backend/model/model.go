package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ============ TASK MODELS ============

// Task represents a single daily task for a user
// Frontend: NormalTasks component expects { id, title, completed }
type Task struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId"`
	Title     string             `bson:"title" json:"title"`
	Completed bool               `bson:"completed" json:"completed"`
	Box       string             `bson:"box" json:"box"` // "left", "right", or "center" for task category
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// CHANGE: DailyTask represents a single task in the daily checklist (1 of 5)
// MongoDB collection: daily_tasks
// Used for: 5-task daily checklist system with cooldown enforcement
type DailyTask struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID      string             `bson:"user_id" json:"user_id"`
	TaskNumber  int                `bson:"task_number" json:"number"` // 1-5
	Completed   bool               `bson:"completed" json:"completed"`
	CompletedAt *time.Time         `bson:"completed_at,omitempty" json:"completed_at,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	ResetAt     time.Time          `bson:"reset_at" json:"reset_at"` // Midnight server time for TTL
}

// CHANGE: DailyTaskProgress tracks user's daily progress (cooldown, completion count)
// MongoDB collection: daily_task_progress
// Used for: Tracking completed tasks count and enforcing 5-minute cooldown
type DailyTaskProgress struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID          string             `bson:"user_id" json:"user_id"`
	CompletedCount  int                `bson:"completed_count" json:"completed_count"` // 0-5
	LastCompletedAt *time.Time         `bson:"last_completed_at,omitempty" json:"last_completed_at,omitempty"`
	LastCooldownEnd *time.Time         `bson:"last_cooldown_end,omitempty" json:"last_cooldown_end,omitempty"`
	NextResetAt     time.Time          `bson:"next_reset_at" json:"next_reset_at"` // Next midnight
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}

// ============ STREAK MODELS ============

// Streak represents a user's weekly check-in progress
// Frontend: DailyStreak component expects { mon, tue, wed, thu, fri, sat, sun }
type Streak struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID      primitive.ObjectID `bson:"userId" json:"userId"`
	Mon         bool               `bson:"mon" json:"mon"`
	Tue         bool               `bson:"tue" json:"tue"`
	Wed         bool               `bson:"wed" json:"wed"`
	Thu         bool               `bson:"thu" json:"thu"`
	Fri         bool               `bson:"fri" json:"fri"`
	Sat         bool               `bson:"sat" json:"sat"`
	Sun         bool               `bson:"sun" json:"sun"`
	LastCheckIn time.Time          `bson:"lastCheckIn" json:"lastCheckIn"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// ============ LEADERBOARD MODELS ============

// LeaderboardUser represents a user's ranking information
// Frontend: Leaderboard component displays { username, email, points, rank }
type LeaderboardUser struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username string             `bson:"username" json:"username"`
	Email    string             `bson:"email" json:"email"`
	Points   int                `bson:"points" json:"points"`
	Rank     int                `bson:"rank" json:"rank"`
}

// ============ USER MODELS (EXISTING) ============

// UserInput for handling user registration/login input (includes password and role)
type UserInput struct {
	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password,omitempty"`
	Role     string `json:"role,omitempty" bson:"role"` // Added for role-based authorization
}

// UserOutput for responses (excludes password)
type UserOutput struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Username string             `json:"username" bson:"username"`
	Email    string             `json:"email" bson:"email"`
	Role     string             `json:"role" bson:"role"` // Added role to output
}

// User for database operations (full struct)
type User struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Username string             `json:"username" bson:"username"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"-" bson:"password,omitempty"`
	Role     string             `json:"role" bson:"role"`
	Points   int                `bson:"points" json:"points"` // Points for leaderboard
}

// BlacklistedToken for logout functionality
// Added for token blacklisting on logout
type BlacklistedToken struct {
	ID    primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Token string             `json:"token" bson:"token"`
}
