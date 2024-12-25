package models

import (
	"time"
)

type SharePermission string

const (
	PermissionRead  SharePermission = "read"
	PermissionWrite SharePermission = "write"
	PermissionAdmin SharePermission = "admin"
)

type Share struct {
	ID          uint            `json:"id" gorm:"primaryKey"`
	DocumentID  uint            `json:"document_id"`
	Document    Document        `json:"-" gorm:"foreignKey:DocumentID"`
	UserID      uint            `json:"user_id"`
	User        User            `json:"user" gorm:"foreignKey:UserID"`
	Permission  SharePermission `json:"permission"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	CreatedByID uint            `json:"created_by_id"`
	CreatedBy   User            `json:"created_by" gorm:"foreignKey:CreatedByID"`
}
