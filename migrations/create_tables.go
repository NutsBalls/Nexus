package migrations

import (
	"github.com/NutsBalls/Nexus/models"
	"github.com/jinzhu/gorm"
)

// Migrate создаёт таблицы в базе данных
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Document{},
		&models.Folder{},
		&models.Tag{},
		&models.DocumentTag{},
	).Error
}
