package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/NutsBalls/Nexus/services"
	"github.com/NutsBalls/Nexus/utils"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type ExportController struct {
	exportService *services.ExportService
}

func NewExportController(exportService *services.ExportService) *ExportController {
	return &ExportController{exportService: exportService}
}

func (ec *ExportController) ExportDocument(c echo.Context) error {
	documentIDStr := c.Param("id")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid document ID"})
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.JWTCustomClaims)

	filename, err := ec.exportService.ExportDocumentToJSON(uint(documentID), claims.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to export document"})
	}

	return c.Attachment(filename, fmt.Sprintf("document_%d.json", documentID))
}
