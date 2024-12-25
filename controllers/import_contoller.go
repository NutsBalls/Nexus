package controllers

import (
	"net/http"

	"github.com/NutsBalls/Nexus/services"
	"github.com/NutsBalls/Nexus/utils"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type ImportController struct {
	importService *services.ImportService
}

func NewImportController(importService *services.ImportService) *ImportController {
	return &ImportController{importService: importService}
}

func (ic *ImportController) ImportDocument(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.JWTCustomClaims)

	file, err := c.FormFile("document")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No file uploaded"})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to open file"})
	}
	defer src.Close()

	document, err := ic.importService.ImportDocumentFromJSON(src, claims.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to import document"})
	}

	return c.JSON(http.StatusOK, document)
}
