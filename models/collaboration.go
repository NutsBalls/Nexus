package models

import (
	"time"

	"gorm.io/gorm"
)

type CollaborationRole string

const (
	RoleViewer CollaborationRole = "viewer"
	RoleEditor CollaborationRole = "editor"
	RoleAdmin  CollaborationRole = "admin"
)

type Collaboration struct {
	gorm.Model
	DocumentID uint              `json:"document_id"`
	UserID     uint              `json:"user_id"`
	Role       CollaborationRole `json:"role"`
	ExpiresAt  *time.Time        `json:"expires_at,omitempty"`
}
