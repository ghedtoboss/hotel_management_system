package controllers

import (
	"encoding/json"
	"hotel_management_system/database"
	"hotel_management_system/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// GetCustomers godoc
// @Summary Get all customers
// @Description Get a list of all users with the role of customer
// @Tags User
// @Produce  json
// @Success 200 {array} models.User
// @Failure 404 {string} string "Customers not found"
// @Failure 500 {string} string "Internal server error"
// @Router /customers [get]
func GetCustomers(w http.ResponseWriter, r *http.Request) {
	var customers []models.User
	if result := database.DB.Where("role = ?", "customer").Find(&customers); result.Error != nil {
		http.Error(w, "Customers not found.", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customers)
}

// GetUser godoc
// @Summary Get user details
// @Description Get details of a specific user by ID
// @Tags User
// @Produce  json
// @Param   user_id  path int  true  "User ID"
// @Success 200 {object} models.User
// @Failure 400 {string} string "Invalid user id"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal server error"
// @Router /users/{user_id} [get]
func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["user_id"])
	if err != nil {
		http.Error(w, "Invalid user id.", http.StatusBadRequest)
		return
	}

	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		http.Error(w, "User not found.", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update details of a specific user by ID
// @Tags User
// @Accept  json
// @Produce  json
// @Param   user_id  path int  true  "User ID"
// @Param   user  body models.User  true  "Updated user data"
// @Success 200 {object} models.User
// @Failure 400 {string} string "Invalid user id"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal server error"
// @Router /users/{user_id} [put]
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["user_id"])
	if err != nil {
		http.Error(w, "Invalid user id.", http.StatusBadRequest)
		return
	}

	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		http.Error(w, "User not found.", http.StatusNotFound)
		return
	}

	var input map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input.", http.StatusInternalServerError)
		return
	}

	if newPassword, ok := input["password"].(string); ok {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash password.", http.StatusInternalServerError)
			return
		}
		user.Password = string(hashedPassword)
	}

	if newEmail, ok := input["email"].(string); ok {
		user.Email = newEmail
	}

	if newUsername, ok := input["username"].(string); ok {
		user.Username = newUsername
	}

	user.UpdatedAt = time.Now()

	if result := database.DB.Save(&user); result.Error != nil {
		http.Error(w, "Failed to update user.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully."})
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user by ID
// @Tags User
// @Param   user_id  path int  true  "User ID"
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string "Invalid user id"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal server error"
// @Router /users/{user_id} [delete]
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["user_id"])
	if err != nil {
		http.Error(w, "Invalid user id.", http.StatusBadRequest)
		return
	}

	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		http.Error(w, "User not found.", http.StatusNotFound)
		return
	}

	if result := database.DB.Delete(&user); result.Error != nil {
		http.Error(w, "Failed to delete user.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully."})
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Get a list of all users
// @Tags User
// @Produce  json
// @Success 200 {array} models.User
// @Failure 404 {string} string "Users not found"
// @Failure 500 {string} string "Internal server error"
// @Router /users [get]
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	if result := database.DB.Find(&users); result.Error != nil {
		http.Error(w, "Users not found.", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}
