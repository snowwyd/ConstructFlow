package interfaces

import (
	"context"
	"service-file/internal/domain"

	"gorm.io/gorm"
)

type DirectoryRepository interface {
	GetFileTree(ctx context.Context, isArchive bool, userID uint) ([]domain.Directory, error)
	UpdateDirectories(ctx context.Context, workflowID uint, directoryIDs []uint) error

	CreateDirectory(ctx context.Context, parentPathID *uint, name string, status string, userID uint) error
	DeleteDirectory(ctx context.Context, directoryID uint, userID uint) error

	CheckUserDirectoryAccess(ctx context.Context, userID, directoryID uint) (bool, error)
	CheckWorkflow(ctx context.Context, workflowID uint) (bool, error)

	DeleteUserRelations(ctx context.Context, userID uint) error

	CheckDirectoriesExist(ctx context.Context, directoryIDs []uint) (bool, error)

	UpdateUserDirectoryRelations(ctx context.Context, userID uint, directoryIDs []uint) error

	WithTx(tx *gorm.DB) DirectoryRepository // Метод для передачи транзакции
	GetDB() *gorm.DB
}

type FileMetadataRepository interface {
	GetFileByID(ctx context.Context, fileID uint) (*domain.File, error)
	GetFilesByID(ctx context.Context, fileIDs []uint32, files *[]domain.File) error
	GetFileInfo(ctx context.Context, fileID uint, userID uint) (*domain.File, error)

	CreateFile(ctx context.Context, directoryID uint, name string, status string, userID uint, minioKey string, size int64, contentType string) error
	UpdateFile(ctx context.Context, file *domain.File) error
	UpdateFileStatus(ctx context.Context, fileID uint, status string, tx *gorm.DB) error
	DeleteFile(ctx context.Context, fileID uint, userID uint) error

	CheckUserFileAccess(ctx context.Context, userID, fileID uint) (bool, error)
	CheckFilesExist(ctx context.Context, fileIDs []uint) (bool, error)

	UpdateUserFileRelations(ctx context.Context, userID uint, fileIDs []uint) error

	DeleteUserRelations(ctx context.Context, userID uint) error

	WithTx(tx *gorm.DB) FileMetadataRepository // Метод для передачи транзакции
	GetDB() *gorm.DB
}
