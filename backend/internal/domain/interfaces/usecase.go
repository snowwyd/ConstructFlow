package interfaces

import (
	"backend/internal/domain"
	"context"
)

type AuthUsecase interface {
	Login(ctx context.Context, login, password string) (token string, err error)
	GetCurrentUser(ctx context.Context, userID uint) (userInfo domain.GetCurrentUserResponse, err error)

	// для админа и локального тестирования
	RegisterUser(ctx context.Context, login, password string, roleID uint) (err error)
	RegisterRole(ctx context.Context, roleName string) (err error)
}

type FileTreeUsecase interface {
	GetFileTree(ctx context.Context, isArchive bool, userID uint) (data domain.GetFileTreeResponse, err error)
	GetFileInfo(ctx context.Context, fileID, userID uint) (fileInfo *domain.FileResponse, err error)

	CreateDirectory(ctx context.Context, directoryID *uint, name string, userID uint) (err error)
	UploadFile(ctx context.Context, directoryID uint, name string, userID uint) (err error)

	DeleteFile(ctx context.Context, fileID uint, userID uint) (err error)
	DeleteDirectory(ctx context.Context, directoryID uint, userID uint) (err error)
}

type ApprovalUsecase interface {
	ApproveFile(ctx context.Context, fileID uint) (err error)
	GetApprovalsByUserID(ctx context.Context, userID uint) (approvals []domain.ApprovalResponse, err error)
	SignApproval(ctx context.Context, approvalID, userID uint) error
	AnnotateApproval(ctx context.Context, approvalID, userID uint, message string) error
	FinalizeApproval(ctx context.Context, approvalID, userID uint) error
}
