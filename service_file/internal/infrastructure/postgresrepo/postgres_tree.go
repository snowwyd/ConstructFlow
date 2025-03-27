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

func (r *FileMetadataRepository) GetDB() *gorm.DB {
	return r.db
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

	if !isArchive && len(directories) == 0 {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrAccessDenied)
	}

	return directories, nil
}

func (r *FileMetadataRepository) GetFileByID(ctx context.Context, fileID uint) (*domain.File, error) {
	const op = "infrastructure.postgres.tree.GetFileByID"

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

func (r *FileMetadataRepository) GetFileInfo(ctx context.Context, fileID uint) (*domain.File, error) {

	var file domain.File
	err := r.db.WithContext(ctx).First(&file, fileID).Error

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, domain.ErrFileNotFound
	case err != nil:
		return nil, fmt.Errorf("failed to retrieve file: %w", domain.ErrInternal)
	}

	return &file, nil
}

func (r *FileMetadataRepository) GetDirectoryByID(ctx context.Context, directoryID uint) (*domain.Directory, error) {
	panic("implement me!")
}

func (r *FileMetadataRepository) CreateFile(ctx context.Context, directoryID uint, name string, status string, userID uint) error {
	const op = "infrastructure.postgres.tree.CreateFile"

	// Начинаем транзакцию
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}

	// Проверяем уникальность имени файла в рамках директории
	var exists bool
	err := tx.Model(&domain.File{}).
		Select("COUNT(*) > 0").
		Where("directory_id = ? AND name = ?", directoryID, name).
		Scan(&exists).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: failed to check file uniqueness: %w", op, err)
	}
	if exists {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, domain.ErrFileAlreadyExists)
	}

	// Создаем новый файл
	newFile := domain.File{
		DirectoryID: directoryID,
		Name:        name,
		Status:      status,
	}

	if err := tx.Create(&newFile).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: failed to create file: %w", op, err)
	}

	// Создаем связь с пользователем через один запрос
	if err := tx.Model(&domain.UserFile{}).
		Create(map[string]interface{}{
			"user_id": userID,
			"file_id": newFile.ID,
		}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: failed to create user-file relation: %w", op, err)
	}

	// Фиксируем транзакцию
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return nil
}

func (r *FileMetadataRepository) CreateDirectory(ctx context.Context, parentPathID *uint, name string, status string, userID uint) error {
	const op = "infrastructure.postgres.file.CreateDirectory"

	// Начинаем транзакцию
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}

	// Проверяем уникальность имени директории
	var exists bool
	err := tx.Model(&domain.Directory{}).
		Select("COUNT(*) > 0").
		Where("parent_path_id = ? AND name = ?", parentPathID, name).
		Scan(&exists).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: failed to check directory uniqueness: %w", op, err)
	}
	if exists {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, domain.ErrDirectoryAlreadyExists)
	}

	// Создаем новую директорию
	newDir := domain.Directory{
		ParentPathID: parentPathID,
		Name:         name,
		Status:       status,
	}

	if err := tx.Create(&newDir).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: failed to create directory: %w", op, err)
	}

	// Создаем связь с пользователем через один запрос
	if err := tx.Model(&domain.UserDirectory{}).
		Create(map[string]interface{}{
			"user_id":      userID,
			"directory_id": newDir.ID,
		}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: failed to create user-directory relation: %w", op, err)
	}

	// Фиксируем транзакцию
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return nil
}

func (r *FileMetadataRepository) DeleteFile(ctx context.Context, fileID uint, userID uint) error {
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Проверка существования файла
	var file domain.File
	if err := tx.First(&file, fileID).Error; err != nil {
		return domain.ErrFileNotFound
	}

	// Проверка доступа пользователя к директории файла
	hasAccess, err := r.CheckUserDirectoryAccess(ctx, userID, file.DirectoryID)
	if err != nil {
		return fmt.Errorf("access check failed: %w", err)
	}
	if !hasAccess {
		return domain.ErrAccessDenied
	}

	// Проверка доступа пользователя к файлу
	hasAccess, err = r.CheckUserFileAccess(ctx, userID, fileID)
	if err != nil {
		return fmt.Errorf("access check failed: %w", err)
	}
	if !hasAccess {
		return domain.ErrAccessDenied
	}

	// Проверка статуса файла
	if file.Status != "draft" {
		return domain.ErrCannotDeleteNonDraftFile
	}

	// Удаление файла
	if err := tx.Where("id = ?", fileID).Delete(&domain.File{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete file: %w", err)
	}

	// Удаление связей из user_files
	if err := tx.Where("file_id = ?", fileID).Delete(&domain.UserFile{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete user-file relations: %w", err)
	}

	tx.Commit()
	return nil
}

func (r *FileMetadataRepository) DeleteDirectory(ctx context.Context, directoryID uint, userID uint) error {
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Рекурсивная функция для проверки наличия файлов со статусом "draft"
	var hasDraftFilesRecursive func(dirID uint) (bool, error)
	hasDraftFilesRecursive = func(dirID uint) (bool, error) {
		// Проверяем файлы в текущей директории
		var hasDraftFiles bool
		err := tx.Model(&domain.File{}).
			Where("directory_id = ?", dirID).
			Where("status != ?", "draft").
			Limit(1).
			Select("COUNT(*) > 0").
			Scan(&hasDraftFiles).Error
		if err != nil {
			return false, fmt.Errorf("failed to check for draft files: %w", err)
		}
		if hasDraftFiles {
			return true, nil
		}

		// Находим все дочерние директории
		var childDirs []domain.Directory
		if err := tx.Where("parent_path_id = ?", dirID).Find(&childDirs).Error; err != nil {
			return false, fmt.Errorf("failed to find child directories: %w", err)
		}

		// Рекурсивно проверяем каждую дочернюю директорию
		for _, childDir := range childDirs {
			hasDrafts, err := hasDraftFilesRecursive(childDir.ID)
			if err != nil {
				return false, err
			}
			if hasDrafts {
				return true, nil
			}
		}

		return false, nil
	}

	// Проверка прав доступа к корневой директории
	hasAccess, err := r.CheckUserDirectoryAccess(ctx, userID, directoryID)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, domain.ErrDirectoryNotFound) {
			return domain.ErrDirectoryNotFound
		}
		return fmt.Errorf("access check failed: %w", err)
	}
	if !hasAccess {
		tx.Rollback()
		return domain.ErrAccessDenied
	}

	// Проверка наличия файлов со статусом "draft" в иерархии
	hasDraftFiles, err := hasDraftFilesRecursive(directoryID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to check for draft files in directory hierarchy: %w", err)
	}
	if hasDraftFiles {
		tx.Rollback()
		return domain.ErrDirectoryContainsNonDraftFiles
	}

	// Рекурсивное удаление директорий и файлов
	var deleteDirectoryRecursive func(dirID uint) error
	deleteDirectoryRecursive = func(dirID uint) error {
		// Удаление всех файлов в текущей директории
		if err := tx.Where("directory_id = ?", dirID).Delete(&domain.File{}).Error; err != nil {
			return fmt.Errorf("failed to delete files in directory: %w", err)
		}

		// Удаление связей файлов с пользователями
		if err := tx.Where("file_id IN (SELECT id FROM files WHERE directory_id = ?)", dirID).Delete(&domain.UserFile{}).Error; err != nil {
			return fmt.Errorf("failed to delete user-file relations: %w", err)
		}

		// Найти все дочерние директории
		var childDirs []domain.Directory
		if err := tx.Where("parent_path_id = ?", dirID).Find(&childDirs).Error; err != nil {
			return fmt.Errorf("failed to find child directories: %w", err)
		}

		// Рекурсивно удалить каждую дочернюю директорию
		for _, childDir := range childDirs {
			if err := deleteDirectoryRecursive(childDir.ID); err != nil {
				return err
			}
		}

		// Удаление текущей директории
		if err := tx.Where("id = ?", dirID).Delete(&domain.Directory{}).Error; err != nil {
			return fmt.Errorf("failed to delete directory: %w", err)
		}

		// Удаление связей пользователя с текущей директорией
		if err := tx.Where("directory_id = ?", dirID).Delete(&domain.UserDirectory{}).Error; err != nil {
			return fmt.Errorf("failed to delete user-directory relations: %w", err)
		}

		return nil
	}

	// Вызываем рекурсивную функцию для удаления
	if err := deleteDirectoryRecursive(directoryID); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *FileMetadataRepository) CheckUserFileAccess(ctx context.Context, userID, fileID uint) (bool, error) {
	var exists bool
	err := r.db.WithContext(ctx).
		Model(&domain.UserFile{}).
		Select("COUNT(*) > 0").
		Where("user_id = ? AND file_id = ?", userID, fileID).
		Scan(&exists).Error

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *FileMetadataRepository) CheckUserDirectoryAccess(ctx context.Context, userID, directoryID uint) (bool, error) {
	var exists bool
	err := r.db.WithContext(ctx).
		Model(&domain.UserDirectory{}).
		Select("COUNT(*) > 0").
		Where("user_id = ? AND directory_id = ?", userID, directoryID).
		Scan(&exists).Error

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *FileMetadataRepository) UpdateFileStatus(ctx context.Context, fileID uint, status string, tx *gorm.DB) error {
	// Проверяем существование файла
	var file domain.File
	if err := tx.WithContext(ctx).First(&file, fileID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrFileNotFound
		}
		return err
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
