package services

import (
	"encoding/json"
	"io"

	"github.com/NutsBalls/Nexus/models"

	"gorm.io/gorm"
)

type ImportService struct {
	db *gorm.DB
}

func NewImportService(db *gorm.DB) *ImportService {
	return &ImportService{db: db}
}

func (is *ImportService) ImportDocumentFromJSON(reader io.Reader, userID uint) (*models.Document, error) {
	var document models.Document
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&document); err != nil {
		return nil, err
	}

	document.ID = 0
	document.UserID = userID

	err := is.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&document).Error; err != nil {
			return err
		}

		// Импортируем теги
		for _, tag := range document.Tags {
			tag.ID = 0
			tag.UserID = userID
			if err := tx.Create(&tag).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &document, nil
}
