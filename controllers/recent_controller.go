package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/NutsBalls/Nexus/models"
	"github.com/NutsBalls/Nexus/utils" // Импортируем JWTCustomClaims

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type RecentController struct {
	DB *gorm.DB
}

func NewRecentController(db *gorm.DB) *RecentController {
	return &RecentController{DB: db}
}

func (rc *RecentController) AddToRecent(c echo.Context) error {
	documentIDStr := c.Param("id")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid document ID"})
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.JWTCustomClaims)

	var recent models.RecentDocument
	result := rc.DB.Where("user_id = ? AND document_id = ?", claims.ID, documentID).
		FirstOrCreate(&recent, models.RecentDocument{
			UserID:     claims.ID,
			DocumentID: uint(documentID),
		})

	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update recent documents"})
	}

	recent.LastAccess = time.Now()
	if err := rc.DB.Save(&recent).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update last access time"})
	}

	return c.JSON(http.StatusOK, recent)
}

func (rc *RecentController) GetRecentDocuments(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.JWTCustomClaims)

	var documents []models.Document
	if err := rc.DB.Joins("JOIN recent_documents ON recent_documents.document_id = documents.id").
		Where("recent_documents.user_id = ?", claims.ID).
		Order("recent_documents.last_access DESC").
		Limit(10).
		Preload("Tags").
		Find(&documents).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch recent documents"})
	}

	return c.JSON(http.StatusOK, documents)
}
