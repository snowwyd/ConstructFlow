package domain

import "gorm.io/gorm"

// Directory модель
type Directory struct {
	gorm.Model
	ParentPathID *uint      `json:"parent_path_id" gorm:"index"` // Указатель для NULL
	Name         string     `json:"name" gorm:"not null"`
	Status       string     `json:"status" gorm:"not null"`
	Version      int        `json:"version" gorm:"default:1"`
	ParentPath   *Directory `gorm:"foreignkey:ParentPathID"`
	WorkflowID   uint       `json:"workflow_id"`

	Files []File `gorm:"foreignKey:DirectoryID"`
}

// File модель
type File struct {
	gorm.Model
	DirectoryID uint   `json:"directory_id" gorm:"not null;index"`
	Name        string `json:"name" gorm:"not null"`
	Status      string `json:"status" gorm:"not null"`
	Version     int    `json:"version" gorm:"default:1"`

	Directory *Directory `gorm:"foreignKey:DirectoryID"`
}
