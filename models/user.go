package models

import "time"

// User represents a user in the system
type User struct {
	ID           uint      `json:"id" example:"1" gorm:"primaryKey"`
	Username     string    `json:"username" example:"testuser" gorm:"unique"`
	Email        string    `json:"email" example:"test@example.com" gorm:"unique"`
	PasswordHash string    `json:"password_hash" example:"$2a$10$..." gorm:"not null"`
	IsAdmin      bool      `json:"is_admin" example:"false" gorm:"default:false"`
	CreatedAt    time.Time `json:"created_at" example:"2023-01-01T00:00:00Z" gorm:"autoCreateTime"`
}
