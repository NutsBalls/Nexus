package controllers

import (
	"net/http"
	"strconv"

	"github.com/NutsBalls/Nexus/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CollaborationController struct {
	DB *gorm.DB
}

func NewCollaborationController(db *gorm.DB) *CollaborationController {
	return &CollaborationController{DB: db}
}

type AddCollaboratorRequest struct {
	Email string                   `json:"email" validate:"required,email"`
	Role  models.CollaborationRole `json:"role" validate:"required"`
}

func (cc *CollaborationController) AddCollaborator(c echo.Context) error {
	documentIDStr := c.Param("id")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid document ID"})
	}

	req := new(AddCollaboratorRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Находим пользователя по email
	var user models.User
	if err := cc.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	collaboration := models.Collaboration{
		DocumentID: uint(documentID),
		UserID:     user.ID,
		Role:       req.Role,
	}

	if err := cc.DB.Create(&collaboration).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add collaborator"})
	}

	return c.JSON(http.StatusOK, collaboration)
}

func (cc *CollaborationController) RemoveCollaborator(c echo.Context) error {
	documentIDStr := c.Param("id")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid document ID"})
	}

	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	if err := cc.DB.Where("document_id = ? AND user_id = ?", documentID, userID).
		Delete(&models.Collaboration{}).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to remove collaborator"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Collaborator removed successfully"})
}

func (cc *CollaborationController) GetCollaborators(c echo.Context) error {
	documentIDStr := c.Param("id")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid document ID"})
	}

	var collaborations []struct {
		models.Collaboration
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	if err := cc.DB.Table("collaborations").
		Select("collaborations.*, users.username, users.email").
		Joins("JOIN users ON users.id = collaborations.user_id").
		Where("collaborations.document_id = ?", documentID).
		Find(&collaborations).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch collaborators"})
	}

	return c.JSON(http.StatusOK, collaborations)
}
