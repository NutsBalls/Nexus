package models

import (
	"time"
)

type Document struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `gorm:"index" json:"deleted_at,omitempty"`

	Title       string    `gorm:"not null" json:"title" example:"Мой документ"`
	Content     string    `json:"content" example:"Содержимое документа"`
	UserID      uint      `json:"user_id" example:"1"`
	User        User      `gorm:"foreignKey:UserID" json:"-"`
	FolderID    *uint     `json:"folder_id,omitempty" example:"2"`
	Folder      Folder    `gorm:"foreignKey:FolderID;constraint:OnDelete:CASCADE;" json:"-"`
	Tags        []Tag     `gorm:"many2many:document_tags;" json:"tags"`
	Versions    []Version `gorm:"foreignKey:DocumentID" json:"versions,omitempty"`
	SharedUsers []User    `gorm:"many2many:document_shares;" json:"shared_users,omitempty"`
	IsPublic    bool      `gorm:"default:false" json:"is_public" example:"false"`
	Shares      []Share   `gorm:"foreignKey:DocumentID" json:"shares"`
}

type Version struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `gorm:"index" json:"deleted_at,omitempty"`

	DocumentID uint   `json:"document_id"`
	Content    string `json:"content"`
	Title      string `json:"title"`
	ChangeLog  string `json:"change_log"`
}
