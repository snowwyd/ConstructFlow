package postgresrepo

import (
	"context"
	"service-file/internal/domain"
	"service-file/internal/domain/interfaces"

	"gorm.io/gorm"
)

type FileMetadataRepository struct {
	db *gorm.DB
}

func NewFileTreeRepository(db *Database) *FileMetadataRepository {
	return &FileMetadataRepository{db: db.db}
}

func (r *FileMetadataRepository) WithTx(tx *gorm.DB) interfaces.FileMetadataRepository {
	return &FileMetadataRepository{db: tx}
}

func (r *FileMetadataRepository) GetFileTree(ctx context.Context, isArchive bool, userID uint) ([]domain.Directory, error) {
	// имитация взятия данных из репозитория
	return []domain.Directory{domain.Directory{
		Name:         "test",
		Status:       "draft",
		Version:      1,
		ParentPathID: nil,
		Files:        []domain.File{domain.File{}},
	}}, nil
}

func (r *FileMetadataRepository) GetFileByID(ctx context.Context, fileID uint) (*domain.File, error) {
	panic("implement me!")
}

func (r *FileMetadataRepository) GetDirectoryByID(ctx context.Context, directoryID uint) (*domain.Directory, error) {
	panic("implement me!")
}
func (r *FileMetadataRepository) CreateFile(ctx context.Context, directoryID uint, name string, status string, userID uint) error {
	panic("implement me!")
}

func (r *FileMetadataRepository) CreateDirectory(ctx context.Context, parentPathID *uint, name string, status string, userID uint) error {
	panic("implement me!")
}

func (r *FileMetadataRepository) DeleteFile(ctx context.Context, fileID uint, userID uint) error {
	panic("implement me!")
}

func (r *FileMetadataRepository) DeleteDirectory(ctx context.Context, directoryID uint, userID uint) error {
	panic("implement me!")
}

func (r *FileMetadataRepository) CheckUserFileAccess(ctx context.Context, userID, fileID uint) (bool, error) {
	panic("implement me!")
}

func (r *FileMetadataRepository) CheckUserDirectoryAccess(ctx context.Context, userID, directoryID uint) (bool, error) {
	panic("implement me!")
}

func (r *FileMetadataRepository) UpdateFileStatus(ctx context.Context, file *domain.File, tx *gorm.DB) error {
	panic("implement me!")
}
