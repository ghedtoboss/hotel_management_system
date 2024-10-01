package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	Email     string `gorm:"unique; not null"`
	Role      string `gorm:"not null"` //"admin", "receptionist", "customer"
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Claims struct {
	Username string `json:"username"`
	UserID   uint   `json:"user_id"`
	Role     string `json:"role"`
	jwt.StandardClaims
}
