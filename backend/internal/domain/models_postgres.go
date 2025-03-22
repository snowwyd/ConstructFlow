package domain

import "gorm.io/gorm"

// Role модель
type Role struct {
	gorm.Model
	RoleName string `gorm:"unique;not null"`
}

// User модель
type User struct {
	gorm.Model
	Login    string `json:"login" gorm:"unique;not null"`
	PassHash []byte `json:"pass_hash" gorm:"not null"`
	RoleID   uint   `gorm:"not null"`

	Directories []Directory `gorm:"many2many:user_directories"`
	Files       []File      `gorm:"many2many:user_files"`
}

// Directory модель
type Directory struct {
	gorm.Model
	ParentPathID *uint      `json:"parent_path_id" gorm:"index"` // Указатель для NULL
	Name         string     `json:"name" gorm:"not null"`
	Status       string     `json:"status" gorm:"not null"`
	Version      int        `json:"version" gorm:"default:1"`
	ParentPath   *Directory `gorm:"foreignkey:ParentPathID"`
	WorkflowID   uint       `json:"workflow_id"`

	Users []User `gorm:"many2many:user_directories;"`
	Files []File `gorm:"foreignKey:DirectoryID"`
}

// UserDirectory связующая таблица
type UserDirectory struct {
	UserID      uint `gorm:"primaryKey;column:user_id;foreignKey:User"`
	DirectoryID uint `gorm:"primaryKey;column:directory_id;foreignKey:Directory"`
}

// File модель
type File struct {
	gorm.Model
	DirectoryID uint   `json:"directory_id" gorm:"not null;index"`
	Name        string `json:"name" gorm:"not null"`
	Status      string `json:"status" gorm:"not null"`
	Version     int    `json:"version" gorm:"default:1"`

	Directory *Directory `gorm:"foreignKey:DirectoryID"`
	Users     []User     `gorm:"many2many:user_files;"`
}

// UserFile связующая таблица
type UserFile struct {
	UserID uint `gorm:"primaryKey;column:user_id;foreignKey:User"`
	FileID uint `gorm:"primaryKey;column:file_id;foreignKey:File"`
}

// Approval модель
type Approval struct {
	gorm.Model
	FileID         uint   `json:"file_id" gorm:"not null;index"`
	Status         string `json:"status" gorm:"not null"`
	WorkflowID     uint   `json:"workflow_id" gorm:"not null"`
	WorkflowOrder  int    `json:"workflow_order" gorm:"not null"`
	AnnotationText string `json:"annotation_text"`
}

// Workflow модель
type Workflow struct {
	gorm.Model
	WorkflowID    uint `json:"workflow_id" gorm:"not null"`
	UserID        uint `json:"user_id" gorm:"not null;index"`
	WorkflowOrder int  `json:"workflow_order" gorm:"not null"`
}
