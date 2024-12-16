package controllers

import (
	"net/http"
	"path/filepath"
	"strconv"

	middleware "github.com/NutsBalls/Nexus/middlewares"
	"github.com/NutsBalls/Nexus/models"
	"github.com/NutsBalls/Nexus/utils"
	"github.com/golang-jwt/jwt"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ShareRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type DocumentController struct {
	DB *gorm.DB
}

func NewDocumentController(db *gorm.DB) *DocumentController {
	return &DocumentController{DB: db}
}

// GetDocuments godoc
// @Summary Получить список всех документов
// @Description Возвращает список всех документов
// @Tags documents
// @Produce json
// @Success 200 {array} models.Document
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /api/documents [get]
func (dc *DocumentController) GetDocuments(c echo.Context) error {
	var documents []models.Document
	if err := dc.DB.Find(&documents).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch documents"})
	}
	return c.JSON(http.StatusOK, documents)
}

// CreateDocument godoc
// @Summary Создать новый документ
// @Description Создает новый документ с переданными данными
// @Tags documents
// @Accept json
// @Produce json
// @Param document body models.Document true "Данные документа"
// @Success 201 {object} models.Document
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/documents [post]
// @Security ApiKeyAuth
func (dc *DocumentController) CreateDocument(c echo.Context) error {
	document := new(models.Document)
	if err := c.Bind(document); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	if err := dc.DB.Create(document).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create document"})
	}

	return c.JSON(http.StatusCreated, document)
}

// GetDocument godoc
// @Summary Получить документ по ID
// @Description Возвращает документ по указанному идентификатору
// @Tags documents
// @Produce json
// @Param id path string true "ID документа"
// @Success 200 {object} models.Document
// @Failure 404 {object} map[string]string
// @Router /api/documents/{id} [get]
// @Security ApiKeyAuth
func (dc *DocumentController) GetDocument(c echo.Context) error {
	id := c.Param("id")
	document := new(models.Document)
	if err := dc.DB.First(document, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Document not found"})
	}
	return c.JSON(http.StatusOK, document)
}

// UpdateDocument godoc
// @Summary Обновить документ
// @Description Обновляет существующий документ
// @Tags documents
// @Accept json
// @Produce json
// @Param id path string true "ID документа"
// @Param document body models.Document true "Обновленные данные документа"
// @Success 200 {object} models.Document
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/documents/{id} [put]
func (dc *DocumentController) UpdateDocument(c echo.Context) error {
	id := c.Param("id")
	document := new(models.Document)
	if err := dc.DB.First(document, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Document not found"})
	}

	if err := c.Bind(document); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	if err := dc.DB.Save(document).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update document"})
	}

	return c.JSON(http.StatusOK, document)
}

// DeleteDocument godoc
// @Summary Удалить документ
// @Description Удаляет документ по указанному идентификатору
// @Tags documents
// @Param id path string true "ID документа"
// @Success 204
// @Failure 500 {object} map[string]string
// @Router /api/documents/{id} [delete]
func (dc *DocumentController) DeleteDocument(c echo.Context) error {
	id := c.Param("id")
	if err := dc.DB.Delete(&models.Document{}, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete document"})
	}
	return c.NoContent(http.StatusNoContent)
}

// CreateVersion godoc
// @Summary Создать новую версию документа
// @Description Создает новую версию для указанного документа
// @Tags document versions
// @Accept json
// @Produce json
// @Param id path string true "ID документа"
// @Param version body models.Version true "Данные версии"
// @Success 201 {object} models.Version
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/documents/{id}/versions [post]
func (dc *DocumentController) CreateVersion(c echo.Context) error {
	documentID := c.Param("id")
	version := new(models.Version)
	if err := c.Bind(version); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Проверяем существование документа
	var document models.Document
	if err := dc.DB.First(&document, documentID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Document not found"})
	}

	version.DocumentID = document.ID
	if err := dc.DB.Create(version).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create version"})
	}

	return c.JSON(http.StatusCreated, version)
}

// GetVersions godoc
// @Summary Получить версии документа
// @Description Возвращает список версий для указанного документа
// @Tags document versions
// @Produce json
// @Param id path string true "ID документа"
// @Success 200 {array} models.Version
// @Failure 500 {object} map[string]string
// @Router /api/documents/{id}/versions [get]
func (dc *DocumentController) GetVersions(c echo.Context) error {
	documentID := c.Param("id")

	var versions []models.Version
	if err := dc.DB.Where("document_id = ?", documentID).Order("created_at desc").Find(&versions).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch versions"})
	}

	return c.JSON(http.StatusOK, versions)
}

// ShareDocument godoc
// @Summary Предоставить доступ к документу
// @Description Предоставляет доступ к документу другому пользователю по email
// @Tags document sharing
// @Accept json
// @Produce json
// @Param id path string true "ID документа"
// @Param share body ShareRequest true "Данные для общего доступа (email пользователя)"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/documents/{id}/share [post]
func (dc *DocumentController) ShareDocument(c echo.Context) error {
	type ShareRequest struct {
		Email string `json:"email" validate:"required,email"`
	}

	documentID := c.Param("id")
	req := new(ShareRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Находим документ
	var document models.Document
	if err := dc.DB.First(&document, documentID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Document not found"})
	}

	// Находим пользователя для совместного доступа
	var shareUser models.User
	if err := dc.DB.Where("email = ?", req.Email).First(&shareUser).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	// Добавляем пользователя в список общего доступа
	if err := dc.DB.Model(&document).Association("SharedUsers").Append(&shareUser); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to share document"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Document shared successfully"})
}

// SearchDocuments godoc
// @Summary Поиск документов
// @Description Поиск документов по запросу с учетом прав доступа
// @Tags documents
// @Produce json
// @Param q query string true "Поисковый запрос"
// @Success 200 {array} models.Document
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/documents/search [get]
func (dc *DocumentController) SearchDocuments(c echo.Context) error {
	query := c.QueryParam("q")
	if query == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Search query is required"})
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*middleware.JWTCustomClaims)

	var documents []models.Document
	if err := dc.DB.Where("user_id = ? AND (title ILIKE ? OR content ILIKE ?)",
		claims.ID, "%"+query+"%", "%"+query+"%").
		Or("id IN (SELECT document_id FROM document_shares WHERE user_id = ?)", claims.ID).
		Preload("Tags").
		Find(&documents).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to search documents"})
	}

	return c.JSON(http.StatusOK, documents)
}

// UploadAttachment godoc
// @Summary Загрузить вложение
// @Description Загружает вложение для указанного документа
// @Tags document attachments
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "ID документа"
// @Param file formData file true "Файл вложения"
// @Success 200 {object} models.Attachment
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/documents/{id}/attachments [post]
func (dc *DocumentController) UploadAttachment(c echo.Context) error {
	documentIDStr := c.Param("id")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid document ID"})
	}

	// Получаем файл из запроса
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No file uploaded"})
	}

	// Генерируем путь для сохранения файла
	filename := filepath.Join("uploads", documentIDStr, file.Filename)

	// Сохраняем файл
	if err := utils.SaveFile(file, filename); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save file"})
	}

	// Создаем запись о вложении в базе данных
	attachment := models.Attachment{
		DocumentID: uint(documentID),
		Filename:   file.Filename,
		Path:       filename,
		Size:       file.Size,
	}

	if err := dc.DB.Create(&attachment).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save attachment info"})
	}

	return c.JSON(http.StatusOK, attachment)
}

// GetAttachments godoc
// @Summary Получить вложения документа
// @Description Возвращает список вложений для указанного документа
// @Tags document attachments
// @Produce json
// @Param id path string true "ID документа"
// @Success 200 {array} models.Attachment
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/documents/{id}/attachments [get]
func (dc *DocumentController) GetAttachments(c echo.Context) error {
	documentIDStr := c.Param("id")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid document ID"})
	}

	var attachments []models.Attachment
	if err := dc.DB.Where("document_id = ?", documentID).Find(&attachments).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch attachments"})
	}

	return c.JSON(http.StatusOK, attachments)
}
