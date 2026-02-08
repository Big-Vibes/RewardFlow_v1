package service

import (
	"context"
	"fmt"
	"rewardpage/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// LeaderboardService handles ranking and points calculation
// Frontend integration: Called by leaderboard_controller to fetch user rankings
type LeaderboardService struct {
	collection *mongo.Collection // users collection to fetch points
}

// NewLeaderboardService creates a new LeaderboardService instance
func NewLeaderboardService(collection *mongo.Collection) *LeaderboardService {
	return &LeaderboardService{collection: collection}
}

// GetLeaderboard returns top ranked users sorted by points (descending)
// Assigns sequential rank numbers starting from 1
// Called by frontend GET /api/leaderboard endpoint
// Parameters:
// - limit: number of top users to return (default 10, max 100)
// Returns: array of LeaderboardUser with rank, username, email, points
func (ls *LeaderboardService) GetLeaderboard(ctx context.Context, limit int64) ([]model.LeaderboardUser, error) {
	if limit <= 0 {
		limit = 10 // Default to top 10
	}
	if limit > 100 {
		limit = 100 // Cap at 100
	}

	// Sort by points descending (highest first), then by username for consistency
	opts := options.Find().SetSort(bson.M{"points": -1, "username": 1}).SetLimit(limit)

	cursor, err := ls.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var leaderboard []model.LeaderboardUser
	if err = cursor.All(ctx, &leaderboard); err != nil {
		return nil, err
	}

	// Assign sequential ranks (1, 2, 3, ...)
	for i := range leaderboard {
		leaderboard[i].Rank = i + 1
	}

	return leaderboard, nil
}

// GetUserRank returns a specific user's rank and points
// Could be used in future for "Your Position" display
func (ls *LeaderboardService) GetUserRank(ctx context.Context, userID string) (*model.LeaderboardUser, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID")
	}

	filter := bson.M{"_id": userObjID}
	var user model.LeaderboardUser

	err = ls.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	// Calculate rank by counting users with more points
	countHigher, err := ls.collection.CountDocuments(ctx, bson.M{"points": bson.M{"$gt": user.Points}})
	if err != nil {
		return nil, err
	}

	user.Rank = int(countHigher) + 1
	return &user, nil
}

// AddPointsToUser adds points to a user (called when task is completed or check-in successful)
// Used internally by controllers when tasks are completed
// Parameters:
// - userID: user who earned points
// - points: number of points to add
func (ls *LeaderboardService) AddPointsToUser(ctx context.Context, userID string, points int) error {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID")
	}

	filter := bson.M{"_id": userObjID}
	update := bson.M{
		"$inc": bson.M{"points": points}, // Increment points
	}

	result, err := ls.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// GetTopStreak returns users sorted by their current week streak count
// Could be used for special "Streak Leaders" leaderboard
func (ls *LeaderboardService) GetTopStreak(ctx context.Context, limit int64, streakService *StreakService) ([]model.LeaderboardUser, error) {
	if limit <= 0 {
		limit = 10
	}

	// Get all users
	cursor, err := ls.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []model.LeaderboardUser
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	// This is a simplified example - in production, you'd want to fetch streak counts
	// and sort in MongoDB for better performance
	// For now, just return top by points
	return ls.GetLeaderboard(ctx, limit)
}

// InitializeUserPoints ensures a user has a points field set to 0 on signup
// Called when a new user registers
func (ls *LeaderboardService) InitializeUserPoints(ctx context.Context, userID string) error {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID")
	}

	filter := bson.M{"_id": userObjID}
	update := bson.M{
		"$set": bson.M{"points": 0}, // Initialize to 0 if not set
	}

	// Use UpdateOne with upsert false to only update if exists
	_, err = ls.collection.UpdateOne(ctx, filter, update)
	return err
}
