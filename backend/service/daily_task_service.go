package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"rewardpage/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CHANGE: DailyTaskService handles 5-task daily checklist with cooldown enforcement
type DailyTaskService struct {
	DB *mongo.Database
}

// CHANGE: Initialize daily task service with TTL index for auto-cleanup
func InitDailyTaskService(db *mongo.Database) {
	DailyTaskServiceInstance = &DailyTaskService{DB: db}

	// CHANGE: Create TTL index on ResetAt field for automatic cleanup of old tasks
	// Prevents MongoDB from storing old daily tasks indefinitely
	dailyTaskCollection := db.Collection("daily_tasks")
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "reset_at", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(86400), // 24 hours
	}
	dailyTaskCollection.Indexes().CreateOne(context.Background(), indexModel)
}

// CHANGE: GetOrCreateDailyTasks retrieves or creates today's 5 tasks for user
// Logic:
// 1. Query tasks created today (between midnight and tomorrow midnight)
// 2. If none exist, create 5 new tasks
// 3. Return array of 5 task objects
func (s *DailyTaskService) GetOrCreateDailyTasks(ctx context.Context, userID string) ([]model.DailyTask, error) {
	collection := s.DB.Collection("daily_tasks")
	today := getTodayMidnight()
	tomorrow := today.AddDate(0, 0, 1)

	// CHANGE: Find tasks created today (between midnight and tomorrow midnight)
	filter := bson.M{
		"user_id": userID,
		"reset_at": bson.M{
			"$gte": today,
			"$lt":  tomorrow,
		},
	}

	opts := options.Find().SetSort(bson.M{"task_number": 1})
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []model.DailyTask
	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	// CHANGE: If no tasks exist for today, create all 5 tasks
	if len(tasks) == 0 {
		tasks, err = s.createFiveTasks(ctx, userID, today, tomorrow)
		if err != nil {
			return nil, err
		}
	}

	return tasks, nil
}

// CHANGE: createFiveTasks creates 5 new task documents in MongoDB
// Each task has:
// - TaskNumber: 1-5 (visual display)
// - Completed: false (initial state)
// - ResetAt: tomorrow midnight (TTL cleanup)
func (s *DailyTaskService) createFiveTasks(ctx context.Context, userID string, today, tomorrow time.Time) ([]model.DailyTask, error) {
	collection := s.DB.Collection("daily_tasks")
	var tasks []model.DailyTask

	// CHANGE: Create exactly 5 task documents
	for i := 1; i <= 5; i++ {
		task := model.DailyTask{
			UserID:     userID,
			TaskNumber: i,
			Completed:  false,
			CreatedAt:  time.Now(),
			ResetAt:    tomorrow,
		}
		result, err := collection.InsertOne(ctx, task)
		if err != nil {
			return nil, err
		}
		task.ID = result.InsertedID.(primitive.ObjectID)
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// CHANGE: CompleteTask marks a task as completed with strict backend validation
// Validation rules:
// 1. Check if user is within cooldown (5 minutes since last task)
// 2. Check if user already completed 5 tasks today
// 3. Check if tasks need daily reset (past midnight)
// 4. Update task.completed = true and task.completedAt = now
// 5. Update progress tracking with new cooldown
func (s *DailyTaskService) CompleteTask(ctx context.Context, userID, taskID string) (map[string]interface{}, error) {
	log.Printf("CompleteTask starting: userID=%s, taskID=%s", userID, taskID)

	// CHANGE: First, check if tasks need daily reset (past midnight)
	if err := s.CheckAndResetDaily(ctx, userID); err != nil {
		log.Printf("CheckAndResetDaily failed: %v", err)
		return nil, err
	}

	// CHANGE: Get current progress to check cooldown and completion limit
	progress, err := s.GetOrCreateProgress(ctx, userID)
	if err != nil {
		log.Printf("GetOrCreateProgress failed: %v", err)
		return nil, err
	}
	log.Printf("Progress retrieved: completedCount=%d, lastCompletedAt=%v", progress.CompletedCount, progress.LastCompletedAt)

	// CHANGE: Check if user already completed 5 tasks today
	if progress.CompletedCount >= 5 {
		return nil, errors.New("all 5 daily tasks already completed")
	}

	// CHANGE: Check if user is within 5-minute cooldown period
	// Backend prevents spam by enforcing this server-side
	now := time.Now()
	if progress.LastCompletedAt != nil {
		elapsed := now.Sub(*progress.LastCompletedAt)
		if elapsed < 5*time.Minute {
			remainingSeconds := int(((5 * time.Minute) - elapsed).Seconds())
			return map[string]interface{}{
				"error":             fmt.Sprintf("cooldown active, wait %d seconds", remainingSeconds),
				"remaining_seconds": remainingSeconds,
			}, errors.New("cooldown period active")
		}
	}

	// CHANGE: Convert taskID string to MongoDB ObjectID
	objID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		log.Printf("Invalid task ID format: %s, error: %v", taskID, err)
		return nil, errors.New("invalid task id format")
	}
	log.Printf("Task ID converted: %s -> %s", taskID, objID.Hex())

	// CHANGE: Update task to completed with timestamp
	taskCollection := s.DB.Collection("daily_tasks")
	now = time.Now()
	update := bson.M{
		"$set": bson.M{
			"completed":    true,
			"completed_at": now,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedTask model.DailyTask
	err = taskCollection.FindOneAndUpdate(ctx, bson.M{"_id": objID}, update, opts).Decode(&updatedTask)
	if err != nil {
		log.Printf("Failed to update task in database: %v", err)
		return nil, err
	}
	log.Printf("Task updated successfully: taskID=%s, completed=%v", objID.Hex(), updatedTask.Completed)

	// CHANGE: Update progress tracking with new completion count and cooldown
	progressCollection := s.DB.Collection("daily_task_progress")
	newCompletedCount := progress.CompletedCount + 1
	nextResetAt := getTodayMidnight().AddDate(0, 0, 1)

	progressUpdate := bson.M{
		"$set": bson.M{
			"completed_count":   newCompletedCount,
			"last_completed_at": now,
			"last_cooldown_end": now.Add(5 * time.Minute),
			"next_reset_at":     nextResetAt,
			"updated_at":        now,
		},
	}

	_, err = progressCollection.UpdateOne(ctx, bson.M{"user_id": userID}, progressUpdate)
	if err != nil {
		return nil, err
	}

	// CHANGE: Return updated state to frontend
	return map[string]interface{}{
		"success":         true,
		"task":            updatedTask,
		"completed_count": newCompletedCount,
		"next_reset_at":   nextResetAt,
		"cooldown_until":  now.Add(5 * time.Minute),
	}, nil
}

// CHANGE: CheckAndResetDaily checks if past midnight and resets all tasks
// This is called on every GET /api/tasks/daily request
// Logic:
// 1. Get current progress and check NextResetAt time
// 2. If NextResetAt < now, we're past midnight
// 3. Delete old tasks from today
// 4. Create 5 new tasks for today
// 5. Reset progress counters
// Used for: Lazy-loaded daily reset on next user request after midnight
func (s *DailyTaskService) CheckAndResetDaily(ctx context.Context, userID string) error {
	progress, err := s.GetOrCreateProgress(ctx, userID)
	if err != nil {
		return err
	}

	now := time.Now()
	today := getTodayMidnight()

	// CHANGE: If next reset time is in the past, we've passed midnight, so reset all tasks
	if progress.NextResetAt.Before(now) {
		collection := s.DB.Collection("daily_tasks")

		// CHANGE: Delete old tasks (before today's midnight)
		_, err := collection.DeleteMany(ctx, bson.M{
			"user_id": userID,
			"reset_at": bson.M{
				"$lt": today,
			},
		})
		if err != nil {
			return err
		}

		// CHANGE: Create new 5 tasks for today
		tomorrow := today.AddDate(0, 0, 1)
		_, err = s.createFiveTasks(ctx, userID, today, tomorrow)
		if err != nil {
			return err
		}

		// CHANGE: Reset progress counters for new day
		progressCollection := s.DB.Collection("daily_task_progress")
		progressUpdate := bson.M{
			"$set": bson.M{
				"completed_count":   0,
				"last_completed_at": nil,
				"last_cooldown_end": nil,
				"next_reset_at":     tomorrow,
				"updated_at":        now,
			},
		}
		_, err = progressCollection.UpdateOne(ctx, bson.M{"user_id": userID}, progressUpdate)
		if err != nil {
			return err
		}

		log.Printf("Daily reset executed for user %s at %v", userID, now)
	}

	return err
}

// CHANGE: GetOrCreateProgress retrieves or creates progress tracking for user
// Returns DailyTaskProgress which contains:
// - CompletedCount: number of tasks completed today (0-5)
// - LastCompletedAt: timestamp of last completed task
// - NextResetAt: next midnight (when tasks reset)
func (s *DailyTaskService) GetOrCreateProgress(ctx context.Context, userID string) (*model.DailyTaskProgress, error) {
	collection := s.DB.Collection("daily_task_progress")
	today := getTodayMidnight()
	tomorrow := today.AddDate(0, 0, 1)

	filter := bson.M{"user_id": userID}
	var progress model.DailyTaskProgress
	err := collection.FindOne(ctx, filter).Decode(&progress)

	if err == mongo.ErrNoDocuments {
		// CHANGE: Create new progress if doesn't exist
		progress = model.DailyTaskProgress{
			UserID:         userID,
			CompletedCount: 0,
			NextResetAt:    tomorrow,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		result, err := collection.InsertOne(ctx, progress)
		if err != nil {
			return nil, err
		}
		progress.ID = result.InsertedID.(primitive.ObjectID)
	} else if err != nil {
		return nil, err
	}

	return &progress, nil
}

// CHANGE: getTodayMidnight returns today's midnight in server time
// Used for: Determining daily reset boundaries and TTL cleanup
func getTodayMidnight() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}
