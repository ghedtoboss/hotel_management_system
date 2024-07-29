package controllers

import (
	"encoding/json"
	"hotel_management_system/database"
	"hotel_management_system/models"
	"net/http"
	"time"
)

type OccupancyInput struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// Occupancy godoc
// @Summary Get hotel occupancy information
// @Description Get the number of total, occupied, and available rooms in a given date range
// @Tags Statistics
// @Accept  json
// @Produce  json
// @Param   input  body  OccupancyInput  true  "Date range for occupancy check"
// @Success 200 {object} map[string]int64
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Internal server error"
// @Router /occupancy [post]
func Occupancy(w http.ResponseWriter, r *http.Request) {

	type Input struct {
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
	}

	var input Input
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input.", http.StatusBadRequest)
		return
	}

	var reservations []models.Reservation
	if result := database.DB.Where("start_date < ? AND end_date > ?", input.EndDate, input.StartDate).Find(&reservations); result.Error != nil {
		http.Error(w, "Failed to fetch reservations.", http.StatusInternalServerError)
		return
	}

	//Determining the number of occupied rooms
	occupiedRooms := make(map[uint]bool)
	for _, reservation := range reservations {
		occupiedRooms[reservation.RoomID] = true
	}

	//Determining the total number of rooms
	var totalRooms int64
	if result := database.DB.Model(&models.Room{}).Count(&totalRooms); result.Error != nil {
		http.Error(w, "Failed to count rooms."+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	occupiedRoomCount := len(occupiedRooms)

	availableRoomCount := totalRooms - int64(occupiedRoomCount)

	result := map[string]int64{
		"total_rooms":     totalRooms,
		"occupied_rooms":  int64(occupiedRoomCount),
		"available_rooms": availableRoomCount,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}




