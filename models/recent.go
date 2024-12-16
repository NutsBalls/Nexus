package models

import (
	"time"

	"gorm.io/gorm"
)

type RecentDocument struct {
	gorm.Model
	UserID     uint      `json:"user_id"`
	DocumentID uint      `json:"document_id"`
	LastAccess time.Time `json:"last_access"`
}
