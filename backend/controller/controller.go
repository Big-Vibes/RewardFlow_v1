package controller

import (
	"encoding/json"
	"net/http"
	"rewardpage/model"
	"rewardpage/service"

	"github.com/gorilla/mux"
)

// GetAlluser retrieves all users from the database
func GetAlluser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users, err := service.UserServiceInstance.GetAllUsers(r.Context())
	if err != nil {
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

// Get1user retrieves a single user by ID
func Get1user(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	userID := vars["id"]

	user, err := service.UserServiceInstance.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// Create1user creates a new user
func Create1user(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user model.UserInput
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if user.Username == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "Please provide all required fields", http.StatusBadRequest)
		return
	}

	err := service.UserServiceInstance.CreateUser(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

// Update1user updates an existing user
func Update1user(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	userID := vars["id"]

	var updateData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := service.UserServiceInstance.UpdateUser(r.Context(), userID, updateData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully"})
}

// Delete1user deletes a single user by ID
func Delete1user(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	userID := vars["id"]

	err := service.UserServiceInstance.DeleteUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}

// DeleteAlluser deletes all users from the database
func DeleteAlluser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	count, err := service.UserServiceInstance.DeleteAllUsers(r.Context())
	if err != nil {
		http.Error(w, "Error deleting users", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "All users deleted successfully",
		"count":   count,
	})
}
