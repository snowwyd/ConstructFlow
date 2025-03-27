package interfaces

import (
	"context"
	"service-core/internal/domain"
)

type FileService interface {
	GetFileWithDirectory(ctx context.Context, fileID uint) (*domain.File, error)
	UpdateFileStatus(ctx context.Context, fileID uint, status string) error
}
