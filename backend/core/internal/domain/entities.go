package domain

// TODO: refactor

// GetCurrentUserResponse godoc
// @Description Информация о текущем пользователе
type GetCurrentUserResponse struct {
	ID    uint   `json:"id" example:"1"`
	Login string `json:"login" example:"john_doe"`
	Role  string `json:"role" example:"user"`
}

// ErrorResponse godoc
// @Description Стандартизированный ответ при ошибке API
type ErrorResponse struct {
	Error struct {
		Code    string `json:"code" example:"NOT_FOUND"`
		Message string `json:"message" example:"Resource not found"`
	} `json:"error"`
}

// ApprovalResponse godoc
// @Description Информация о процессе одобрения файла
type ApprovalResponse struct {
	ID                uint   `json:"approval_id" example:"101"`
	FileID            uint   `json:"file_id" example:"789"`
	FileName          string `json:"file_name" example:"report.pdf"`
	Status            string `json:"status" example:"on approval"`
	WorkflowOrder     int    `json:"workflow_order" example:"2"`
	WorkflowUserCount int    `json:"workflow_user_count"`
}

type WorkflowResponse struct {
	WorkflowID     uint   `json:"workflow_id"`
	WorkflowName   string `json:"workflow_name"`
	WorkflowLength int    `json:"workflow_length"`
}

type ExtendedWorkflowResponse struct {
	WorkflowName string          `json:"workflow_name"`
	Stages       []WorkflowStage `json:"stages"`
}

type WorkflowStage struct {
	UserID uint `json:"user_id"`
	Order  int  `json:"order"`
}

type RoleResponse struct {
	RoleID   uint   `json:"role_id"`
	RoleName string `json:"role_name"`
}

type RoleData struct {
	RoleName string     `json:"role_name"`
	Users    []UserData `json:"users"`
}

type UserData struct {
	UserID uint   `json:"user_id"`
	Login  string `json:"login"`
}

// Directory модель
type Directory struct {
	ID           uint
	ParentPathID *uint
	Name         string
	Status       string
	Version      int
	ParentPath   *Directory
	WorkflowID   uint

	Files []File
}

// File модель
type File struct {
	ID          uint
	DirectoryID uint
	Name        string
	Status      string
	Version     int

	Directory *Directory
}
