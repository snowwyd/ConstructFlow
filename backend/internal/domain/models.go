package domain

type User struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	PassHash []byte `json:"pass_hash"`
	Role     string `json:"role"`
}

type Folder struct {
	ID       int    `json:"id"`
	ParentID int    `json:"parent_id"`
	Name     string `json:"name"`
	Status   string `json:"status"`
}

type File struct {
	ID       int    `json:"id"`
	FolderID int    `json:"folder_id"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	Version  int    `json:"version"`
}

type Approval struct {
	ID         int    `json:"id"`
	FileID     int    `json:"file_id"`
	UserID     int    `json:"user_id"`
	Status     string `json:"status"`
	WorkflowID int    `json:"workflow_id"`
	Order      int    `json:"order"`
}

type Annotation struct {
	ID         int    `json:"id"`
	ApprovalID int    `json:"approval_id"`
	Text       string `json:"text"`
}

type Workflow struct {
	ID      int   `json:"id"`
	UserIDs []int `json:"user_ids"`
}
