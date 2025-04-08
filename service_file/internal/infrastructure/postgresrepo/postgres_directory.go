package postgresrepo

import (
	"context"
	"errors"
	"fmt"
	"service-file/internal/domain"
	"service-file/internal/domain/interfaces"

	"gorm.io/gorm"
)

type DirectoryRepository struct {
	db *gorm.DB
}

func NewDirectoryRepository(db *Database) *DirectoryRepository {
	return &DirectoryRepository{db: db.db}
}

func (r *DirectoryRepository) WithTx(tx *gorm.DB) interfaces.DirectoryRepository {
	return &DirectoryRepository{db: tx}
}

func (r *DirectoryRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *DirectoryRepository) GetFileTree(ctx context.Context, isArchive bool, userID uint) ([]domain.Directory, error) {
	const op = "infrastructure.postgresrepo.directory.GetFileTree"

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
		query = query.Where("directories.status = ?", "archive")
	} else {
		query = query.Joins("JOIN user_directories ON user_directories.directory_id = directories.id").
			Where("user_directories.user_id = ?", userID).
			Where("directories.status != ?", "archive")
	}

	if err := query.Find(&directories).Error; err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if !isArchive && len(directories) == 0 {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrAccessDenied)
	}

	return directories, nil
}

func (r *DirectoryRepository) CreateDirectory(ctx context.Context, parentPathID *uint, name string, status string, userID uint) error {
	const op = "infrastructure.postgresrepo.directory.CreateDirectory"

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
	newDirectory := domain.Directory{
		ParentPathID: parentPathID,
		Name:         name,
		Status:       status,
	}
	if err := tx.Create(&newDirectory).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, domain.ErrDirectoryAlreadyExists)
	}

	userDirectory := domain.UserDirectory{
		UserID:      userID,
		DirectoryID: newDirectory.ID,
	}
	if err := tx.Create(&userDirectory).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *DirectoryRepository) DeleteDirectory(ctx context.Context, directoryID, userID uint) error {
	const op = "infrastructure.postgresrepo.directory.DeleteDirectory"
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Получаем файл и проверяем его существование
	var directory domain.Directory
	if err := tx.First(&directory, directoryID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, domain.ErrDirectoryNotFound)
	}

	// Проверяем, существует ли директория и имеет ли пользователь к ней доступ
	var hasAccess bool
	err := tx.Raw(`
		SELECT COUNT(*) > 0
		FROM user_directories ud
		JOIN directories d ON ud.directory_id = d.id
		WHERE ud.user_id = ? AND d.id = ?
	`, userID, directoryID).Scan(&hasAccess).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}
	if !hasAccess {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, domain.ErrAccessDenied)
	}

	// Проверяем, есть ли в текущей директории или её потомках файлы со статусом != "draft"
	var hasNonDraftFiles bool
	err = tx.Raw(`
		SELECT EXISTS (
			SELECT 1 FROM files 
			WHERE directory_id IN (
				WITH RECURSIVE subdirs AS (
					SELECT id FROM directories WHERE id = ?
					UNION ALL
					SELECT d.id FROM directories d 
					JOIN subdirs s ON d.parent_path_id = s.id
				) SELECT id FROM subdirs
			) AND status != 'draft'
		)
	`, directoryID).Scan(&hasNonDraftFiles).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}
	if hasNonDraftFiles {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, domain.ErrDirectoryContainsNonDraftFiles)
	}

	// Удаляем файлы и их связи одним запросом
	err = tx.Exec(`
		DELETE FROM files 
		WHERE directory_id IN (
			WITH RECURSIVE subdirs AS (
				SELECT id FROM directories WHERE id = ?
				UNION ALL
				SELECT d.id FROM directories d 
				JOIN subdirs s ON d.parent_path_id = s.id
			) SELECT id FROM subdirs
		)
	`, directoryID).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	// Удаляем директории и их связи одним запросом
	err = tx.Exec(`
		DELETE FROM directories 
		WHERE id IN (
			WITH RECURSIVE subdirs AS (
				SELECT id FROM directories WHERE id = ?
				UNION ALL
				SELECT d.id FROM directories d 
				JOIN subdirs s ON d.parent_path_id = s.id
			) SELECT id FROM subdirs
		)
	`, directoryID).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	err = tx.Exec(`
		DELETE FROM user_directories 
		WHERE directory_id IN (
			WITH RECURSIVE subdirs AS (
				SELECT id FROM directories WHERE id = ?
				UNION ALL
				SELECT d.id FROM directories d 
				JOIN subdirs s ON d.parent_path_id = s.id
			) SELECT id FROM subdirs
		)
	`, directoryID).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *DirectoryRepository) CheckUserDirectoryAccess(ctx context.Context, userID, directoryID uint) (bool, error) {
	const op = "infrastructure.postgresrepo.directory.CheckUserDirectoryAccess"

	var dir domain.Directory
	if err := r.db.WithContext(ctx).First(&dir, directoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, fmt.Errorf("%s: %w", op, domain.ErrDirectoryNotFound)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	var exists bool
	err := r.db.WithContext(ctx).
		Model(&domain.UserDirectory{}).
		Select("COUNT(*) > 0").
		Where("user_id = ? AND directory_id = ?", userID, directoryID).
		Scan(&exists).Error

	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return exists, nil
}

func (directoryRepository *DirectoryRepository) CheckWorkflow(ctx context.Context, workflowID uint) (bool, error) {
	const op = "infrastructure.postgresrepo.directory.CheckWorkflowExists"

	var count int64
	err := directoryRepository.db.WithContext(ctx).
		Model(&domain.Directory{}).
		Where("workflow_id = ?", workflowID).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return count > 0, nil
}

func (directoryRepository *DirectoryRepository) DeleteUserRelations(ctx context.Context, userID uint) error {
	const op = "infrastructure.postgresrepo.directory.DeleteUserRelations"

	tx := directoryRepository.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	result := tx.Where("user_id = ?", userID).Delete(&domain.UserDirectory{})
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
