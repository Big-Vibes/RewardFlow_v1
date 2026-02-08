package service

import (
	"context"
	"fmt"
	"rewardpage/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TaskService handles all task-related database operations
// Frontend integration: Called by task_controller to handle task CRUD and completion
type TaskService struct {
	collection *mongo.Collection
}

// NewTaskService creates a new TaskService instance
func NewTaskService(collection *mongo.Collection) *TaskService {
	return &TaskService{collection: collection}
}

// GetTasksByUserID retrieves all tasks for a specific user
// Used by frontend GET /api/tasks endpoint to populate NormalTasks component
func (ts *TaskService) GetTasksByUserID(ctx context.Context, userID string) ([]model.Task, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID")
	}

	filter := bson.M{"userId": userObjID}
	cursor, err := ts.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []model.Task
	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// CompleteTask marks a task as completed for a user
// Called by frontend POST /api/tasks/complete endpoint
// Parameters:
// - userID: extracted from JWT claims in middleware
// - taskID: sent in request body
// Returns error if task not found or user doesn't own task
func (ts *TaskService) CompleteTask(ctx context.Context, userID, taskID string) error {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID")
	}

	taskObjID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return fmt.Errorf("invalid task ID")
	}

	// Filter ensures user can only complete their own tasks (security)
	filter := bson.M{
		"_id":    taskObjID,
		"userId": userObjID,
	}

	update := bson.M{
		"$set": bson.M{
			"completed": true,
			"updatedAt": time.Now(),
		},
	}

	result, err := ts.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("task not found or does not belong to user")
	}

	return nil
}

// CreateTask adds a new task for a user
// Could be used in future for admin-created daily tasks
func (ts *TaskService) CreateTask(ctx context.Context, task *model.Task) error {
	task.ID = primitive.NewObjectID()
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	task.Completed = false

	_, err := ts.collection.InsertOne(ctx, task)
	return err
}

// ResetDailyTasks resets all tasks for all users (called daily)
// Could be run via cron job to reset the task checklist each day
func (ts *TaskService) ResetDailyTasks(ctx context.Context) error {
	filter := bson.M{}
	update := bson.M{
		"$set": bson.M{
			"completed": false,
			"updatedAt": time.Now(),
		},
	}

	_, err := ts.collection.UpdateMany(ctx, filter, update)
	return err
}

// DeleteTask removes a task (for cleanup or admin use)
func (ts *TaskService) DeleteTask(ctx context.Context, userID, taskID string) error {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID")
	}

	taskObjID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return fmt.Errorf("invalid task ID")
	}

	filter := bson.M{
		"_id":    taskObjID,
		"userId": userObjID,
	}

	result, err := ts.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}
