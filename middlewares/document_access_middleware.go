package middlewares

import (
	"net/http"
	"strconv"

	"github.com/NutsBalls/Nexus/models"
	"github.com/NutsBalls/Nexus/utils"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func DocumentAccessMiddleware(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			documentIDStr := c.Param("id")
			documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid document ID"})
			}

			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(*utils.JWTCustomClaims)

			var document models.Document
			if err := db.First(&document, documentID).Error; err != nil {
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Document not found"})
			}

			if document.UserID == claims.ID {
				return next(c)
			}

			var collaborationExists bool
			err = db.Raw(`
                                SELECT EXISTS(
                                        SELECT 1 FROM collaborations 
                                        WHERE document_id = ? AND user_id = ?
                                )`, documentID, claims.ID).Scan(&collaborationExists).Error

			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Access check failed"})
			}

			if collaborationExists {
				return next(c)
			}

			return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied"})
		}
	}
}
