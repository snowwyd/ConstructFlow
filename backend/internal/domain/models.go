package domain

import "gorm.io/gorm"

// TODO: переопределить в соответствии с правками
type User struct {
	gorm.Model
	Login    string `json:"login" gorm:"unique;not null"`
	PassHash []byte `json:"pass_hash" gorm:"not null"`
	Role     string `json:"role" gorm:"not null"`
}

type Folder struct {
	gorm.Model
	ParentID int    `json:"parent_id"`
	Name     string `json:"name" gorm:"not null"`
	Status   string `json:"status" gorm:"not null"`
}

type File struct {
	gorm.Model
	FolderID uint   `json:"folder_id" gorm:"not null;index"` // Ссылка на папку
	Name     string `json:"name" gorm:"not null"`
	Status   string `json:"status" gorm:"not null"`
	Version  int    `json:"version" gorm:"default:1"` // Начальная версия файла
}

type Approval struct {
	gorm.Model
	FileID     int    `json:"file_id" gorm:"not null;index"` // Ссылка на файл
	UserID     int    `json:"user_id" gorm:"not null;index"` // Ссылка на пользователя
	Status     string `json:"status" gorm:"not null"`
	WorkflowID int    `json:"workflow_id" gorm:"not null;index"` // Ссылка на рабочий процесс
	Order      int    `json:"order" gorm:"not null"`             // Порядок согласования
}

type Annotation struct {
	gorm.Model
	ApprovalID int    `json:"approval_id" gorm:"not null;index"` // Ссылка на согласование
	Text       string `json:"text" gorm:"not null"`
}

type Workflow struct {
	gorm.Model
	UserID int `json:"user_id" gorm:"not null;index"` // Ссылка на пользователя
	Order  int `json:"order" gorm:"not null"`         // Порядок в рабочем процессе
}
