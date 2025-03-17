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

	Directories []Directory `gorm:"many2many:user_directories;"`
	Files       []File      `gorm:"many2many:user_files;"`
}

// Directory модель
type Directory struct {
	gorm.Model
	ParentPathID *uint      `json:"parent_path_id" gorm:"index"` // Указатель для NULL
	Name         string     `json:"name" gorm:"not null"`
	Status       string     `json:"status" gorm:"not null"`
	Version      int        `json:"version" gorm:"default:1"`
	ParentPath   *Directory `gorm:"foreignkey:ParentPathID"`

	Users []User `gorm:"many2many:user_directories;"`
	Files []File `gorm:"foreignKey:DirectoryID"`
}

// UserDirectory связующая таблица
type UserDirectory struct {
	UserID      uint `gorm:"primaryKey"`
	DirectoryID uint `gorm:"primaryKey"`
}

// File модель
type File struct {
	gorm.Model
	DirectoryID uint   `json:"directory_id" gorm:"not null;index"`
	Name        string `json:"name" gorm:"not null"`
	Status      string `json:"status" gorm:"not null"`
	Version     int    `json:"version" gorm:"default:1"`

	Users []User `gorm:"many2many:user_files;"`
}

// UserFile связующая таблица
type UserFile struct {
	UserID uint `gorm:"primaryKey"`
	FileID uint `gorm:"primaryKey"`
}

// Approval модель
type Approval struct {
	gorm.Model
	FileID     uint   `json:"file_id" gorm:"not null;index"`
	UserID     uint   `json:"user_id" gorm:"not null;index"`
	Status     string `json:"status" gorm:"not null"`
	WorkflowID uint   `json:"workflow_id" gorm:"not null;index"`
	Order      int    `json:"order" gorm:"not null"`
}

// Annotation модель
type Annotation struct {
	gorm.Model
	ApprovalID uint   `json:"approval_id" gorm:"not null;index"`
	Text       string `json:"text" gorm:"not null"`
}

// Workflow модель
type Workflow struct {
	gorm.Model
	UserID uint `json:"user_id" gorm:"not null;index"`
	Order  int  `json:"order" gorm:"not null"`
}
