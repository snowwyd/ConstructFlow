package domain

type GetCurrentUserResponse struct {
	ID    uint   `json:"id"`
	Login string `json:"login"`
	Role  string `json:"role"`
}

type GetFileTreeResponse struct {
	Data []DirectoryResponse `json:"data"`
}

type DirectoryResponse struct {
	NameFolder string         `json:"name_folder"`
	Status     string         `json:"status"`
	ParentID   uint           `json:"parent_id"`
	Files      []FileResponse `json:"files"`
}

type FileResponse struct {
	ID       uint   `json:"id"`
	NameFile string `json:"name_file"`
	Status   string `json:"status"`
	ParentID uint   `json:"parent_id"`
}
