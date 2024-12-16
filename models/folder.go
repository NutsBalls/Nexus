package models

import (
	"time"
)

type Folder struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `gorm:"index" json:"deleted_at,omitempty"`

	Name      string     `gorm:"not null" json:"name" example:"Личные документы"`
	UserID    uint       `json:"user_id" example:"1"`
	User      User       `gorm:"foreignKey:UserID" json:"-"`
	ParentID  *uint      `json:"parent_id,omitempty"`
	Parent    *Folder    `gorm:"foreignKey:ParentID" json:"-"`
	Children  []*Folder  `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Documents []Document `gorm:"foreignKey:FolderID" json:"documents,omitempty"`
}
