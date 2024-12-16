package services

import (
	"github.com/NutsBalls/Nexus/models"

	"gorm.io/gorm"
)

type DocumentService struct {
	db *gorm.DB
}

func NewDocumentService(db *gorm.DB) *DocumentService {
	return &DocumentService{db: db}
}

func (s *DocumentService) CreateDocumentWithTags(document *models.Document, tagNames []string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Создаем документ
		if err := tx.Create(document).Error; err != nil {
			return err
		}

		// Обрабатываем теги
		for _, tagName := range tagNames {
			var tag models.Tag
			// Ищем существующий тег или создаем новый
			err := tx.Where("name = ? AND user_id = ?", tagName, document.UserID).
				FirstOrCreate(&tag, models.Tag{
					Name:   tagName,
					UserID: document.UserID,
				}).Error
			if err != nil {
				return err
			}

			// Связываем тег с документом
			if err := tx.Model(document).Association("Tags").Append(&tag); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *DocumentService) GetDocumentWithDetails(id uint, userID uint) (*models.Document, error) {
	var document models.Document
	err := s.db.Preload("Tags").
		Preload("Versions").
		Preload("SharedUsers").
		Where("id = ? AND (user_id = ? OR id IN (SELECT document_id FROM document_shares WHERE user_id = ?))",
			id, userID, userID).
		First(&document).Error
	return &document, err
}
