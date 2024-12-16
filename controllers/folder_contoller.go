package controllers

import (
	"net/http"

	middleware "github.com/NutsBalls/Nexus/middlewares"
	"github.com/NutsBalls/Nexus/models"
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
	folder := new(models.Folder)
	if err := c.Bind(folder); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Получаем пользователя из JWT
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*middleware.JWTCustomClaims)
	folder.UserID = claims.ID

	if err := fc.DB.Create(folder).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create folder"})
	}

	return c.JSON(http.StatusCreated, folder)
}

func (fc *FolderController) GetFolders(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*middleware.JWTCustomClaims)

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
	id := c.Param("id")

	// Проверяем существование папки и права доступа
	var folder models.Folder
	if err := fc.DB.First(&folder, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Folder not found"})
	}

	// Удаляем папку
	if err := fc.DB.Delete(&folder).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete folder"})
	}

	return c.NoContent(http.StatusNoContent)
}
