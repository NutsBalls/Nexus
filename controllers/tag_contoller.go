package controllers

import (
	"net/http"

	"github.com/NutsBalls/Nexus/models"
	"github.com/NutsBalls/Nexus/utils"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type TagController struct {
	DB *gorm.DB
}

func NewTagController(db *gorm.DB) *TagController {
	return &TagController{DB: db}
}

func (tc *TagController) CreateTag(c echo.Context) error {
	tag := new(models.Tag)
	if err := c.Bind(tag); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.JWTCustomClaims)
	tag.UserID = claims.ID

	if err := tc.DB.Create(tag).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create tag"})
	}

	return c.JSON(http.StatusCreated, tag)
}

func (tc *TagController) GetTags(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.JWTCustomClaims)

	var tags []models.Tag
	if err := tc.DB.Where("user_id = ?", claims.ID).Find(&tags).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch tags"})
	}

	return c.JSON(http.StatusOK, tags)
}

func (tc *TagController) SearchByTag(c echo.Context) error {
	tagName := c.QueryParam("tag")
	if tagName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Tag name is required"})
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.JWTCustomClaims)

	var documents []models.Document
	if err := tc.DB.Joins("JOIN document_tags ON document_tags.document_id = documents.id").
		Joins("JOIN tags ON tags.id = document_tags.tag_id").
		Where("tags.name = ? AND (documents.user_id = ? OR documents.id IN (SELECT document_id FROM document_shares WHERE user_id = ?))",
			tagName, claims.ID, claims.ID).
		Preload("Tags").
		Find(&documents).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to search documents by tag"})
	}

	return c.JSON(http.StatusOK, documents)
}
