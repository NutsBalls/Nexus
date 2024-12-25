package controllers

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/NutsBalls/Nexus/models"
	"github.com/NutsBalls/Nexus/utils"
	"github.com/golang-jwt/jwt"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type DocumentController struct {
	DB *gorm.DB
}

type CreateDocumentRequest struct {
	Title    string `json:"title" example:"Мой документ"`
	Content  string `json:"content" example:"Содержимое документа"`
	FolderID *uint  `json:"folder_id,omitempty" example:"1"`
	IsPublic bool   `json:"is_public" example:"false"`
}

func NewDocumentController(db *gorm.DB) *DocumentController {
	return &DocumentController{DB: db}
}

func (dc *DocumentController) GetDocuments(c echo.Context) error {
	var documents []models.Document
	if err := dc.DB.Find(&documents).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch documents"})
	}
	return c.JSON(http.StatusOK, documents)
}

func (dc *DocumentController) CreateDocument(c echo.Context) error {
	type CreateDocumentRequest struct {
		Title    string `json:"title" binding:"required"`
		Content  string `json:"content"`
		FolderID *uint  `json:"folder_id,omitempty"`
		IsPublic bool   `json:"is_public"`
	}

	req := new(CreateDocumentRequest)
	if err := c.Bind(req); err != nil {
		log.Printf("Failed to bind data: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload",
		})
	}

	claims := c.Get("claims").(*utils.JWTCustomClaims)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "User not authenticated",
		})
	}

	document := &models.Document{
		Title:    req.Title,
		Content:  req.Content,
		UserID:   claims.ID,
		FolderID: req.FolderID,
		IsPublic: req.IsPublic,
	}

	if document.FolderID != nil {
		var folder models.Folder
		if err := dc.DB.First(&folder, *document.FolderID).Error; err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Folder not found",
			})
		}
	}

	if err := dc.DB.Create(document).Error; err != nil {
		log.Printf("Failed to create document: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create document",
		})
	}

	return c.JSON(http.StatusCreated, document)
}

func (dc *DocumentController) GetDocument(c echo.Context) error {
	id := c.Param("id")
	document := new(models.Document)
	if err := dc.DB.First(document, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Document not found"})
	}
	return c.JSON(http.StatusOK, document)
}

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

func (dc *DocumentController) DeleteDocument(c echo.Context) error {
	documentID := c.Param("id")
	log.Printf("Attempting to delete document with ID: %s", documentID)

	var document models.Document
	if err := dc.DB.First(&document, documentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Document not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.JWTCustomClaims)
	if document.UserID != claims.ID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied"})
	}

	tx := dc.DB.Begin()

	var attachments []models.Attachment
	if err := tx.Where("document_id = ?", documentID).Find(&attachments).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch attachments"})
	}

	for _, attachment := range attachments {
		filePath := filepath.Join("uploads", attachment.Path)
		if err := os.Remove(filePath); err != nil {
			log.Printf("Warning: Failed to delete file: %v", err)
		}
	}

	if err := tx.Where("document_id = ?", documentID).Delete(&models.Attachment{}).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete attachments"})
	}

	if err := tx.Delete(&document).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete document"})
	}

	if err := tx.Commit().Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to commit transaction"})
	}

	return c.NoContent(http.StatusNoContent)
}

func (dc *DocumentController) CreateVersion(c echo.Context) error {
	documentID := c.Param("id")
	version := new(models.Version)
	if err := c.Bind(version); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
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

func (dc *DocumentController) GetVersions(c echo.Context) error {
	documentID := c.Param("id")

	var versions []models.Version
	if err := dc.DB.Where("document_id = ?", documentID).Order("created_at desc").Find(&versions).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch versions"})
	}

	return c.JSON(http.StatusOK, versions)
}

func (dc *DocumentController) SearchDocuments(c echo.Context) error {
	query := c.QueryParam("q")
	if query == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Search query is required"})
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.JWTCustomClaims)

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

func (dc *DocumentController) UploadAttachment(c echo.Context) error {
	documentIDStr := c.Param("id")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid document ID"})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No file uploaded"})
	}

	if err := os.MkdirAll("uploads", 0755); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create uploads directory"})
	}

	filePath := filepath.Join("uploads", file.Filename)

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to open uploaded file"})
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create destination file"})
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save file"})
	}

	attachment := models.Attachment{
		DocumentID: uint(documentID),
		Filename:   file.Filename,
		Path:       file.Filename,
		Size:       file.Size,
	}

	if err := dc.DB.Create(&attachment).Error; err != nil {
		os.Remove(filePath)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save attachment info"})
	}

	return c.JSON(http.StatusOK, attachment)
}

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

func (dc *DocumentController) DownloadAttachment(c echo.Context) error {
	path := c.Param("*")
	if path == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid path"})
	}

	filePath := filepath.Join("uploads", path)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "File not found"})
	}

	var attachment models.Attachment
	if err := dc.DB.Where("path = ?", path).First(&attachment).Error; err != nil {
		return c.File(filePath)
	}

	return c.Attachment(filePath, attachment.Filename)
}

func (dc *DocumentController) DeleteAttachment(c echo.Context) error {
	attachmentID := c.Param("id")

	var attachment models.Attachment
	if err := dc.DB.First(&attachment, attachmentID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Attachment not found"})
	}

	var document models.Document
	if err := dc.DB.First(&document, attachment.DocumentID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Document not found"})
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.JWTCustomClaims)
	if document.UserID != claims.ID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied"})
	}

	filePath := filepath.Join("uploads", attachment.Path)
	if err := os.Remove(filePath); err != nil {
		log.Printf("Failed to delete file: %v", err)
	}

	if err := dc.DB.Delete(&attachment).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete attachment"})
	}

	return c.NoContent(http.StatusNoContent)
}
