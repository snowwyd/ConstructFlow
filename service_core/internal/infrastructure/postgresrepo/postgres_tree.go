package postgresrepo

import (
	"context"
	"errors"
	"fmt"
	"service-core/internal/domain"
	"service-core/internal/domain/interfaces"

	"gorm.io/gorm"
)

type FileTreeRepository struct {
	db *gorm.DB
}

func NewFileTreeRepository(db *Database) *FileTreeRepository {
	return &FileTreeRepository{db: db.db}
}

// func (r *FileTreeRepository) WithTx(tx *gorm.DB) interfaces.FileTreeRepository {
// 	return &FileTreeRepository{db: tx}
// }

func (r *FileTreeRepository) GetDB() *gorm.DB {
	return r.db
}
func (r *FileTreeRepository) WithTx(tx *gorm.DB) interfaces.FileTreeRepository {
	return &FileTreeRepository{db: tx}
}

// GetDirectoriesWithFiles отдает список всех директорий и файлов, соответствующих входным условиям
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

	if isArchive {
		query = query.Where("status = ?", "archive")
	} else {
		query = query.Joins("JOIN user_directories ON user_directories.directory_id = directories.id").
			Where("user_directories.user_id = ?", userID).
			Where("directories.status != ?", "archive")
	}

	if err := query.Find(&directories).Error; err != nil {
		return nil, fmt.Errorf("failed to get directories: %w", domain.ErrInternal)
	}

	// Проверяем случай, когда для активных директорий нет доступных записей
	if !isArchive && len(directories) == 0 {
		return nil, domain.ErrAccessDenied
	}

	return directories, nil
}

func (r *FileTreeRepository) GetFileInfo(ctx context.Context, fileID uint) (*domain.File, error) {

	var file domain.File
	err := r.db.WithContext(ctx).First(&file, fileID).Error

	// Обрабатываем возможные ошибки
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, domain.ErrFileNotFound
	case err != nil:
		return nil, fmt.Errorf("failed to retrieve file: %w", domain.ErrInternal)
	}

	return &file, nil
}

// CreateDirectory записывает новую сущность Directory в БД
func (r *FileTreeRepository) CreateDirectory(ctx context.Context, parentPathID *uint, name string, status string, userID uint) error {
	tx := r.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	newDir := domain.Directory{
		ParentPathID: parentPathID,
		Name:         name,
		Status:       status,
	}

	if err := tx.Create(&newDir).Error; err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Создаем связь с пользователем
	userDir := domain.UserDirectory{
		UserID:      userID,
		DirectoryID: newDir.ID,
	}

	if err := tx.Create(&userDir).Error; err != nil {
		return fmt.Errorf("failed to create user-directory relation: %w", err)
	}

	tx.Commit()
	return nil
}

// CreateFile записывает новую сущность File в БД
func (r *FileTreeRepository) CreateFile(ctx context.Context, directoryID uint, name string, status string, userID uint) error {
	tx := r.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	newFile := domain.File{
		DirectoryID: directoryID,
		Name:        name,
		Status:      status,
	}

	if err := tx.Create(&newFile).Error; err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	// Создаем связь с пользователем
	userFile := domain.UserFile{
		UserID: userID,
		FileID: newFile.ID,
	}

	if err := tx.Create(&userFile).Error; err != nil {
		return fmt.Errorf("failed to create user-file relation: %w", err)
	}

	tx.Commit()
	return nil
}

// DeleteDirectory рекурсивно удаляет все файлы и директории внутри указанной
func (r *FileTreeRepository) DeleteDirectory(ctx context.Context, directoryID uint, userID uint) error {
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

// DeleteFile удаляет файл из БД
func (r *FileTreeRepository) DeleteFile(ctx context.Context, fileID uint, userID uint) error {
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

// CheckUserDirectoryAccess проверяет доступ пользователя к директории
func (r *FileTreeRepository) CheckUserDirectoryAccess(ctx context.Context, userID, directoryID uint) (bool, error) {

	// Проверка существования директории
	var directoryExists bool
	err := r.db.WithContext(ctx).
		Model(&domain.Directory{}).
		Select("COUNT(*) > 0").
		Where("id = ?", directoryID).
		Find(&directoryExists).Error
	if err != nil {
		return false, fmt.Errorf("directory existence check failed: %w", domain.ErrInternal)
	}
	if !directoryExists {
		return false, domain.ErrDirectoryNotFound
	}

	// Проверка связи пользователя с директорией
	var count int64
	err = r.db.WithContext(ctx).
		Model(&domain.UserDirectory{}).
		Where("user_id = ? AND directory_id = ?", userID, directoryID).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("user access check failed: %w", domain.ErrInternal)
	}

	return count > 0, nil
}

// CheckUserFileAccess проверяет доступ пользователя к файлу
func (r *FileTreeRepository) CheckUserFileAccess(ctx context.Context, userID, fileID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.File{}). // Используем модель File для проверки статуса
		Joins("JOIN user_files ON user_files.file_id = files.id").
		Where("user_files.user_id = ?", userID).
		Where("files.id = ?", fileID).
		Where("files.status = ?", "draft"). // Добавляем проверку статуса
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetFileWithDirectory получает файл с директорией
// Кастомные ошибки: ErrFileNotFound
func (r *FileTreeRepository) GetFileWithDirectory(ctx context.Context, fileID uint, tx *gorm.DB) (*domain.File, error) {
	// Проверка валидности входных данных
	if fileID == 0 {
		return nil, domain.ErrInvalidID
	}

	var file domain.File
	query := tx.WithContext(ctx).
		Preload("Directory").
		Where("id = ?", fileID).
		First(&file)

	// Обработка ошибок
	if query.Error != nil {
		if errors.Is(query.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrFileNotFound
		}
		// Все остальные ошибки считаются внутренними
		return nil, fmt.Errorf("failed to get file with directory: %w", domain.ErrInternal)
	}

	return &file, nil
}

// UpdateFileStatus обновляет статус файла
func (r *FileTreeRepository) UpdateFileStatus(ctx context.Context, file *domain.File, tx *gorm.DB) error {
	return tx.WithContext(ctx).
		Model(file).
		Update("status", file.Status).
		Error
}
