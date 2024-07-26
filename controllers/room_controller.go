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

// CreateRoom godoc
// @Summary Create a new room
// @Description Create a new room with number, type, status, and price
// @Tags Room
// @Accept  json
// @Produce  json
// @Param   number  body string  true  "Room Number"
// @Param   type    body string  true  "Room Type"
// @Param   status  body string  true  "Room Status"
// @Param   price   body float64 true  "Room Price"
// @Success 201 {string} string "Room created successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Internal server error"
// @Router /rooms [post]
func CreateRoom(w http.ResponseWriter, r *http.Request) {
	var room models.Room
	err := json.NewDecoder(r.Body).Decode(&room)
	if err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	room.CreatedAt = time.Now()
	room.UpdateAt = time.Now()

	if result := database.DB.Create(&room); result.Error != nil {
		http.Error(w, "Failed to create room: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Room created successfully"})
}

// UpdateRoom godoc
// @Summary Update an existing room
// @Description Update an existing room with number, type, status, and price
// @Tags Room
// @Accept  json
// @Produce  json
// @Param   room_id  path int  true  "Room ID"
// @Param   number  body string  true  "Room Number"
// @Param   type    body string  true  "Room Type"
// @Param   status  body string  true  "Room Status"
// @Param   price   body float64 true  "Room Price"
// @Success 200 {string} string "Room updated successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Room not found"
// @Failure 500 {string} string "Internal server error"
// @Router /rooms/{room_id} [put]
func UpdateRoom(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	roomID := params["room_id"]

	var room models.Room
	if result := database.DB.First(&room, roomID); result.Error != nil {
		http.Error(w, "Room not found: "+result.Error.Error(), http.StatusNotFound)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&room)
	if err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	room.UpdateAt = time.Now()

	if result := database.DB.Save(&room); result.Error != nil {
		http.Error(w, "Failed to update room "+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": ":Room updated succesfully."})

}

// DeleteRoom godoc
// @Summary Delete a room
// @Description Delete a room by ID
// @Tags Room
// @Param   room_id  path int  true  "Room ID"
// @Success 204 {string} string "No Content"
// @Failure 404 {string} string "Room not found"
// @Failure 500 {string} string "Internal server error"
// @Router /rooms/{room_id} [delete]
func DeleteRoom(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	roomID, err := strconv.Atoi(params["room_id"])
	if err != nil {
		http.Error(w, "Invalid room id", http.StatusBadRequest)
		return
	}

	var room models.Room
	if result := database.DB.First(&room, roomID); result.Error != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	if result := database.DB.Delete(&room); result.Error != nil {
		http.Error(w, "Failed to delete room: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(map[string]string{"message": "Room deleted successfully."})
}

// GetRooms godoc
// @Summary Get all rooms
// @Description Get a list of all rooms
// @Tags Room
// @Produce  json
// @Success 200 {array} models.Room
// @Failure 500 {string} string "Internal server error"
// @Router /rooms [get]
func GetRooms(w http.ResponseWriter, r *http.Request) {

	var rooms []models.Room
	if result := database.DB.Find(&rooms); result.Error != nil {
		http.Error(w, "Rooms not found.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rooms)
}

// GetRoomDetails godoc
// @Summary Get room details
// @Description Get details of a specific room
// @Tags Room
// @Produce  json
// @Param   room_id  path int  true  "Room ID"
// @Success 200 {object} models.Room
// @Failure 400 {string} string "Invalid room ID"
// @Failure 404 {string} string "Room not found"
// @Router /rooms/{room_id} [get]
func GetRoomDetails(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	roomID, err := strconv.Atoi(params["room_id"])
	if err != nil {
		http.Error(w, "Invalid room id", http.StatusBadRequest)
		return
	}

	var room models.Room
	if result := database.DB.Find(&room, roomID); result.Error != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(room)

}
