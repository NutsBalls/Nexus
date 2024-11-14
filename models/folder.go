package models

import "time"

// Folder represents a folder for organizing documents
type Folder struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"foreignKey:UserID"`
	Name      string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
