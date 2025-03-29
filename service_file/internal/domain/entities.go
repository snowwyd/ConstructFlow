package domain

type GetFileTreeResponse struct {
	Data []DirectoryResponse `json:"data"`
}

type DirectoryResponse struct {
	ID           uint           `json:"id" example:"123"`
	NameFolder   string         `json:"name_folder" example:"Documents"`
	Status       string         `json:"status" example:"active"`
	ParentPathID *uint          `json:"parent_path_id,omitempty" example:"456"`
	Files        []FileResponse `json:"files"`
}

type FileResponse struct {
	ID          uint   `json:"id" example:"789"`
	NameFile    string `json:"name_file" example:"report.pdf"`
	Status      string `json:"status" example:"draft"`
	DirectoryID uint   `json:"directory_id" example:"123"`
}

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code" example:"NOT_FOUND"`
		Message string `json:"message" example:"Resource not found"`
	} `json:"error"`
}
