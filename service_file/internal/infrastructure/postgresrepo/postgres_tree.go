package postgresrepo

import (
	"context"
	"errors"
	"fmt"
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
	const op = "infrastructure.postgres.tree.GetFileTree"

	var directories []domain.Directory

	query := r.db.WithContext(ctx).
		Select("directories.*").
		Preload("Files", func(db *gorm.DB) *gorm.DB {
			if isArchive {
				return db.Where("status = ?", "archive")
			}
			return db.Where("status != ?", "archive")
		})

	if isArchive {
		query = query.Where("status = ?", "archive")
	} else {
		query = query.Joins("JOIN user_directories ON user_directories.directory_id = directories.id").
			Where("user_directories.user_id = ?", userID).
			Where("directories.status != ?", "archive")
	}

	if err := query.Find(&directories).Error; err != nil {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrInternal)
	}

	// Проверяем случай, когда для активных директорий нет доступных записей
	if !isArchive && len(directories) == 0 {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrAccessDenied)
	}

	return directories, nil
}

func (r *FileMetadataRepository) GetFileByID(ctx context.Context, fileID uint) (*domain.File, error) {
	const op = "infrastructure.postgres.tree.GetFileByID"

	var file domain.File
	err := r.db.WithContext(ctx).First(&file, fileID).Error

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, fmt.Errorf("%s: %w", op, domain.ErrFileNotFound)
	case err != nil:
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &file, nil
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
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.File{}). // Используем модель File для проверки статуса
		Joins("JOIN user_files ON user_files.file_id = files.id").
		Where("user_files.user_id = ?", userID).
		Where("files.id = ?", fileID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *FileMetadataRepository) CheckUserDirectoryAccess(ctx context.Context, userID, directoryID uint) (bool, error) {
	panic("implement me!")
}

func (r *FileMetadataRepository) UpdateFileStatus(ctx context.Context, file *domain.File, tx *gorm.DB) error {
	panic("implement me!")
}
