package domain

// GetCurrentUserResponse godoc
// @Description Информация о текущем пользователе
type GetCurrentUserResponse struct {
	ID    uint   `json:"id" example:"1"`
	Login string `json:"login" example:"john_doe"`
	Role  string `json:"role" example:"user"`
}

// GetFileTreeResponse godoc
// @Description Ответ с древовидной структурой файлов и директорий
type GetFileTreeResponse struct {
	Data []DirectoryResponse `json:"data"`
}

// DirectoryResponse godoc
// @Description Детальная информация о директории
type DirectoryResponse struct {
	ID           uint           `json:"id" example:"123"`
	NameFolder   string         `json:"name_folder" example:"Documents"`
	Status       string         `json:"status" example:"active"`
	ParentPathID *uint          `json:"parent_path_id,omitempty" example:"456"`
	Files        []FileResponse `json:"files"`
}

// FileResponse godoc
// @Description Детальная информация о файле
type FileResponse struct {
	ID          uint   `json:"id" example:"789"`
	NameFile    string `json:"name_file" example:"report.pdf"`
	Status      string `json:"status" example:"draft"`
	DirectoryID uint   `json:"directory_id" example:"123"`
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
	ID            uint   `json:"id" example:"101"`
	FileID        uint   `json:"file_id" example:"789"`
	FileName      string `json:"file_name" example:"report.pdf"`
	Status        string `json:"status" example:"on approval"`
	WorkflowOrder int    `json:"workflow_order" example:"2"`
}
