package service

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbName = "userdb"
const colName = "users"
const blacklistColName = "blacklisted_tokens" // Added for logout functionality
const tasksColName = "tasks"                  // Added for task management
const streaksColName = "streaks"              // Added for daily streak tracking
const leaderboardColName = "leaderboard"      // Reference to users collection for ranking

var UserServiceInstance *UserService
var BlacklistServiceInstance *BlacklistService     // Added for token blacklisting
var TaskServiceInstance *TaskService               // Added for task operations
var StreakServiceInstance *StreakService           // Added for streak operations
var DailyTaskServiceInstance *DailyTaskService     // Added for daily task checklist
var LeaderboardServiceInstance *LeaderboardService // Added for leaderboard ranking
var mongoClient *mongo.Client                      // CHANGE: Store mongo client for GetDB() access

// InitializeDB initializes MongoDB connection and all service instances
// Creates collections for users, tasks, streaks, and token blacklist
// Called once on application startup
func InitializeDB() error {
	clientOption := options.Client().ApplyURI(os.Getenv("MONGO_URI"))

	client, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		return err
	}

	// CHANGE: Store mongo client for GetDB() access
	mongoClient = client

	fmt.Println("MongoDB connection success")

	// Initialize collections
	userCollection := client.Database(dbName).Collection(colName)
	fmt.Println("User collection instance is ready")

	blacklistCollection := client.Database(dbName).Collection(blacklistColName)
	fmt.Println("Blacklist collection instance is ready")

	// Initialize task collection for daily task checklist
	tasksCollection := client.Database(dbName).Collection(tasksColName)
	fmt.Println("Tasks collection instance is ready")

	// Initialize streaks collection for weekly check-in grid
	streaksCollection := client.Database(dbName).Collection(streaksColName)
	fmt.Println("Streaks collection instance is ready")

	// Initialize service instances
	UserServiceInstance = NewUserService(userCollection)
	BlacklistServiceInstance = NewBlacklistService(blacklistCollection)
	TaskServiceInstance = NewTaskService(tasksCollection)
	StreakServiceInstance = NewStreakService(streaksCollection)
	LeaderboardServiceInstance = NewLeaderboardService(userCollection) // Uses users collection for points

	return nil
}

// CHANGE: GetDB returns the MongoDB database instance for DailyTaskService
// Used to initialize daily task service with proper MongoDB access
func GetDB() *mongo.Database {
	return mongoClient.Database(dbName)
}
