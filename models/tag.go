package models

import (
	"time"
)

type Tag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `gorm:"index" json:"deleted_at,omitempty"`

	Name      string     `gorm:"uniqueIndex;not null" json:"name" example:"Important"`
	UserID    uint       `json:"user_id" example:"1"`
	User      User       `gorm:"foreignKey:UserID" json:"-"`
	Documents []Document `gorm:"many2many:document_tags;" json:"documents,omitempty"`
}
