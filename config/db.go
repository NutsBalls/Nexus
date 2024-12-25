package config

import (
	"fmt"
	"log"

	"github.com/NutsBalls/Nexus/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(cfg *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Printf("Ошибка подключения к базе данных: %v", err)
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Folder{}); err != nil {
		log.Printf("Ошибка миграции базы данных для User и Folder: %v", err)
		return nil, err
	}

	if err := db.AutoMigrate(
		&models.Document{},
		&models.Tag{},
		&models.Version{},
		&models.Share{},
	); err != nil {
		log.Printf("Ошибка миграции базы данных для остальных моделей: %v", err)
		return nil, err
	}

	log.Println("Успешное подключение к базе данных!")
	return db, nil
}
