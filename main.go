package main

import (
	"github.com/NutsBalls/Nexus/config"
	"github.com/NutsBalls/Nexus/controllers"
	_ "github.com/NutsBalls/Nexus/docs"
	"github.com/NutsBalls/Nexus/models"
	"github.com/NutsBalls/Nexus/services"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/swaggo/echo-swagger"
	"log"
)

// @title Nexus
// @version 1.0
// @description Описание вашего API.
// @host localhost:8080
// @BasePath /

var db *gorm.DB

func main() {
	// Load config from .env
	config.LoadConfig()

	// Connect to database
	var err error
	db, err = gorm.Open("postgres", config.GetDatabaseURL())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.AutoMigrate(&models.User{}, &models.Folder{}, &models.Document{}, &models.Tag{}, &models.DocumentTag{})

	// Create Echo instance
	e := echo.New()
	// @Summary Регистрация нового пользователя
	// @Description Создает учетную запись пользователя
	// @Tags users
	// @Accept json
	// @Produce json
	// @Param user body User true "User data"
	// @Success 200 {object} User
	// @Failure 400 {object} map[string]string "Bad Request"
	// @Router /api/register [post]
	e.POST("/api/register", func(c echo.Context) error {
		// Логика обработчика
		return c.JSON(200, "User registered")
	})

	// Enable CORS
	e.Use(middleware.CORS())

	// Initialize services
	userService := services.UserService{DB: db}
	userController := controllers.UserController{UserService: userService}

	// Routes
	e.POST("/api/register", userController.Register)
	e.POST("/api/login", userController.Login)

	// Swagger documentation route
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
