package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `gorm:"index" json:"deleted_at,omitempty"`

	Username  string     `gorm:"uniqueIndex;not null" json:"username" example:"johndoe"`
	Email     string     `gorm:"uniqueIndex;not null" json:"email" example:"john@example.com"`
	Password  string     `gorm:"not null" json:"-"`
	Documents []Document `gorm:"foreignKey:UserID" json:"documents,omitempty"`
	Folders   []Folder   `gorm:"foreignKey:UserID" json:"folders,omitempty"`
}
