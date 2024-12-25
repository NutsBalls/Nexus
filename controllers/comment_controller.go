package controllers

import (
	"net/http"
	"strconv"

	"github.com/NutsBalls/Nexus/models"
	"github.com/NutsBalls/Nexus/utils"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CommentController struct {
	DB *gorm.DB
}

func NewCommentController(db *gorm.DB) *CommentController {
	return &CommentController{DB: db}
}

func (cc *CommentController) AddComment(c echo.Context) error {
	documentIDStr := c.Param("id")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid document ID"})
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.JWTCustomClaims)

	comment := new(models.Comment)
	if err := c.Bind(comment); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	comment.DocumentID = uint(documentID)
	comment.UserID = claims.ID

	if err := cc.DB.Create(comment).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add comment"})
	}

	return c.JSON(http.StatusCreated, comment)
}

func (cc *CommentController) GetComments(c echo.Context) error {
	documentIDStr := c.Param("id")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid document ID"})
	}

	var comments []struct {
		models.Comment
		Username string `json:"username"`
	}

	if err := cc.DB.Table("comments").
		Select("comments.*, users.username").
		Joins("JOIN users ON users.id = comments.user_id").
		Where("comments.document_id = ?", documentID).
		Order("comments.created_at DESC").
		Find(&comments).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch comments"})
	}

	return c.JSON(http.StatusOK, comments)
}

func (cc *CommentController) DeleteComment(c echo.Context) error {
	commentIDStr := c.Param("commentId")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid comment ID"})
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.JWTCustomClaims)

	result := cc.DB.Where("id = ? AND user_id = ?", commentID, claims.ID).Delete(&models.Comment{})
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete comment"})
	}
	if result.RowsAffected == 0 {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Not authorized to delete this comment"})
	}

	return c.NoContent(http.StatusNoContent)
}
