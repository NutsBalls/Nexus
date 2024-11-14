package controllers

import (
	"github.com/NutsBalls/Nexus/config"
	"github.com/NutsBalls/Nexus/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

// CreateDocument godoc
// @Summary Создает новый документ
// @Description Добавляет новый документ в базу данных
// @Tags documents
// @Accept json
// @Produce json
// @Param document body models.Document true "Document data"
// @Success 200 {object} models.Document
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/documents [post]
func CreateDocument(c echo.Context) error {
	document := new(models.Document)
	if err := c.Bind(document); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if err := config.DB.Create(&document).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not create document"})
	}

	return c.JSON(http.StatusOK, document)
}

// GetDocuments godoc
// @Summary Получает список документов
// @Description Возвращает все документы из базы данных
// @Tags documents
// @Produce json
// @Success 200 {array} models.Document
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/documents [get]
func GetDocuments(c echo.Context) error {
	var documents []models.Document
	if err := config.DB.Find(&documents).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not fetch documents"})
	}

	return c.JSON(http.StatusOK, documents)
}

// UpdateDocument godoc
// @Summary Обновляет документ
// @Description Изменяет данные указанного документа по ID
// @Tags documents
// @Accept json
// @Produce json
// @Param id path int true "Document ID"
// @Param document body models.Document true "Updated document data"
// @Success 200 {object} models.Document
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 404 {object} map[string]string "Not Found"
// @Router /api/documents/{id} [put]
func UpdateDocument(c echo.Context) error {
	id := c.Param("id")
	document := new(models.Document)

	if err := config.DB.First(&document, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Document not found"})
	}

	if err := c.Bind(document); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	config.DB.Save(&document)
	return c.JSON(http.StatusOK, document)
}

// DeleteDocument godoc
// @Summary Удаляет документ
// @Description Удаляет документ по ID
// @Tags documents
// @Produce json
// @Param id path int true "Document ID"
// @Success 200 {object} map[string]string "Document deleted successfully"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/documents/{id} [delete]
func DeleteDocument(c echo.Context) error {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Document{}, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not delete document"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Document deleted successfully"})
}
