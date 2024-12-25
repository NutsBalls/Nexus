package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/NutsBalls/Nexus/config"
	"github.com/NutsBalls/Nexus/controllers"
	_ "github.com/NutsBalls/Nexus/docs"
	"github.com/NutsBalls/Nexus/middlewares"
	"github.com/NutsBalls/Nexus/services"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

var db *gorm.DB

func main() {
	e := echo.New()

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders: []string{echo.HeaderAuthorization, echo.HeaderContentType},
	}))

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	uploadsDir := filepath.Join(currentDir, "uploads")
	e.Static("/uploads", uploadsDir)

	exportService := services.NewExportService(db)
	importService := services.NewImportService(db)

	// Используемые контроллеры
	userController := controllers.NewUserController(db, cfg.JWTSecret)
	documentController := controllers.NewDocumentController(db)
	shareController := controllers.NewShareController(db)
	exportController := controllers.NewExportController(exportService)
	importController := controllers.NewImportController(importService)
	folderController := controllers.NewFolderController(db)

	// Аутентификация
	e.POST("/api/register", userController.Register)
	e.POST("/api/login", userController.Login)

	// Возможность делиться с другими пользователями
	shareGroup := e.Group("/api/shares")
	shareGroup.Use(middlewares.DocumentAccessMiddleware(db))
	shareGroup.GET("/shared-with-me", shareController.GetSharedWithMe)
	shareGroup.GET("/shared-by-me", shareController.GetSharedByMe)
	shareGroup.GET("/:id/access", shareController.CheckDocumentAccess)

	// Защищенные маршруты
	api := e.Group("/api")
	api.Use(middlewares.JWTMiddleware(cfg.JWTSecret))

	// documents
	api.GET("/documents", documentController.GetDocuments)
	api.GET("/documents/:id", documentController.GetDocument)
	api.GET("/documents/:id/attachments", documentController.GetAttachments)
	api.GET("/download/*", documentController.DownloadAttachment)
	api.GET("/documents/search", documentController.SearchDocuments)
	api.POST("/documents", documentController.CreateDocument)
	api.POST("/documents/:id/attachments", documentController.UploadAttachment)
	api.PUT("/documents/:id", documentController.UpdateDocument)
	api.DELETE("/attachments/:id", documentController.DeleteAttachment)
	api.DELETE("/documents/:id", documentController.DeleteDocument)

	api.POST("/documents/:id/share", shareController.ShareDocument)
	api.GET("/documents/:id/shares", shareController.GetDocumentShares, middlewares.DocumentAccessMiddleware(db))
	api.DELETE("/shares/:id", shareController.RemoveShare, middlewares.DocumentAccessMiddleware(db))

	// folders
	api.GET("/folders", folderController.GetFolders)
	api.GET("/folders/:id/documents", documentController.GetFolderDocuments)
	api.POST("/folders", folderController.CreateFolder)
	api.PUT("/folders/:id", folderController.UpdateFolder)
	api.DELETE("/folders/:id", folderController.DeleteFolder)

	// import/export
	api.GET("/documents/:id/export", exportController.ExportDocument)
	api.POST("/documents/import", importController.ImportDocument)

	documentGroup := api.Group("/documents/:id")
	documentGroup.Use(middlewares.DocumentAccessMiddleware(db))
	documentGroup.PUT("", documentController.UpdateDocument)
	documentGroup.DELETE("", documentController.DeleteDocument)

	e.Logger.Fatal(e.Start("0.0.0.0:" + cfg.ServerPort))

}
