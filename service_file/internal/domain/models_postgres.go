package domain

import "gorm.io/gorm"

// Directory модель
type Directory struct {
	gorm.Model
	ParentPathID *uint  `gorm:"index"` // Указатель для NULL
	Name         string `gorm:"not null"`
	Status       string `gorm:"not null"`
	Version      int    `gorm:"default:1"`
	WorkflowID   uint   `gorm:"not null"` // Логическая связь с core-service

	// Родительская директория
	ParentPath *Directory `gorm:"foreignKey:ParentPathID"`

	// Связь с файлами
	Files []File `gorm:"foreignKey:DirectoryID"`
}

// File модель
type File struct {
	gorm.Model
	DirectoryID uint   `gorm:"not null;index"`
	Name        string `gorm:"not null"`
	Status      string `gorm:"not null"`
	Version     int    `gorm:"default:1"`

	// Ссылка на MinIO
	MinioObjectKey string `gorm:"not null"`

	// Связь с директорией
	Directory Directory `gorm:"foreignKey:DirectoryID"`
}

// UserDirectory связующая таблица
type UserDirectory struct {
	UserID      uint `gorm:"primaryKey;column:user_id;foreignKey:User"`
	DirectoryID uint `gorm:"primaryKey;column:directory_id;foreignKey:Directory"`
}

// UserFile связующая таблица
type UserFile struct {
	UserID uint `gorm:"primaryKey;column:user_id;foreignKey:User"`
	FileID uint `gorm:"primaryKey;column:file_id;foreignKey:File"`
}
