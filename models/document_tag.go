package models

// DocumentTag represents the many-to-many relationship between documents and tags
type DocumentTag struct {
	DocumentID uint `gorm:"primaryKey"`
	TagID      uint `gorm:"primaryKey"`
}
