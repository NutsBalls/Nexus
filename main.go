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

	userController := controllers.NewUserController(db, cfg.JWTSecret)
	documentController := controllers.NewDocumentController(db)
	shareController := controllers.NewShareController(db)

	e.POST("/api/register", userController.Register)
	e.POST("/api/login", userController.Login)

	tagController := controllers.NewTagController(db)
	folderController := controllers.NewFolderController(db)

	api := e.Group("/api")
	api.Use(middlewares.JWTMiddleware(cfg.JWTSecret))

	shareGroup := api.Group("/shares")

	api.GET("/shares/shared-with-me", shareController.GetSharedWithMe)
	api.GET("/shares/shared-by-me", shareController.GetSharedByMe)
	shareGroup.GET("/:id/access", shareController.CheckDocumentAccess, middlewares.DocumentAccessMiddleware(db))

	api.POST("/documents/:id/share", shareController.ShareDocument, middlewares.DocumentAccessMiddleware(db))
	api.GET("/documents/:id/shares", shareController.GetDocumentShares, middlewares.DocumentAccessMiddleware(db))
	api.DELETE("/shares/:id", shareController.RemoveShare, middlewares.DocumentAccessMiddleware(db))

	api.GET("/documents", documentController.GetDocuments)
	api.POST("/documents", documentController.CreateDocument)
	api.GET("/documents/:id", documentController.GetDocument)
	api.PUT("/documents/:id", documentController.UpdateDocument)
	api.DELETE("/attachments/:id", documentController.DeleteAttachment)
	api.DELETE("/documents/:id", documentController.DeleteDocument)
	api.GET("/documents/search", documentController.SearchDocuments)
	api.POST("/documents/:id/versions", documentController.CreateVersion)
	api.GET("/documents/:id/versions", documentController.GetVersions)
	api.POST("/documents/:id/attachments", documentController.UploadAttachment)
	api.GET("/documents/:id/attachments", documentController.GetAttachments)
	api.GET("/download/*", documentController.DownloadAttachment)

	api.POST("/folders", folderController.CreateFolder)
	api.GET("/folders", folderController.GetFolders)
	api.GET("/folders/:id/documents", documentController.GetFolderDocuments)
	api.PUT("/folders/:id", folderController.UpdateFolder)
	api.DELETE("/folders/:id", folderController.DeleteFolder)

	api.POST("/tags", tagController.CreateTag)
	api.GET("/tags", tagController.GetTags)

	api.POST("/documents/:id/attachments", documentController.UploadAttachment)
	api.GET("/documents/:id/attachments", documentController.GetAttachments)

	exportService := services.NewExportService(db)
	importService := services.NewImportService(db)

	exportController := controllers.NewExportController(exportService)
	importController := controllers.NewImportController(importService)

	api.GET("/documents/:id/export", exportController.ExportDocument)
	api.POST("/documents/import", importController.ImportDocument)

	api.GET("/search/tags", tagController.SearchByTag)

	commentController := controllers.NewCommentController(db)
	notificationController := controllers.NewNotificationController(db)

	api.POST("/documents/:id/comments", commentController.AddComment)
	api.GET("/documents/:id/comments", commentController.GetComments)
	api.DELETE("/documents/:id/comments/:commentId", commentController.DeleteComment)

	api.GET("/notifications", notificationController.GetNotifications)
	api.PUT("/notifications/:id/read", notificationController.MarkAsRead)
	api.PUT("/notifications/read-all", notificationController.MarkAllAsRead)

	documentGroup := api.Group("/documents/:id")
	documentGroup.Use(middlewares.DocumentAccessMiddleware(db))

	documentGroup.PUT("", documentController.UpdateDocument)
	documentGroup.DELETE("", documentController.DeleteDocument)

	uploadsDir := filepath.Join(currentDir, "uploads")
	e.Static("/uploads", uploadsDir)
	e.Logger.Fatal(e.Start("0.0.0.0:" + cfg.ServerPort))
}
