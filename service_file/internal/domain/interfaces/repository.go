package interfaces

import (
	"context"
	"service-file/internal/domain"

	"gorm.io/gorm"
)

type FileMetadataRepository interface {
	GetFileByID(ctx context.Context, fileID uint, tx *gorm.DB) (*domain.File, error)
	GetDirectoryByID(ctx context.Context, directoryID uint) (*domain.Directory, error)
	GetFileTree(ctx context.Context, isArchive bool, userID uint) ([]domain.Directory, error)

	CreateFile(ctx context.Context, directoryID uint, name string, status string, userID uint) error
	CreateDirectory(ctx context.Context, parentPathID *uint, name string, status string, userID uint) error

	DeleteFile(ctx context.Context, fileID uint, userID uint) error
	DeleteDirectory(ctx context.Context, directoryID uint, userID uint) error

	CheckUserFileAccess(ctx context.Context, userID, fileID uint) (bool, error)
	CheckUserDirectoryAccess(ctx context.Context, userID, directoryID uint) (bool, error)

	UpdateFileStatus(ctx context.Context, fileID uint, status string, tx *gorm.DB) error
	WithTx(tx *gorm.DB) FileMetadataRepository // Метод для передачи транзакции
	GetDB() *gorm.DB
}
