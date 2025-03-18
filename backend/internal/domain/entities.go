package domain

type GetCurrentUserResponse struct {
	ID    uint   `json:"id"`
	Login string `json:"login"`
	Role  string `json:"role"`
}

type GetFileTreeResponse struct {
	Data []DirectoryResponse `json:"data"`
}

type FileResponse struct {
	ID          uint   `json:"id"`
	NameFile    string `json:"name_file"`
	Status      string `json:"status"`
	DirectoryID uint   `json:"directory_id"`
}

type DirectoryResponse struct {
	NameFolder   string         `json:"name_folder"`
	Status       string         `json:"status"`
	ParentPathID *uint          `json:"parent_path_id,omitempty"`
	Files        []FileResponse `json:"files"`
}
