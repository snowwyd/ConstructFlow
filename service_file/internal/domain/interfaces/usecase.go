package interfaces

import (
	"context"
	"service-file/internal/domain"
)

type FileTreeUsecase interface {
	// REST usecases
	GetFileByID(ctx context.Context, fileID uint) (domain.FileResponse, error)
	GetDirectoryByID(ctx context.Context, directoryID uint) (*domain.DirectoryResponse, error)
	GetFileTree(ctx context.Context, isArchive bool, userID uint) (data domain.GetFileTreeResponse, err error)

	CreateFile(ctx context.Context, directoryID uint, name string, userID uint) (err error)
	CreateDirectory(ctx context.Context, directoryID *uint, name string, userID uint) (err error)

	DeleteFile(ctx context.Context, fileID uint, userID uint) (err error)
	DeleteDirectory(ctx context.Context, directoryID uint, userID uint) (err error)

	// gRPC usecases
	UpdateFileStatus(ctx context.Context, fileID uint, status string) error
	CheckAccessToFile(ctx context.Context, fileID, userID uint) (bool, error)
	CheckAccessToDirectory(ctx context.Context, directoryID, userID uint) (bool, error)
}
