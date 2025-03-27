// В internal/infrastructure/grpc/file_repository.go
package grpc

import (
	"context"
	"service-core/internal/domain"
	"service-core/internal/domain/interfaces"

	"gorm.io/gorm"
)

type FileRepositoryImpl struct {
	client *FileGRPCClient
}

func NewFileRepository(client *FileGRPCClient) interfaces.FileTreeRepository {
	return &FileRepositoryImpl{client: client}
}

func (r *FileRepositoryImpl) GetFileWithDirectory(ctx context.Context, fileID uint) (*domain.File, error) {
	return r.client.GetFileWithDirectory(ctx, fileID)
}

func (r *FileRepositoryImpl) GetDirectoriesWithFiles(ctx context.Context, isArchive bool, userID uint) ([]*domain.Directory, error) {
	panic("implement me!")
}
func (r *FileRepositoryImpl) GetFileInfo(ctx context.Context, fileID uint) (*domain.File, error) {
	panic("implement me!")
}

func (r *FileRepositoryImpl) CreateDirectory(ctx context.Context, parentPathID *uint, name string, status string, userID uint) error {
	panic("implement me!")
}
func (r *FileRepositoryImpl) CreateFile(ctx context.Context, directoryID uint, name string, status string, userID uint) error {
	panic("implement me!")
}

func (r *FileRepositoryImpl) DeleteDirectory(ctx context.Context, directoryID uint, userID uint) error {
	panic("implement me!")
}
func (r *FileRepositoryImpl) DeleteFile(ctx context.Context, fileID uint, userID uint) error {
	panic("implement me!")
}

func (r *FileRepositoryImpl) CheckUserDirectoryAccess(ctx context.Context, userID, directoryID uint) (bool, error) {
	panic("implement me!")
}
func (r *FileRepositoryImpl) CheckUserFileAccess(ctx context.Context, userID, fileID uint) (bool, error) {
	panic("implement me!")
}

func (r *FileRepositoryImpl) WithTx(tx *gorm.DB) interfaces.FileTreeRepository {
	panic("implement me!")
} // Метод для передачи транзакции
func (r *FileRepositoryImpl) GetDB() *gorm.DB {
	panic("implement me!")
}

func (r *FileRepositoryImpl) UpdateFileStatus(ctx context.Context, file *domain.File, tx *gorm.DB) error {
	panic("implement me!")

}
