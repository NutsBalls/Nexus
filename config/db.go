package config

import (
	"log"
	"os"

	"github.com/NutsBalls/Nexus/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB инициализирует подключение к базе данных с использованием DATABASE_URL
func InitDB() {
	// Получаем строку подключения из переменной окружения
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set in the environment")
	}

	// Подключаемся к базе данных с использованием GORM и строки подключения
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Автоматическая миграция моделей в базе данных
	err = db.AutoMigrate(&models.User{}, &models.Document{}) // Можно добавить другие модели, если нужно
	if err != nil {
		log.Fatalf("failed to migrate models: %v", err)
	}

	// Присваиваем глобальную переменную DB для использования в других частях приложения
	DB = db
}
