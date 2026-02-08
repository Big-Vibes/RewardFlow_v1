package service

import (
	"context"
	"fmt"
	"rewardpage/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	collection *mongo.Collection
}

// Added BlacklistService for token blacklisting
type BlacklistService struct {
	collection *mongo.Collection
}

func NewUserService(collection *mongo.Collection) *UserService {
	return &UserService{collection: collection}
}

// Added NewBlacklistService for initializing blacklist service
func NewBlacklistService(collection *mongo.Collection) *BlacklistService {
	return &BlacklistService{collection: collection}
}

// CreateUser creates a new user with password hashing and duplicate email check
func (us *UserService) CreateUser(ctx context.Context, user model.UserInput) error {
	// Check if email already exists
	count, err := us.collection.CountDocuments(ctx, bson.M{"email": user.Email})
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	// Set default role if not provided
	// Added for role-based authorization
	if user.Role == "" {
		user.Role = "user"
	}

	_, err = us.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

// GetAllUsers retrieves all users from the database
func (us *UserService) GetAllUsers(ctx context.Context) ([]bson.M, error) {
	cursor, err := us.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var users []bson.M
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	return users, nil
}

// GetUserByID retrieves a single user by ID
func (us *UserService) GetUserByID(ctx context.Context, userID string) (bson.M, error) {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID")
	}

	filter := bson.M{"_id": id}
	var user bson.M

	err = us.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser updates a user by ID
func (us *UserService) UpdateUser(ctx context.Context, userID string, updateData map[string]interface{}) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID")
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": updateData}

	result, err := us.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// DeleteUser deletes a user by ID
func (us *UserService) DeleteUser(ctx context.Context, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID")
	}

	filter := bson.M{"_id": id}
	result, err := us.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// DeleteAllUsers deletes all users from the database
func (us *UserService) DeleteAllUsers(ctx context.Context) (int64, error) {
	result, err := us.collection.DeleteMany(ctx, bson.D{})
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}

// FindUserByEmail retrieves a user by email address
// Added to fix undefined method error in auth_controller.go
// This method queries the database for a user matching the provided email
func (us *UserService) FindUserByEmail(ctx context.Context, email string, user *model.User) error {
	filter := bson.M{"email": email}
	err := us.collection.FindOne(ctx, filter).Decode(user)
	if err != nil {
		return err
	}
	return nil
}

// Added FindUserByID for internal use, returns model.User
func (us *UserService) FindUserByID(ctx context.Context, userID string, user *model.User) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID")
	}

	filter := bson.M{"_id": id}
	err = us.collection.FindOne(ctx, filter).Decode(user)
	if err != nil {
		return err
	}
	return nil
}

// Added BlacklistToken for logout functionality
func (bs *BlacklistService) BlacklistToken(ctx context.Context, token string) error {
	blacklisted := model.BlacklistedToken{Token: token}
	_, err := bs.collection.InsertOne(ctx, blacklisted)
	return err
}

// Added IsTokenBlacklisted for checking if token is blacklisted
func (bs *BlacklistService) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	filter := bson.M{"token": token}
	count, err := bs.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
