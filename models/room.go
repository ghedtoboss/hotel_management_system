package models

import (
	"time"

	"gorm.io/gorm"
)

type Room struct {
	ID        uint    `gorm:"primaryKey"`
	Number    string  `gorm:"unique;not null"`
	Type      string  `gorm:"not null"` //"single", "double", "suite"
	Status    string  `gorm:"not null"` //"available", "occupied", "cleaning"
	Price     float64 `gorm:"not null"`
	CreatedAt time.Time
	UpdateAt  time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
