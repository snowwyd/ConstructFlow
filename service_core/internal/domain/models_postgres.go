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
