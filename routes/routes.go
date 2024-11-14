package routes

import (
	"github.com/NutsBalls/Nexus/controllers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"os"
)

func InitRoutes(e *echo.Echo) {
	// Маршруты для аутентификации
	e.POST("/register", controllers.Register)
	e.POST("/login", controllers.Login)

	// Защищенные маршруты для документов
	r := e.Group("/documents")
	r.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))
	r.POST("", controllers.CreateDocument)
	r.GET("", controllers.GetDocuments)
	r.PUT("/:id", controllers.UpdateDocument)
	r.DELETE("/:id", controllers.DeleteDocument)
}
