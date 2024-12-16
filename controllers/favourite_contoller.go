package controllers

import (
	"net/http"
	"strconv"

	middleware "github.com/NutsBalls/Nexus/middlewares"
	"github.com/NutsBalls/Nexus/models"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type FavoriteController struct {
	DB *gorm.DB
}

func NewFavoriteController(db *gorm.DB) *FavoriteController {
	return &FavoriteController{DB: db}
}

func (fc *FavoriteController) AddToFavorites(c echo.Context) error {
	documentIDStr := c.Param("id")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid document ID"})
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*middleware.JWTCustomClaims)

	favorite := models.Favorite{
		UserID:     claims.ID,
		DocumentID: uint(documentID),
	}

	if err := fc.DB.Create(&favorite).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add to favorites"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Added to favorites"})
}

func (fc *FavoriteController) RemoveFromFavorites(c echo.Context) error {
	documentIDStr := c.Param("id")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid document ID"})
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*middleware.JWTCustomClaims)

	if err := fc.DB.Where("user_id = ? AND document_id = ?", claims.ID, documentID).
		Delete(&models.Favorite{}).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to remove from favorites"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Removed from favorites"})
}

func (fc *FavoriteController) GetFavorites(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*middleware.JWTCustomClaims)

	var documents []models.Document
	if err := fc.DB.Joins("JOIN favorites ON favorites.document_id = documents.id").
		Where("favorites.user_id = ?", claims.ID).
		Preload("Tags").
		Find(&documents).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch favorites"})
	}

	return c.JSON(http.StatusOK, documents)
}
