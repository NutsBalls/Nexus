package controllers

import (
	"net/http"

	"github.com/NutsBalls/Nexus/config"
	"github.com/NutsBalls/Nexus/models"
	"github.com/labstack/echo/v4"
)

// DocumentController - структура контроллера для документов
type DocumentController struct{}

// CreateDocument - создает новый документ
func (dc *DocumentController) CreateDocument(c echo.Context) error {
	document := new(models.Document)
	if err := c.Bind(document); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Неправильный ввод"})
	}

	if err := config.DB.Create(&document).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Не удалось создать документ"})
	}

	return c.JSON(http.StatusOK, document)
}

// GetDocuments - получает все документы
func (dc *DocumentController) GetDocuments(c echo.Context) error {
	var documents []models.Document
	if err := config.DB.Find(&documents).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Не удалось получить документы"})
	}

	return c.JSON(http.StatusOK, documents)
}

// UpdateDocument - обновляет документ по ID
func (dc *DocumentController) UpdateDocument(c echo.Context) error {
	id := c.Param("id")
	document := new(models.Document)

	if err := config.DB.First(&document, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Документ не найден"})
	}

	if err := c.Bind(document); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Неправильный ввод"})
	}

	config.DB.Save(&document)
	return c.JSON(http.StatusOK, document)
}

// DeleteDocument - удаляет документ по ID
func (dc *DocumentController) DeleteDocument(c echo.Context) error {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Document{}, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Не удалось удалить документ"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Документ успешно удален"})
}
