package models

import (
	"time"
)

type Reservation struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	UserID    uint `gorm:"int" json:"user_id"`
	RoomID    uint `gorm:"int" json:"room_id"`
	StartDate time.Time
	EndDate   time.Time
	Status    string `gorm:"string" json:"status"` //pending, confirmed, checked-in, checked-out, cancelled, no-show
	CreatedAt time.Time
	UpdatedAt time.Time
}
