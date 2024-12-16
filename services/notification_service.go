package services

import (
	"github.com/NutsBalls/Nexus/models"

	"gorm.io/gorm"
)

type NotificationService struct {
	db *gorm.DB
}

func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{db: db}
}

func (ns *NotificationService) CreateNotification(userID uint, senderID uint, documentID uint, notificationType models.NotificationType, content string) error {
	notification := models.Notification{
		UserID:     userID,
		SenderID:   senderID,
		DocumentID: documentID,
		Type:       notificationType,
		Content:    content,
	}

	return ns.db.Create(&notification).Error
}

func (ns *NotificationService) NotifyCollaborators(documentID uint, senderID uint, notificationType models.NotificationType, content string) error {
	var collaborators []models.Collaboration
	if err := ns.db.Where("document_id = ? AND user_id != ?", documentID, senderID).Find(&collaborators).Error; err != nil {
		return err
	}

	for _, collaborator := range collaborators {
		if err := ns.CreateNotification(collaborator.UserID, senderID, documentID, notificationType, content); err != nil {
			return err
		}
	}

	return nil
}
