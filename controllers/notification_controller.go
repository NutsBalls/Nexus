package controllers

import (
	"net/http"
	"strconv"

	"github.com/NutsBalls/Nexus/models"
	"github.com/NutsBalls/Nexus/utils" // Импортируем JWTCustomClaims

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type NotificationController struct {
	DB *gorm.DB
}

func NewNotificationController(db *gorm.DB) *NotificationController {
	return &NotificationController{DB: db}
}

func (nc *NotificationController) GetNotifications(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.JWTCustomClaims)

	var notifications []struct {
		models.Notification
		SenderName string `json:"sender_name"`
	}

	if err := nc.DB.Table("notifications").
		Select("notifications.*, users.username as sender_name").
		Joins("JOIN users ON users.id = notifications.sender_id").
		Where("notifications.user_id = ?", claims.ID).
		Order("notifications.created_at DESC").
		Find(&notifications).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch notifications"})
	}

	return c.JSON(http.StatusOK, notifications)
}

func (nc *NotificationController) MarkAsRead(c echo.Context) error {
	notificationIDStr := c.Param("id")
	notificationID, err := strconv.ParseUint(notificationIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid notification ID"})
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.JWTCustomClaims)

	result := nc.DB.Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", notificationID, claims.ID).
		Update("is_read", true)

	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to mark notification as read"})
	}

	if result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Notification not found or already read"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Notification marked as read"})
}

func (nc *NotificationController) MarkAllAsRead(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.JWTCustomClaims)

	if err := nc.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", claims.ID, false).
		Update("is_read", true).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to mark notifications as read"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "All unread notifications marked as read"})
}
