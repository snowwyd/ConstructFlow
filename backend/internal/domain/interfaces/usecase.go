package interfaces

import (
	"backend/internal/domain"
	"context"
)

type AuthUsecase interface {
	Login(ctx context.Context, login, password string) (token string, err error)
	GetCurrentUser(ctx context.Context, userID uint) (userInfo domain.GetCurrentUserResponse, err error)

	// для админа и локального тестирования
	RegisterUser(ctx context.Context, login, password string, roleID uint) (userID uint, err error)
	RegisterRole(ctx context.Context, roleName string) (roleID uint, err error)
}

type FileTreeUsecase interface {
	GetFileTree(ctx context.Context, isArchive bool) (data domain.GetFileTreeResponse, err error)
	GetFileInfo(ctx context.Context, fileID uint) (file domain.File, err error)

	UploadFile(ctx context.Context, directoryID uint, name string) (err error)
	UploadDirectory(ctx context.Context, directoryID uint, name string) (err error)

	DeleteFile(ctx context.Context, fileID uint) (err error)
	DeleteDirectory(ctx context.Context, directoryID uint) (err error)
}
