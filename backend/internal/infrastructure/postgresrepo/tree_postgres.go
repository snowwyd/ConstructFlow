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
		Select("directories.*").
		Preload("Files", func(db *gorm.DB) *gorm.DB {
			if isArchive {
				return db.Where("status = ?", "archive")
			}
			return db.Where("status != ?", "archive")
		})

	// Условия в зависимости от isArchive
	if isArchive {
		// Для архивных директорий берем все без привязки к пользователю
		query = query.Where("status = ?", "archive")
	} else {
		// Для активных директорий проверяем связь с пользователем
		query = query.Joins("JOIN user_directories ON user_directories.directory_id = directories.id").
			Where("user_directories.user_id = ?", userID).
			Where("directories.status != ?", "archive")
	}

	if err := query.Find(&directories).Error; err != nil {
		return nil, fmt.Errorf("failed to get directories: %w", err)
	}

	return directories, nil
}

func (r *FileTreeRepository) GetFileInfo(ctx context.Context, fileID uint) (*domain.File, error) {
	var file domain.File
	if err := r.db.WithContext(ctx).First(&file, fileID).Error; err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", domain.ErrFileNotFound)
	}
	return &file, nil
}

func (r *FileTreeRepository) CreateDirectory(ctx context.Context, parentPathID *uint, name string, status string, userID uint) (uint, error) {
	tx := r.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	newDir := domain.Directory{
		ParentPathID: parentPathID,
		Name:         name,
		Status:       status,
	}

	if err := tx.Create(&newDir).Error; err != nil {
		return 0, fmt.Errorf("failed to create directory: %w", err)
	}

	// Создаем связь с пользователем
	userDir := domain.UserDirectory{
		UserID:      userID,
		DirectoryID: newDir.ID,
	}

	if err := tx.Create(&userDir).Error; err != nil {
		return 0, fmt.Errorf("failed to create user-directory relation: %w", err)
	}

	tx.Commit()
	return newDir.ID, nil
}

func (r *FileTreeRepository) CreateFile(ctx context.Context, directoryID uint, name string, status string, userID uint) (uint, error) {
	tx := r.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	newFile := domain.File{
		DirectoryID: directoryID,
		Name:        name,
		Status:      status,
	}

	if err := tx.Create(&newFile).Error; err != nil {
		return 0, fmt.Errorf("failed to create file: %w", err)
	}

	// Создаем связь с пользователем
	userFile := domain.UserFile{
		UserID: userID,
		FileID: newFile.ID,
	}

	if err := tx.Create(&userFile).Error; err != nil {
		return 0, fmt.Errorf("failed to create user-file relation: %w", err)
	}

	tx.Commit()
	return newFile.ID, nil
}

func (r *FileTreeRepository) DeleteDirectory(ctx context.Context, directoryID uint, userID uint) error {
	tx := r.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Проверка существования директории
	var dir domain.Directory
	if err := tx.First(&dir, directoryID).Error; err != nil {
		return fmt.Errorf("directory not found: %w", err)
	}

	// Проверка доступа пользователя к директории
	hasAccess, err := r.CheckUserDirectoryAccess(ctx, userID, directoryID)
	if err != nil {
		return fmt.Errorf("access check failed: %w", err)
	}
	if !hasAccess {
		return domain.ErrAccessDenied
	}

	// Проверка доступа к родительской директории
	if dir.ParentPathID != nil {
		parentAccess, err := r.CheckUserDirectoryAccess(ctx, userID, *dir.ParentPathID)
		if err != nil {
			return fmt.Errorf("parent access check failed: %w", err)
		}
		if !parentAccess {
			return domain.ErrAccessDenied
		}
	}

	// Удаление всех файлов, связанных с этой директорией
	if err := tx.Where("directory_id = ?", directoryID).Delete(&domain.File{}).Error; err != nil {
		return fmt.Errorf("failed to delete files in directory: %w", err)
	}

	// Удаление связей файлов с пользователями
	if err := tx.Where("file_id IN (SELECT id FROM files WHERE directory_id = ?)", directoryID).Delete(&domain.UserFile{}).Error; err != nil {
		return fmt.Errorf("failed to delete user-file relations: %w", err)
	}

	// Удаление директории
	if err := tx.Where("id = ?", directoryID).Delete(&domain.Directory{}).Error; err != nil {
		return fmt.Errorf("failed to delete directory: %w", err)
	}

	// Удаление связей пользователя с директорией
	if err := tx.Where("directory_id = ?", directoryID).Delete(&domain.UserDirectory{}).Error; err != nil {
		return fmt.Errorf("failed to delete user-directory relations: %w", err)
	}

	tx.Commit()
	return nil
}

func (r *FileTreeRepository) DeleteFile(ctx context.Context, fileID uint, userID uint) error {
	tx := r.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Проверка существования файла
	var file domain.File
	if err := tx.First(&file, fileID).Error; err != nil {
		return fmt.Errorf("file not found: %w", err)
	}

	// Проверка доступа пользователя к директории файла
	hasAccess, err := r.CheckUserDirectoryAccess(ctx, userID, file.DirectoryID)
	if err != nil {
		return fmt.Errorf("access check failed: %w", err)
	}
	if !hasAccess {
		return domain.ErrAccessDenied
	}

	// Удаление файла и связей
	if err := tx.Where("id = ?", fileID).Delete(&domain.File{}).Error; err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	// Удаление связей из user_files
	if err := tx.Where("file_id = ?", fileID).Delete(&domain.UserFile{}).Error; err != nil {
		return fmt.Errorf("failed to delete user-file relations: %w", err)
	}

	tx.Commit()
	return nil
}

func (r *FileTreeRepository) CheckUserDirectoryAccess(ctx context.Context, userID, directoryID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.UserDirectory{}).
		Where("user_id = ? AND directory_id = ?", userID, directoryID).
		Count(&count).Error
	return count > 0, err
}

func (r *FileTreeRepository) CheckUserFileAccess(ctx context.Context, userID, fileID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.UserFile{}).
		Where("user_id = ? AND file_id = ?", userID, fileID).
		Count(&count).Error
	return count > 0, err
}
