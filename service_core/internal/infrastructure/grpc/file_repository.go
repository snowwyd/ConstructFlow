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
