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

// GetTotalRevenue godoc
// @Summary Get total revenue for a date range
// @Description Get the total revenue of the hotel for a given date range
// @Tags Statistics
// @Accept  json
// @Produce  json
// @Param   input  body  RevenueInput  true  "Date range for revenue calculation"
// @Success 200 {object} map[string]float64
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Internal server error"
// @Router /revenue/total [post]
func GetTotalRevenue(w http.ResponseWriter, r *http.Request) {
	type RevenueInput struct {
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
	}

	var input RevenueInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input.", http.StatusBadRequest)
		return
	}

	var totalRevenue float64
	if result := database.DB.Model(&models.Reservation{}).
		Select("sum(rooms.price) as revenue").
		Joins("left join rooms on reservations.room_id = rooms.id").
		Where("reservations.start_date >= ? AND reservations.end_date <= ? AND reservations.status IN ?", input.StartDate, input.EndDate, []string{"confirmed", "checked-in", "checked-out"}).
		Scan(&totalRevenue); result.Error != nil {
		http.Error(w, "Failed to calculate total revenue.", http.StatusInternalServerError)
		return
	}

	result := map[string]float64{
		"total_revenue": totalRevenue,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// GetDailyRevenue godoc
// @Summary Get daily revenue for a date range
// @Description Get the daily revenue of the hotel for a given date range
// @Tags Statistics
// @Accept  json
// @Produce  json
// @Param   input  body  RevenueInput  true  "Date range for revenue calculation"
// @Success 200 {object} map[string]float64
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Internal server error"
// @Router /revenue/daily [post]
func GetDailyRevenue(w http.ResponseWriter, r *http.Request) {
	type RevenueInput struct {
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
	}

	var input RevenueInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input.", http.StatusBadRequest)
		return
	}

	var dailyRevenues []struct {
		Date    time.Time
		Revenue float64
	}

	if result := database.DB.Model(&models.Reservation{}).
		Select("date(start_date) as date, sum(rooms.price) as revenue").
		Joins("left join rooms on reservations.room_id = rooms.id").
		Where("reservations.start_date >= ? AND reservations.end_date <= ? AND reservations.status IN ?", input.StartDate, input.EndDate, []string{"confirmed", "checked-in", "checked-out"}).
		Group("date(start_date)").
		Scan(&dailyRevenues); result.Error != nil {
		http.Error(w, "Failed to calculate daily revenues.", http.StatusInternalServerError)
		return
	}

	result := make(map[string]float64)
	for _, dailyRevenue := range dailyRevenues {
		result[dailyRevenue.Date.Format("2006-01-02")] = dailyRevenue.Revenue
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// GetMonthlyRevenue godoc
// @Summary Get monthly revenue for a date range
// @Description Get the monthly revenue of the hotel for a given date range
// @Tags Statistics
// @Accept  json
// @Produce  json
// @Param   input  body  RevenueInput  true  "Date range for revenue calculation"
// @Success 200 {object} map[string]float64
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Internal server error"
// @Router /revenue/monthly [post]
func GetMonthlyRevenue(w http.ResponseWriter, r *http.Request) {
	type RevenueInput struct {
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
	}

	var input RevenueInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input.", http.StatusBadRequest)
		return
	}

	var monthlyRevenues []struct {
		Month   string
		Revenue float64
	}

	if result := database.DB.Model(&models.Reservation{}).
		Select("DATE_FORMAT(start_date, '%Y-%m') as month, SUM(rooms.price) as revenue").
		Joins("left join rooms on reservations.room_id = rooms.id").
		Where("reservations.start_date >= ? AND reservations.end_date <= ? AND reservations.status IN ?", input.StartDate, input.EndDate, []string{"confirmed", "checked-in", "checked-out"}).
		Group("month").
		Scan(&monthlyRevenues); result.Error != nil {
		http.Error(w, "Failed to calculate monthly revenues.", http.StatusInternalServerError)
		return
	}

	result := make(map[string]float64)
	for _, monthlyRevenue := range monthlyRevenues {
		result[monthlyRevenue.Month] = monthlyRevenue.Revenue
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
