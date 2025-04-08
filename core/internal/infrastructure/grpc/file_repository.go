// В internal/infrastructure/grpc/file_repository.go
package grpc

import (
	"context"
	"service-core/internal/domain"
	"service-core/internal/domain/interfaces"
)

type FileRepositoryImpl struct {
	client *FileGRPCClient
}

// TODO: убрать костыль

func NewFileService(client *FileGRPCClient) interfaces.FileService {
	return &FileRepositoryImpl{client: client}
}

func (r *FileRepositoryImpl) GetFileWithDirectory(ctx context.Context, fileID uint) (*domain.File, error) {
	return r.client.GetFileWithDirectory(ctx, fileID)
}

func (r *FileRepositoryImpl) UpdateFileStatus(ctx context.Context, fileID uint, status string) error {
	return r.client.UpdateFileStatus(ctx, fileID, status)
}

func (r *FileRepositoryImpl) GetFilesInfo(ctx context.Context, fileIDs []uint32) (map[uint32]string, error) {
	return r.client.GetFilesInfo(ctx, fileIDs)
}

func (r *FileRepositoryImpl) CheckWorkflow(ctx context.Context, workflowID uint) (bool, error) {
	return r.client.CheckWorkflow(ctx, workflowID)
}

func (r *FileRepositoryImpl) AssignWorkflow(ctx context.Context, workflowID uint, directoryIDs []uint32) error {
	return r.client.AssignWorkflow(ctx, workflowID, directoryIDs)
}

func (r *FileRepositoryImpl) DeleteUserRelations(ctx context.Context, userID uint) error {
	return r.client.DeleteUserRelations(ctx, userID)
}
