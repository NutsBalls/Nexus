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

type ShareController struct {
	DB *gorm.DB
}

func NewShareController(db *gorm.DB) *ShareController {
	return &ShareController{DB: db}
}

func (sc *ShareController) ShareDocument(c echo.Context) error {
	documentID := c.Param("id")

	type ShareRequest struct {
		UserEmail  string                 `json:"user_email"`
		Permission models.SharePermission `json:"permission"`
	}

	var req ShareRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.JWTCustomClaims)

	var document models.Document
	if err := sc.DB.First(&document, documentID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Document not found"})
	}

	var existingShare models.Share
	err := sc.DB.Where("document_id = ? AND user_id = ? AND permission = ?",
		documentID, claims.ID, models.PermissionAdmin).First(&existingShare).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to check access rights"})
	}

	isAdmin := err == nil
	if document.UserID != claims.ID && !isAdmin {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied"})
	}

	if req.UserEmail == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "User email cannot be empty"})
	}

	var targetUser models.User
	if err := sc.DB.Where("email = ?", req.UserEmail).First(&targetUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to find user"})
	}

	share := models.Share{
		DocumentID:  document.ID,
		UserID:      targetUser.ID,
		Permission:  req.Permission,
		CreatedByID: claims.ID,
	}

	if err := sc.DB.Where("document_id = ? AND user_id = ?", document.ID, targetUser.ID).First(&models.Share{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if createErr := sc.DB.Create(&share).Error; createErr != nil {
				log.Printf("Ошибка при создании share: %v", createErr)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to share document"})
			}
			return c.JSON(http.StatusOK, share)
		} else {
			log.Printf("Ошибка при поиске share: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to check existing share"})
		}
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Share already exists for this user and document"})
	}
}

func (sc *ShareController) GetDocumentShares(c echo.Context) error {
	documentID := c.Param("id")

	var shares []models.Share
	if err := sc.DB.Preload("User").Where("document_id = ?", documentID).Find(&shares).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch shares"})
	}

	return c.JSON(http.StatusOK, shares)
}

func (sc *ShareController) RemoveShare(c echo.Context) error {
	shareID := c.Param("id")

	var share models.Share
	if err := sc.DB.First(&share, shareID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Share not found"})
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.JWTCustomClaims)

	var document models.Document
	if err := sc.DB.First(&document, share.DocumentID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Document not found"})
	}

	if document.UserID != claims.ID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied"})
	}

	if err := sc.DB.Where("document_id = ? AND user_id = ?", share.DocumentID, share.UserID).Delete(&models.Share{}).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to remove share"})
	}

	return c.NoContent(http.StatusNoContent)
}

func (sc *ShareController) GetSharedWithMe(c echo.Context) error {
	userToken, ok := c.Get("user").(*jwt.Token)
	if !ok || userToken == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}
	claims, ok := userToken.Claims.(*utils.JWTCustomClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token claims"})
	}

	userID := claims.ID
	var shares []models.Share
	if err := sc.DB.Preload("Document").Where("user_id = ?", userID).Find(&shares).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch shared documents"})
	}

	var sharedDocuments []models.Document
	for _, share := range shares {
		sharedDocuments = append(sharedDocuments, share.Document)
	}

	return c.JSON(http.StatusOK, sharedDocuments)
}

func (sc *ShareController) GetSharedByMe(c echo.Context) error {
	userToken, ok := c.Get("user").(*jwt.Token)
	if !ok || userToken == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}
	claims, ok := userToken.Claims.(*utils.JWTCustomClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token claims"})
	}

	var shares []models.Share
	if err := sc.DB.Preload("Document").Preload("User").Where("created_by_id = ?", claims.ID).Find(&shares).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch shared documents"})
	}

	return c.JSON(http.StatusOK, shares)
}

func (sc *ShareController) CheckDocumentAccess(c echo.Context) error {
	documentID := c.Param("id")

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.JWTCustomClaims)

	var document models.Document
	if err := sc.DB.First(&document, documentID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Document not found"})
	}

	isOwner := document.UserID == claims.ID

	var share models.Share
	err := sc.DB.Where("document_id = ? AND user_id = ?", documentID, claims.ID).First(&share).Error

	hasAccess := err == nil

	return c.JSON(http.StatusOK, map[string]bool{
		"isOwner":   isOwner,
		"hasAccess": hasAccess,
	})
}
