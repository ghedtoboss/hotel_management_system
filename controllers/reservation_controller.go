package controllers

import (
	"encoding/json"
	"hotel_management_system/database"
	"hotel_management_system/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// CreateReservation godoc
// @Summary Create a new reservation
// @Description Create a new reservation for a room
// @Tags Reservation
// @Accept  json
// @Produce  json
// @Param   reservation  body models.Reservation  true  "Reservation data"
// @Success 201 {object} models.Reservation
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Internal server error"
// @Router /reservations [post]
func CreateReservation(w http.ResponseWriter, r *http.Request) {
	var input map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input.", http.StatusBadRequest)
		return
	}

	roomNumber, ok := input["room_number"].(string)
	if !ok {
		http.Error(w, "Invalid room number", http.StatusBadRequest)
		return
	}
	var room models.Room
	if result := database.DB.Where("number = ?", roomNumber).First(&room); result.Error != nil {
		http.Error(w, "Room not found.", http.StatusNotFound)
		return
	}

	startDateStr, ok := input["start_date"].(string)
	if !ok {
		http.Error(w, "Invalid start date", http.StatusBadRequest)
		return
	}
	endDateStr, ok := input["end_date"].(string)
	if !ok {
		http.Error(w, "Invalid end date", http.StatusBadRequest)
		return
	}
	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		http.Error(w, "Invalid start date format.", http.StatusBadRequest)
		return
	}
	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		http.Error(w, "Invalid end date format.", http.StatusBadRequest)
		return
	}

	userID, ok := input["user_id"].(float64)
	if !ok {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	reservation := models.Reservation{
		UserID:    uint(userID),
		RoomID:    room.ID,
		StartDate: startDate,
		EndDate:   endDate,
		Status:    "confirmed",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Check for conflicting reservations
	var existingReservations []models.Reservation
	database.DB.Where("room_id = ? AND start_date < ? AND end_date > ?", room.ID, endDate, startDate).Find(&existingReservations)
	if len(existingReservations) > 0 {
		http.Error(w, "Reservation dates conflict with an existing reservation", http.StatusConflict)
		return
	}

	if result := database.DB.Create(&reservation); result.Error != nil {
		http.Error(w, "Failed to create reservation.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(reservation)
}

// UpdateReservation godoc
// @Summary Update an existing reservation
// @Description Update an existing reservation
// @Tags Reservation
// @Accept  json
// @Produce  json
// @Param   reservation_id  path int  true  "Reservation ID"
// @Param   reservation  body models.Reservation  true  "Updated reservation data"
// @Success 200 {object} models.Reservation
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Reservation not found"
// @Failure 500 {string} string "Internal server error"
// @Router /reservations/{reservation_id} [put]
func UpdateReservation(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	reservID, err := strconv.Atoi(params["reservation_id"])
	if err != nil {
		http.Error(w, "Invalid reservation id", http.StatusBadRequest)
		return
	}

	var reservation models.Reservation
	if result := database.DB.First(&reservation, reservID); result.Error != nil {
		http.Error(w, "Reservation not found.", http.StatusNoContent)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&reservation)
	if err != nil {
		http.Error(w, "Invalid input.", http.StatusInternalServerError)
		return
	}

	reservation.UpdatedAt = time.Now()

	if result := database.DB.Save(&reservation); result.Error != nil {
		http.Error(w, "Failed to update reservation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Reservation updated successfully."})
}

// DeleteReservation godoc
// @Summary Delete a reservation
// @Description Delete a reservation by ID
// @Tags Reservation
// @Param   reservation_id  path int  true  "Reservation ID"
// @Success 204 {string} string "No Content"
// @Failure 404 {string} string "Reservation not found"
// @Failure 500 {string} string "Internal server error"
// @Router /reservations/{reservation_id} [delete]
func DeleteReservation(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	reservID, err := strconv.Atoi(params["reservation_id"])
	if err != nil {
		http.Error(w, "Invalid reservation id", http.StatusBadRequest)
		return
	}

	var reservation models.Reservation
	if result := database.DB.First(&reservation, reservID); result.Error != nil {
		http.Error(w, "Reservation not found.", http.StatusNoContent)
		return
	}

	if result := database.DB.Delete(&reservation); result.Error != nil {
		http.Error(w, "Failed to delete reservation.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(map[string]string{"message": "Reservation deleted successfully."})
}

// GetReservations godoc
// @Summary Get all reservations
// @Description Get a list of all reservations
// @Tags Reservation
// @Produce  json
// @Success 200 {array} models.Reservation
// @Failure 404 {string} string "Reservations not found"
// @Failure 500 {string} string "Internal server error"
// @Router /reservations [get]
func GetReservations(w http.ResponseWriter, r *http.Request) {
	var reservations []models.Reservation
	if result := database.DB.Find(&reservations); result.Error != nil {
		http.Error(w, "Reservations not found.", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reservations)
}

// GetReservationDetails godoc
// @Summary Get reservation details
// @Description Get details of a specific reservation
// @Tags Reservation
// @Produce  json
// @Param   reservation_id  path int  true  "Reservation ID"
// @Success 200 {object} models.Reservation
// @Failure 400 {string} string "Invalid reservation ID"
// @Failure 404 {string} string "Reservation not found"
// @Failure 500 {string} string "Internal server error"
// @Router /reservations/{reservation_id} [get]
func GetReservationDetails(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	reservID, err := strconv.Atoi(params["reservation_id"])
	if err != nil {
		http.Error(w, "Invalid reservation id.", http.StatusBadRequest)
		return
	}

	var reservation models.Reservation
	if result := database.DB.Find(&reservation, reservID); result.Error != nil {
		http.Error(w, "Reservation not found.", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reservation)
}

// UpdateReservationStatus godoc
// @Summary Update reservation status
// @Description Update the status of a reservation
// @Tags Reservation
// @Accept  json
// @Produce  json
// @Param   reservation_id  path int  true  "Reservation ID"
// @Param   status  body string  true  "New status"
// @Success 200 {object} models.Reservation
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Reservation not found"
// @Failure 500 {string} string "Internal server error"
// @Router /reservations/{reservation_id}/status [put]
func UpdateReservationStatus(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	reservID, err := strconv.Atoi(params["reservation_id"])
	if err != nil {
		http.Error(w, "Invalid reservation id", http.StatusBadRequest)
		return
	}

	var input struct {
		Status string `json:"status"`
	}
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var reservation models.Reservation
	if result := database.DB.First(&reservation, reservID); result.Error != nil {
		http.Error(w, "Reservation not found", http.StatusNotFound)
		return
	}

	validStatuses := []string{"pending", "confirmed", "checked-in", "checked-out", "cancelled", "no-show"}
	isValidStatus := false
	for _, status := range validStatuses {
		if input.Status == status {
			isValidStatus = true
			break
		}
	}

	if !isValidStatus {
		http.Error(w, "Invalid status value", http.StatusBadRequest)
		return
	}

	reservation.Status = input.Status
	reservation.UpdatedAt = time.Now()

	if result := database.DB.Save(&reservation); result.Error != nil {
		http.Error(w, "Failed to update reservation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reservation)
}
