package models

import "time"

type Document struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"userID"`
	FolderID  uint      `json:"folderID"`
	Title     string    `json:"title" validate:"required"`
	Content   string    `json:"content" validate:"required"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
