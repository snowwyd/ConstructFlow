package postgresrepo

import (
	"backend/internal/domain"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type FileTreeRepository struct {
	db *gorm.DB
}

func NewFileTreeRepository(db *Database) *FileTreeRepository {
	return &FileTreeRepository{db: db.db}
}

func (r *FileTreeRepository) GetDirectoriesWithFiles(ctx context.Context, isArchive bool, userID uint) ([]*domain.Directory, error) {
	var directories []*domain.Directory

	query := r.db.WithContext(ctx).
		Preload("Files", func(db *gorm.DB) *gorm.DB {
			if isArchive {
				return db.Where("status = ?", "archive")
			}
			return db.Where("status != ?", "archive")
		}).
		Where("status != ?", "deleted") // Пример фильтрации удаленных директорий

	if err := query.Find(&directories).Error; err != nil {
		return nil, fmt.Errorf("failed to get directories: %w", err)
	}

	return directories, nil
}

func (r *FileTreeRepository) CreateDirectory(ctx context.Context, parentPathID *uint, name string, status string) (uint, error) {
	newDir := domain.Directory{
		ParentPathID: parentPathID,
		Name:         name,
		Status:       status,
	}

	if err := r.db.WithContext(ctx).Create(&newDir).Error; err != nil {
		return 0, fmt.Errorf("failed to create directory: %w", err)
	}

	return newDir.ID, nil
}

func (r *FileTreeRepository) CreateFile(ctx context.Context, directoryID uint, name string, status string) (uint, error) {
	newFile := domain.File{
		DirectoryID: directoryID,
		Name:        name,
		Status:      status,
	}

	if err := r.db.WithContext(ctx).Create(&newFile).Error; err != nil {
		return 0, fmt.Errorf("failed to create file: %w", err)
	}

	return newFile.ID, nil
}
