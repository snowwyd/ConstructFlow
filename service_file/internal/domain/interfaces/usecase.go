package interfaces

import (
	"context"
	"service-file/internal/domain"
)

type FileTreeUsecase interface {
	GetFileTree(ctx context.Context, isArchive bool, userID uint) (data domain.GetFileTreeResponse, err error)
	GetFileInfo(ctx context.Context, fileID, userID uint) (fileInfo *domain.FileResponse, err error)

	CreateDirectory(ctx context.Context, directoryID *uint, name string, userID uint) (err error)
	UploadFile(ctx context.Context, directoryID uint, name string, userID uint) (err error)

	DeleteFile(ctx context.Context, fileID uint, userID uint) (err error)
	DeleteDirectory(ctx context.Context, directoryID uint, userID uint) (err error)
}
