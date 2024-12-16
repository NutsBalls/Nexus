package services

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/NutsBalls/Nexus/models"
	"gorm.io/gorm"
)

type ExportService struct {
	db *gorm.DB
}

func NewExportService(db *gorm.DB) *ExportService {
	return &ExportService{db: db}
}

func (es *ExportService) ExportDocumentToJSON(documentID uint, userID uint) (string, error) {
	var document models.Document
	if err := es.db.Preload("Tags").
		Preload("Versions").
		Where("id = ? AND (user_id = ? OR id IN (SELECT document_id FROM document_shares WHERE user_id = ?))",
			documentID, userID, userID).
		First(&document).Error; err != nil {
		return "", err
	}

	filename := fmt.Sprintf("exports/document_%d.json", documentID)

	// Создаем директорию если не существует
	os.MkdirAll("exports", 0755)

	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(document); err != nil {
		return "", err
	}

	return filename, nil
}
