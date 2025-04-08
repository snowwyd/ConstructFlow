package interfaces

import (
	"context"
	"service-core/internal/domain"
)

type FileService interface {
	GetFileWithDirectory(ctx context.Context, fileID uint) (*domain.File, error)
	UpdateFileStatus(ctx context.Context, fileID uint, status string) error
	GetFilesInfo(ctx context.Context, fileIDs []uint32) (map[uint32]string, error)

	CheckWorkflow(ctx context.Context, workflowID uint) (bool, error)
	AssignWorkflow(ctx context.Context, workflowID uint, directoryIDs []uint32) error

	DeleteUserRelations(ctx context.Context, userID uint) error
}
