package config

import (
	"fmt"
	"log"
	"os"

	"github.com/NutsBalls/Nexus/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Миграции для всех таблиц
	db.AutoMigrate(
		&models.User{},
		&models.Document{},
		&models.Folder{},
		&models.Tag{},
		&models.DocumentTag{},
	)

	DB = db
}
