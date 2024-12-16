package main

import (
	"log"

	"github.com/NutsBalls/Nexus/config"
	"github.com/NutsBalls/Nexus/controllers"
	_ "github.com/NutsBalls/Nexus/docs"
	"github.com/NutsBalls/Nexus/middlewares"
	"github.com/NutsBalls/Nexus/services"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Nexus
// @version 1.0
// @description Описание вашего API.
// @host localhost:8080
// @BasePath /

var db *gorm.DB

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Инициализация базы данных
	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Инициализация контроллеров
	userController := controllers.NewUserController(db, cfg.JWTSecret)
	documentController := controllers.NewDocumentController(db)

	// Публичные маршруты
	e.POST("/api/register", userController.Register)
	e.POST("/api/login", userController.Login)

	// Инициализация новых контроллеров
	tagController := controllers.NewTagController(db)
	folderController := controllers.NewFolderController(db)
	collaborationController := controllers.NewCollaborationController(db)

	// Защищенные маршруты
	api := e.Group("/api")
	api.Use(middlewares.JWTMiddleware(cfg.JWTSecret)) // Добавляем middleware для JWT

	// Маршруты для документов
	api.GET("/documents", documentController.GetDocuments)
	api.POST("/documents", documentController.CreateDocument)
	api.GET("/documents/:id", documentController.GetDocument)
	api.PUT("/documents/:id", documentController.UpdateDocument)
	api.DELETE("/documents/:id", documentController.DeleteDocument)
	api.GET("/documents/search", documentController.SearchDocuments)
	api.POST("/documents/:id/share", documentController.ShareDocument)
	api.POST("/documents/:id/versions", documentController.CreateVersion)
	api.GET("/documents/:id/versions", documentController.GetVersions)

	// Маршруты для папок
	api.POST("/folders", folderController.CreateFolder)
	api.GET("/folders", folderController.GetFolders)
	api.PUT("/folders/:id", folderController.UpdateFolder)
	api.DELETE("/folders/:id", folderController.DeleteFolder)

	// Маршруты для тегов
	api.POST("/tags", tagController.CreateTag)
	api.GET("/tags", tagController.GetTags)

	favoriteController := controllers.NewFavoriteController(db)
	recentController := controllers.NewRecentController(db)

	// Маршруты для избранного
	api.POST("/documents/:id/favorite", favoriteController.AddToFavorites)
	api.DELETE("/documents/:id/favorite", favoriteController.RemoveFromFavorites)
	api.GET("/favorites", favoriteController.GetFavorites)

	// Маршруты для недавних документов
	api.POST("/documents/:id/recent", recentController.AddToRecent)
	api.GET("/recent", recentController.GetRecentDocuments)

	// Маршруты для вложений
	api.POST("/documents/:id/attachments", documentController.UploadAttachment)
	api.GET("/documents/:id/attachments", documentController.GetAttachments)

	exportService := services.NewExportService(db)
	importService := services.NewImportService(db)

	// Инициализация контроллеров
	exportController := controllers.NewExportController(exportService)
	importController := controllers.NewImportController(importService)

	// Экспорт/Импорт маршруты
	api.GET("/documents/:id/export", exportController.ExportDocument)
	api.POST("/documents/import", importController.ImportDocument)

	// Поиск по тегам
	api.GET("/search/tags", tagController.SearchByTag)

	commentController := controllers.NewCommentController(db)
	notificationController := controllers.NewNotificationController(db)

	// Маршруты для комментариев
	api.POST("/documents/:id/comments", commentController.AddComment)
	api.GET("/documents/:id/comments", commentController.GetComments)
	api.DELETE("/documents/:id/comments/:commentId", commentController.DeleteComment)

	// Маршруты для уведомлений
	api.GET("/notifications", notificationController.GetNotifications)
	api.PUT("/notifications/:id/read", notificationController.MarkAsRead)
	api.PUT("/notifications/read-all", notificationController.MarkAllAsRead)

	// Добавим middleware для проверки прав доступа к документам
	documentGroup := api.Group("/documents/:id")
	documentGroup.Use(middlewares.DocumentAccessMiddleware(db))

	// Защищенные маршруты документов
	documentGroup.PUT("", documentController.UpdateDocument)
	documentGroup.DELETE("", documentController.DeleteDocument)
	documentGroup.POST("/collaborators", collaborationController.AddCollaborator)
	documentGroup.GET("/collaborators", collaborationController.GetCollaborators)
	documentGroup.DELETE("/collaborators/:userId", collaborationController.RemoveCollaborator)

	// Запуск сервера
	e.Logger.Fatal(e.Start(":" + cfg.ServerPort))
}
