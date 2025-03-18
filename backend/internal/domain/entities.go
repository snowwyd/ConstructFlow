package domain

// GetCurrentUserResponse godoc
// @Description Информация о пользователе
type GetCurrentUserResponse struct {
	ID    uint   `json:"id"`
	Login string `json:"login"`
	Role  string `json:"role"`
}

// GetFileTreeResponse godoc
// @Description Структура ответа для дерева файлов
type GetFileTreeResponse struct {
	Data []DirectoryResponse `json:"data"`
}

// DirectoryResponse godoc
// @Description Информация о директории
type DirectoryResponse struct {
	ID           uint           `json:"id"`
	NameFolder   string         `json:"name_folder"`
	Status       string         `json:"status"`
	ParentPathID *uint          `json:"parent_path_id,omitempty"`
	Files        []FileResponse `json:"files"`
}

// FileResponse godoc
// @Description Информация о файле
type FileResponse struct {
	ID          uint   `json:"id"`
	NameFile    string `json:"name_file"`
	Status      string `json:"status"`
	DirectoryID uint   `json:"directory_id"`
}

// ErrorResponse godoc
// @Description Структура ошибки API
type ErrorResponse struct {
	Error struct {
		Code    string `json:"code" example:"INVALID_REQUEST"`
		Message string `json:"message" example:"Invalid request body"`
	} `json:"error"`
}
