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

func NewFileMetadataRepository(db *Database) *FileMetadataRepository {
	return &FileMetadataRepository{db: db.db}
}

func (r *FileMetadataRepository) WithTx(tx *gorm.DB) interfaces.FileMetadataRepository {
	return &FileMetadataRepository{db: tx}
}

func (r *FileMetadataRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *FileMetadataRepository) GetFileByID(ctx context.Context, fileID uint) (*domain.File, error) {
	const op = "infrastructure.postgresrepo.file.GetFileByID"

	var file domain.File
	query := r.db.WithContext(ctx).
		Preload("Directory").
		Where("id = ?", fileID).
		First(&file)

	if query.Error != nil {
		if errors.Is(query.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%s: %w", op, domain.ErrFileNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, domain.ErrInternal)
	}

	return &file, nil
}

func (r *FileMetadataRepository) GetFilesByID(ctx context.Context, fileIDs []uint32, files *[]domain.File) error {
	return r.db.WithContext(ctx).
		Where("id IN (?)", fileIDs).
		Find(files).Error
}

func (r *FileMetadataRepository) GetFileInfo(ctx context.Context, fileID, userID uint) (*domain.File, error) {
	const op = "infrastructure.postgresrepo.file.GetFileInfo"

	type fileResult struct {
		domain.File
		AccessGranted bool `gorm:"column:access_granted"`
	}

	var result fileResult
	err := r.db.WithContext(ctx).
		Model(&domain.File{}).
		Select("files.*, CASE WHEN files.status = ? THEN TRUE ELSE EXISTS (SELECT 1 FROM user_files WHERE user_id = ? AND file_id = files.id) END as access_granted", "archive", userID).
		Where("files.id = ?", fileID).
		First(&result).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrFileNotFound
	} else if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if !result.AccessGranted {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrAccessDenied)
	}

	return &result.File, nil
}

func (r *FileMetadataRepository) CreateFile(ctx context.Context, directoryID uint, name string, status string, userID uint, minioKey string) error {
	const op = "infrastructure.postgresrepo.file.CreateFile"

	// Начинаем транзакцию
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		// Если транзакция не зафиксирована, откатываем
		if r := recover(); r != nil {
			tx.Rollback()
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	// Создаем новый файл
	newFile := domain.File{
		DirectoryID:    directoryID,
		Name:           name,
		Status:         status,
		MinioObjectKey: minioKey,
	}
	if err := tx.Create(&newFile).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, domain.ErrFileAlreadyExists)
	}

	userFile := domain.UserFile{
		UserID: userID,
		FileID: newFile.ID,
	}
	if err := tx.Create(&userFile).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *FileMetadataRepository) UpdateFileStatus(ctx context.Context, fileID uint, status string, tx *gorm.DB) error {
	const op = "infrastructure.postgresrepo.file.UpdateFileStatus"

	var file domain.File
	if err := tx.WithContext(ctx).First(&file, fileID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%s: %w", op, domain.ErrFileNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	// Обновляем статус файла
	file.Status = status

	// Используем явное условие для обновления
	return tx.WithContext(ctx).
		Model(&domain.File{}).   // Указываем модель, а не экземпляр
		Where("id = ?", fileID). // Указываем условие для обновления
		Update("status", status).
		Error
}

func (r *FileMetadataRepository) UpdateFile(ctx context.Context, file *domain.File) error {
	const op = "infrastructure.postgresrepo.file.UpdateFile"
	tx := r.db.WithContext(ctx).Begin()

	// Обновляем метаданные
	if err := tx.Model(&domain.File{}).
		Where("id = ?", file.ID).
		Updates(map[string]interface{}{
			"minio_object_key": file.MinioObjectKey,
			"version":          file.Version,
			"updated_at":       file.UpdatedAt,
		}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	return tx.Commit().Error
}

func (r *FileMetadataRepository) DeleteFile(ctx context.Context, fileID uint, userID uint) error {
	const op = "infrastructure.postgresrepo.file.DeleteFile"

	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Получаем файл и проверяем его существование
	var file domain.File
	if err := tx.First(&file, fileID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, domain.ErrFileNotFound)
	}

	// Объединённая проверка доступа: одновременно проверяем,
	// имеет ли пользователь доступ к директории и к файлу.
	var accessResult struct {
		HasDirectoryAccess bool `gorm:"column:has_directory_access"`
		HasFileAccess      bool `gorm:"column:has_file_access"`
	}
	accessQuery := `
		SELECT
			(SELECT COUNT(*) > 0 FROM user_directories WHERE user_id = ? AND directory_id = ?) AS has_directory_access,
			(SELECT COUNT(*) > 0 FROM user_files WHERE user_id = ? AND file_id = ?) AS has_file_access
	`
	if err := tx.Raw(accessQuery, userID, file.DirectoryID, userID, fileID).Scan(&accessResult).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}
	if !accessResult.HasDirectoryAccess || !accessResult.HasFileAccess {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, domain.ErrAccessDenied)
	}

	// Проверяем, что файл находится в статусе "draft"
	if file.Status != "draft" {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, domain.ErrCannotDeleteNonDraftFile)
	}

	if err := tx.Delete(&domain.File{}, fileID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := tx.Where("file_id = ?", fileID).Delete(&domain.UserFile{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *FileMetadataRepository) CheckUserFileAccess(ctx context.Context, userID, fileID uint) (bool, error) {
	const op = "infrastructure.postgresrepo.file.CheckUserFileAccess"

	var exists bool
	err := r.db.WithContext(ctx).
		Model(&domain.UserFile{}).
		Select("COUNT(*) > 0").
		Where("user_id = ? AND file_id = ?", userID, fileID).
		Scan(&exists).Error

	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return exists, nil
}
