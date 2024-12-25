package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/NutsBalls/Nexus/models"
	"github.com/NutsBalls/Nexus/utils"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type FolderController struct {
	DB *gorm.DB
}

func NewFolderController(db *gorm.DB) *FolderController {
	return &FolderController{DB: db}
}

func (fc *FolderController) CreateFolder(c echo.Context) error {
	type CreateFolderRequest struct {
		Name string `json:"name" binding:"required"`
	}

	req := new(CreateFolderRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request data",
		})
	}

	claims := c.Get("claims").(*utils.JWTCustomClaims)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "User not authenticated",
		})
	}

	folder := &models.Folder{
		Name:   req.Name,
		UserID: claims.ID,
	}

	if err := fc.DB.Create(folder).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create folder",
		})
	}

	return c.JSON(http.StatusCreated, folder)
}

func (dc *DocumentController) GetFolderDocuments(c echo.Context) error {
	folderID := c.Param("id")

	var documents []models.Document
	if err := dc.DB.Where("folder_id = ?", folderID).Find(&documents).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch documents",
		})
	}

	return c.JSON(http.StatusOK, documents)
}

func (fc *FolderController) GetFolders(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.JWTCustomClaims)

	var folders []models.Folder
	if err := fc.DB.Where("user_id = ?", claims.ID).Find(&folders).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch folders"})
	}

	return c.JSON(http.StatusOK, folders)
}

func (fc *FolderController) UpdateFolder(c echo.Context) error {
	id := c.Param("id")
	folder := new(models.Folder)

	if err := fc.DB.First(folder, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Folder not found"})
	}

	if err := c.Bind(folder); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := fc.DB.Save(folder).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update folder"})
	}

	return c.JSON(http.StatusOK, folder)
}

func (fc *FolderController) DeleteFolder(c echo.Context) error {
	log.Printf("DeleteFolder called with context: %v", c.Path())

	folderID := c.Param("id")
	log.Printf("Attempting to delete folder with ID: %s", folderID)

	var folder models.Folder
	if err := fc.DB.First(&folder, folderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Folder not found with ID: %s", folderID)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Folder not found"})
		}
		log.Printf("Database error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.JWTCustomClaims)

	if folder.UserID != claims.ID {
		log.Printf("Access denied: folder user ID %d doesn't match token user ID %d", folder.UserID, claims.ID)
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied"})
	}

	if err := fc.DB.Unscoped().Delete(&folder).Error; err != nil {
		log.Printf("Error deleting folder: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete folder"})
	}

	log.Printf("Successfully deleted folder with ID: %s", folderID)
	return c.NoContent(http.StatusNoContent)
}
