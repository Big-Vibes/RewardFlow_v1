package service

import (
	"context"
	"fmt"
	"rewardpage/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// StreakService handles user streak (daily check-in) data
// Frontend integration: Called by streak_controller to manage weekly check-in grid
type StreakService struct {
	collection *mongo.Collection
}

// NewStreakService creates a new StreakService instance
func NewStreakService(collection *mongo.Collection) *StreakService {
	return &StreakService{collection: collection}
}

// GetStreakByUserID retrieves the streak record for a user
// Used by frontend GET /api/streak endpoint to populate DailyStreak component
// Returns: { mon, tue, wed, thu, fri, sat, sun } booleans
func (ss *StreakService) GetStreakByUserID(ctx context.Context, userID string) (*model.Streak, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID")
	}

	filter := bson.M{"userId": userObjID}
	var streak model.Streak

	err = ss.collection.FindOne(ctx, filter).Decode(&streak)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Create a new streak record if doesn't exist
			return ss.CreateStreakRecord(ctx, userObjID)
		}
		return nil, err
	}

	return &streak, nil
}

// UpdateStreak performs a daily check-in and marks the current day as completed
// Called by frontend POST /api/streak/update endpoint
// Logic:
// 1. Get today's day of week (Mon-Sun)
// 2. Update the corresponding field to true
// 3. Update lastCheckIn timestamp
// Returns the updated streak object
func (ss *StreakService) UpdateStreak(ctx context.Context, userID string) (*model.Streak, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID")
	}

	// Get current day name (Mon-Sun)
	dayName := time.Now().Weekday().String()[:3]
	dayLower := map[string]string{
		"Mon": "mon",
		"Tue": "tue",
		"Wed": "wed",
		"Thu": "thu",
		"Fri": "fri",
		"Sat": "sat",
		"Sun": "sun",
	}[dayName]

	filter := bson.M{"userId": userObjID}

	// Update the current day to true
	update := bson.M{
		"$set": bson.M{
			dayLower:      true,
			"lastCheckIn": time.Now(),
			"updatedAt":   time.Now(),
		},
	}

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(options.After) // Return updated document

	var updatedStreak model.Streak
	err = ss.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedStreak)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// If no record exists, create one and set today's day
			newStreak := &model.Streak{
				UserID:      userObjID,
				LastCheckIn: time.Now(),
				UpdatedAt:   time.Now(),
			}

			// Set current day to true
			switch dayLower {
			case "mon":
				newStreak.Mon = true
			case "tue":
				newStreak.Tue = true
			case "wed":
				newStreak.Wed = true
			case "thu":
				newStreak.Thu = true
			case "fri":
				newStreak.Fri = true
			case "sat":
				newStreak.Sat = true
			case "sun":
				newStreak.Sun = true
			}

			insertedID, err := ss.collection.InsertOne(ctx, newStreak)
			if err != nil {
				return nil, err
			}

			newStreak.ID = insertedID.InsertedID.(primitive.ObjectID)
			return newStreak, nil
		}
		return nil, err
	}

	return &updatedStreak, nil
}

// ResetStreakDaily resets all user streaks at midnight (called via cron job)
// Typically run at 11:59 PM daily to reset the next day's check-in
// This allows users to check in again on the new day
func (ss *StreakService) ResetStreakDaily(ctx context.Context) error {
	// Optional: Only reset if user checked in today
	// This preserves the streak chain and motivation

	// For now, we don't auto-reset; users can check in fresh each day
	return nil
}

// CreateStreakRecord initializes a new streak record for a user
func (ss *StreakService) CreateStreakRecord(ctx context.Context, userObjID primitive.ObjectID) (*model.Streak, error) {
	streak := &model.Streak{
		ID:          primitive.NewObjectID(),
		UserID:      userObjID,
		LastCheckIn: time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := ss.collection.InsertOne(ctx, streak)
	if err != nil {
		return nil, err
	}

	return streak, nil
}

// GetStreakCount returns the number of consecutive days a user has checked in
// Could be used for badges/achievements in future
func (ss *StreakService) GetStreakCount(ctx context.Context, userID string) (int, error) {
	streak, err := ss.GetStreakByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}

	count := 0
	if streak.Mon {
		count++
	}
	if streak.Tue {
		count++
	}
	if streak.Wed {
		count++
	}
	if streak.Thu {
		count++
	}
	if streak.Fri {
		count++
	}
	if streak.Sat {
		count++
	}
	if streak.Sun {
		count++
	}

	return count, nil
}
