package controllers

import (
	"encoding/json"
	"net/http"
	"time"
)

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

}
