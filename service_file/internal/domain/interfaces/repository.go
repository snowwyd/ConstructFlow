package interfaces

import (
	"context"
	"service-file/internal/domain"

	"gorm.io/gorm"
)

type FileTreeRepository interface {
	GetDirectoriesWithFiles(ctx context.Context, isArchive bool, userID uint) ([]*domain.Directory, error)
	GetFileInfo(ctx context.Context, fileID uint) (*domain.File, error)

	CreateDirectory(ctx context.Context, parentPathID *uint, name string, status string, userID uint) error
	CreateFile(ctx context.Context, directoryID uint, name string, status string, userID uint) error

	DeleteDirectory(ctx context.Context, directoryID uint, userID uint) error
	DeleteFile(ctx context.Context, fileID uint, userID uint) error

	CheckUserDirectoryAccess(ctx context.Context, userID, directoryID uint) (bool, error)
	CheckUserFileAccess(ctx context.Context, userID, fileID uint) (bool, error)

	WithTx(tx *gorm.DB) FileTreeRepository // Метод для передачи транзакции
	GetDB() *gorm.DB

	GetFileWithDirectory(ctx context.Context, fileID uint, tx *gorm.DB) (*domain.File, error)
	UpdateFileStatus(ctx context.Context, file *domain.File, tx *gorm.DB) error
}
