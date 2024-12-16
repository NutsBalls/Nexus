package models

import (
	"gorm.io/gorm"
)

type NotificationType string

const (
	NotificationShare   NotificationType = "share"
	NotificationComment NotificationType = "comment"
	NotificationMention NotificationType = "mention"
)

type Notification struct {
	gorm.Model
	UserID     uint             `json:"user_id"`
	Type       NotificationType `json:"type"`
	DocumentID uint             `json:"document_id"`
	SenderID   uint             `json:"sender_id"`
	Content    string           `json:"content"`
	IsRead     bool             `json:"is_read" gorm:"default:false"`
}
