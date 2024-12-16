package models

import (
	"gorm.io/gorm"
)

type Favorite struct {
	gorm.Model
	UserID     uint `json:"user_id"`
	DocumentID uint `json:"document_id"`
}
