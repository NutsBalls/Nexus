package models

import "time"

// Document represents a document belonging to a user
type Document struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"foreignKey:UserID"`
	FolderID  uint      `gorm:"foreignKey:FolderID"`
	Title     string    `gorm:"not null"`
	Content   string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
