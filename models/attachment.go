package models

import "time"

type Attachment struct {
	ID         uint       `json:"id" gorm:"primaryKey"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty" gorm:"index"`
	DocumentID uint       `json:"document_id"`
	Filename   string     `json:"filename"`
	Path       string     `json:"path"`
	Size       int64      `json:"size"`
}
