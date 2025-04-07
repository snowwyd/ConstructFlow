package interfaces

import (
	"context"
	"service-file/internal/domain"

	"github.com/minio/minio-go/v7"
)

type DirectoryUsecase interface {
	GetFileTree(ctx context.Context, isArchive bool, userID uint) (domain.GetFileTreeResponse, error)
	CreateDirectory(ctx context.Context, parentPathID *uint, name string, userID uint) error
	DeleteDirectory(ctx context.Context, directoryID uint, userID uint) error
}

type FileUsecase interface {
	GetFileInfo(ctx context.Context, fileID uint, userID uint) (*domain.FileResponse, error)
	CreateFile(ctx context.Context, directoryID uint, name string, fileData []byte, contentType string, userID uint) error
	DownloadFileDirect(ctx context.Context, fileID uint, userID uint) (*domain.File, *minio.Object, error)
	DeleteFile(ctx context.Context, fileID uint, userID uint) error
	UpdateFile(ctx context.Context, fileID uint, newData []byte, userID uint) error

	ConvertSTPToGLTF(ctx context.Context, fileID uint, userID uint) (string, error)
}

type GRPCUsecase interface {
	GetFileByID(ctx context.Context, fileID uint) (*domain.File, error)
	UpdateFileStatus(ctx context.Context, fileID uint, status string) error
	GetFilesByID(ctx context.Context, fileIDs []uint32) ([]domain.File, error)
	CheckWorkflow(ctx context.Context, workflowID uint) (bool, error)
}
