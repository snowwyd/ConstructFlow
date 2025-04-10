package postgresrepo

import (
	"context"
	"errors"
	"fmt"
	"service-file/internal/domain"
	"service-file/internal/domain/interfaces"
	"sort"

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

func (r *DirectoryRepository) GetFileTreeWorking(ctx context.Context, userID uint) ([]domain.Directory, error) {
	const op = "infrastructure.postgresrepo.directory.GetFileTreeWorking"

	var directories []domain.Directory

	query := r.db.WithContext(ctx).
		Select("directories.*").
		Preload("Files", func(db *gorm.DB) *gorm.DB {
			return db.Joins("JOIN user_files AS uf ON uf.file_id = files.id").
				Where("status != ? AND uf.user_id = ?", "archive", userID)
		})

	query = query.Joins("JOIN user_directories ON user_directories.directory_id = directories.id").
		Where("user_directories.user_id = ?", userID).
		Where("directories.status != ?", "archive")

	if err := query.Find(&directories).Error; err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return directories, nil
}

func (r *DirectoryRepository) GetFileTreeArchive(ctx context.Context, userID uint) ([]domain.Directory, error) {
	const op = "infrastructure.postgresrepo.directory.GetFileTreeArchive"

	var directories []domain.Directory

	query := r.db.WithContext(ctx).
		Select("directories.*").
		Preload("Files", func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", "archive")
		})

	query = query.Where("directories.status = ? ", "archive")

	if err := query.Find(&directories).Error; err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return directories, nil
}

func (r *DirectoryRepository) GetUserFileTree(ctx context.Context, userID uint) ([]domain.DirectoryUserResponse, error) {
	const op = "infrastructure.postgresrepo.directory.GerUserFileTree"

	query := `
        SELECT 
            d.id AS directory_id,
            d.name AS name_folder,
            d.parent_path_id AS parent_path_id,
            (ud.user_id IS NOT NULL) AS user_has_access,
            f.id AS file_id,
            f.name AS name_file,
            f.directory_id AS file_directory_id,
            (uf.user_id IS NOT NULL) AS file_user_has_access
        FROM directories d
        LEFT JOIN user_directories ud ON ud.directory_id = d.id AND ud.user_id = ?
        LEFT JOIN files f ON f.directory_id = d.id
        LEFT JOIN user_files uf ON uf.file_id = f.id AND uf.user_id = ?
    `

	var rawResults []struct {
		DirectoryID       uint    `json:"directory_id"`
		NameFolder        string  `json:"name_folder"`
		ParentPathID      *uint   `json:"parent_path_id"`
		UserHasAccess     bool    `json:"user_has_access"`
		FileID            *uint   `json:"file_id"`
		NameFile          *string `json:"name_file"`
		FileDirectoryID   *uint   `json:"file_directory_id"`
		FileUserHasAccess *bool   `json:"file_user_has_access"`
	}

	if err := r.db.WithContext(ctx).Raw(query, userID, userID).Scan(&rawResults).Error; err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	directoryMap := make(map[uint]*domain.DirectoryUserResponse)
	for _, row := range rawResults {
		if _, exists := directoryMap[row.DirectoryID]; !exists {
			directoryMap[row.DirectoryID] = &domain.DirectoryUserResponse{
				ID:            row.DirectoryID,
				NameFolder:    row.NameFolder,
				ParentPathID:  row.ParentPathID,
				UserHasAccess: row.UserHasAccess,
				Files:         []domain.FileUserResponse{},
			}
		}

		if row.FileID != nil {
			directoryMap[row.DirectoryID].Files = append(directoryMap[row.DirectoryID].Files, domain.FileUserResponse{
				ID:            *row.FileID,
				NameFile:      *row.NameFile,
				DirectoryID:   *row.FileDirectoryID,
				UserHasAccess: *row.FileUserHasAccess,
			})
		}
	}

	var result []domain.DirectoryUserResponse
	for _, dir := range directoryMap {
		result = append(result, *dir)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})

	return result, nil
}

func (r *DirectoryRepository) GetWorkflowFileTree(ctx context.Context, workflowID uint) ([]domain.DirectoryWorkflowResponse, error) {
	const op = "infrastructure.postgresrepo.directory.GetWorkflowFileTree"

	query := `
        SELECT 
            d.id AS directory_id,
            d.name AS name_folder,
            d.parent_path_id AS parent_path_id,
            (d.workflow_id = ?) AS current_workflow_assigned
        FROM directories d
    `

	var rawResults []struct {
		DirectoryID             uint   `json:"directory_id"`
		NameFolder              string `json:"name_folder"`
		ParentPathID            *uint  `json:"parent_path_id"`
		CurrentWorkflowAssigned bool   `json:"current_workflow_assigned"`
	}

	if err := r.db.WithContext(ctx).Raw(query, workflowID).Scan(&rawResults).Error; err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	directoryMap := make(map[uint]*domain.DirectoryWorkflowResponse)
	for _, row := range rawResults {
		if _, exists := directoryMap[row.DirectoryID]; !exists {
			directoryMap[row.DirectoryID] = &domain.DirectoryWorkflowResponse{
				ID:                      row.DirectoryID,
				NameFolder:              row.NameFolder,
				ParentPathID:            row.ParentPathID,
				CurrentWorkflowAssigned: row.CurrentWorkflowAssigned,
			}
		}
	}

	var result []domain.DirectoryWorkflowResponse
	for _, dir := range directoryMap {
		result = append(result, *dir)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})

	return result, nil
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

func (directoryRepository *DirectoryRepository) CheckDirectoriesExist(ctx context.Context, directoryIDs []uint) (bool, error) {
	const op = "infrastructure.postgresrepo.directory.CheckDirectoriesExist"

	if len(directoryIDs) == 0 {
		return false, fmt.Errorf("%s: empty user list", op)
	}

	var count int64
	err := directoryRepository.db.WithContext(ctx).
		Model(&domain.Directory{}).
		Where("id IN (?)", directoryIDs).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return count == int64(len(directoryIDs)), nil
}

func (r *DirectoryRepository) UpdateDirectories(ctx context.Context, workflowID uint, directoryIDs []uint) error {
	const op = "infrastructure.postgresrepo.directory.UpdateDirectories"

	// Начало транзакции
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Обновление workflow_id для всех указанных directory_ids
	result := tx.Model(&domain.Directory{}).
		Where("id IN ?", directoryIDs).
		Update("workflow_id", workflowID)

	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, result.Error)
	}

	// Проверка, были ли затронуты строки
	if result.RowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, domain.ErrDirectoryNotFound)
	}

	// Фиксация изменений
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *DirectoryRepository) UpdateUserDirectoryRelations(ctx context.Context, userID uint, directoryIDs []uint) error {
	const op = "infrastructure.postgresrepo.directory.UpdateUserDirectoryRelations"

	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Where("user_id = ?", userID).Delete(&domain.UserDirectory{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: failed to delete user directory relations: %w", op, err)
	}

	var newRelations []domain.UserDirectory
	for _, directoryID := range directoryIDs {
		newRelations = append(newRelations, domain.UserDirectory{
			UserID:      userID,
			DirectoryID: directoryID,
		})
	}

	if len(newRelations) > 0 {
		if err := tx.Create(&newRelations).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("%s: failed to create user directory relations: %w", op, err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
