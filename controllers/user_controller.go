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

// GetProfile godoc
// @Summary Get user profile
// @Description Get the profile information of the currently logged-in user
// @Tags Profile
// @Produce  json
// @Success 200 {object} models.User
// @Failure 401 {string} string "Unauthorized"
// @Router /profile [get]
func GetProfile(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*models.Claims)

	var user models.User
	if result := database.DB.First(&user, claims.UserID); result.Error != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update the profile information of the currently logged-in user
// @Tags Profile
// @Accept  json
// @Produce  json
// @Param   user  body models.User  true  "Updated user data"
// @Success 200 {object} models.User
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /profile [put]
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*models.Claims)
	var user models.User
	if result := database.DB.First(&user, claims.UserID); result.Error != nil {
		http.Error(w, "User not found.", http.StatusNotFound)
		return
	}

	var input models.User
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input.", http.StatusBadRequest)
		return
	}

	user.Email = input.Email
	user.Username = input.Username
	user.UpdatedAt = time.Now()

	if result := database.DB.Save(&user); result != nil {
		http.Error(w, "Failed to update profile.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// UpdatePassword godoc
// @Summary Update user password
// @Description Update the password of the currently logged-in user
// @Tags Profile
// @Accept  json
// @Produce  json
// @Param   password_data  body map[string]string  true  "Old and new passwords"
// @Success 200 {string} string "Password updated successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /profile/password [put]
func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*models.Claims)
	var user models.User
	if result := database.DB.First(&user, claims.UserID); result.Error != nil {
		http.Error(w, "User not found.", http.StatusNotFound)
		return
	}

	var passwordData map[string]string
	err := json.NewDecoder(r.Body).Decode(&passwordData)
	if err != nil {
		http.Error(w, "Invalid input.", http.StatusBadRequest)
		return
	}

	oldPassword := passwordData["old_password"]
	newPassword := passwordData["new_password"]

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		http.Error(w, "Old password is incorrect.", http.StatusUnauthorized)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password.", http.StatusInternalServerError)
		return
	}

	user.Password = string(hashedPassword)
	user.UpdatedAt = time.Now()

	if result := database.DB.Save(&user); result.Error != nil {
		http.Error(w, "Failed to update password.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Password updated successfully."})
}
