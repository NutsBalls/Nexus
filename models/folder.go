package models

import (
	"time"
)

type Folder struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Name      string     `json:"name" binding:"required"`
	UserID    uint       `json:"user_id"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}
