package models

import "time"

// Tag represents a tag that can be attached to documents
type Tag struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"unique"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
