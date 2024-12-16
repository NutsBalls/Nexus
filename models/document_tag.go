package models

type DocumentTag struct {
	DocumentID uint `gorm:"primaryKey"`
	TagID      uint `gorm:"primaryKey"`
}
