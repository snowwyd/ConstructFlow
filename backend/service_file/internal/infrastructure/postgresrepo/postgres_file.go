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
	const op = "infrastructure.postgresrepo.file.GetFilesByID"

	if err := r.db.WithContext(ctx).
		Where("id IN (?)", fileIDs).
		Find(files).Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
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

func (r *FileMetadataRepository) CreateFile(ctx context.Context, directoryID uint, name string, status string, userID uint, minioKey string, size int64, contentType string) error {
	const op = "infrastructure.postgresrepo.file.CreateFile"

	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	newFile := domain.File{
		DirectoryID:    directoryID,
		Name:           name,
		Status:         status,
		MinioObjectKey: minioKey,
		Size:           size,
		ContentType:    contentType,
		Version:        1,
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

	file.Status = status

	if err := tx.WithContext(ctx).
		Model(&domain.File{}).
		Where("id = ?", fileID).
		Update("status", status).
		Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *FileMetadataRepository) UpdateFile(ctx context.Context, file *domain.File) error {
	const op = "infrastructure.postgresrepo.file.UpdateFile"
	tx := r.db.WithContext(ctx).Begin()

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

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *FileMetadataRepository) DeleteFile(ctx context.Context, fileID uint, userID uint) error {
	const op = "infrastructure.postgresrepo.file.DeleteFile"

	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

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

func (fileMetadataRepo *FileMetadataRepository) DeleteUserRelations(ctx context.Context, userID uint) error {
	const op = "infrastructure.postgresrepo.file.DeleteUserRelations"

	tx := fileMetadataRepo.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	result := tx.Where("user_id = ?", userID).Delete(&domain.UserFile{})
	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, result.Error)
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, domain.ErrNoRelationsFound)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (fileMetadataRepo *FileMetadataRepository) CheckFilesExist(ctx context.Context, fileIDs []uint) (bool, error) {
	const op = "infrastructure.postgresrepo.directory.CheckFilesExist"

	if len(fileIDs) == 0 {
		return false, fmt.Errorf("%s: empty user list", op)
	}

	var count int64
	err := fileMetadataRepo.db.WithContext(ctx).
		Model(&domain.File{}).
		Where("id IN (?)", fileIDs).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return count == int64(len(fileIDs)), nil
}

func (r *FileMetadataRepository) UpdateUserFileRelations(ctx context.Context, userID uint, fileIDs []uint) error {
	const op = "infrastructure.postgresrepo.file.UpdateUserFileRelations"

	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Where("user_id = ?", userID).Delete(&domain.UserFile{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: failed to delete user file relations: %w", op, err)
	}

	var newRelations []domain.UserFile
	for _, fileID := range fileIDs {
		newRelations = append(newRelations, domain.UserFile{
			UserID: userID,
			FileID: fileID,
		})
	}

	if len(newRelations) > 0 {
		if err := tx.Create(&newRelations).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("%s: failed to create user file relations: %w", op, err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
