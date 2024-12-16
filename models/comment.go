package models

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	DocumentID uint   `json:"document_id"`
	UserID     uint   `json:"user_id"`
	Content    string `json:"content"`
	ParentID   *uint  `json:"parent_id,omitempty"`
}
