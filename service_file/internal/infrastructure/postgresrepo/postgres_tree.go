package postgresrepo

import (
	"context"
	"service-file/internal/domain"
	"service-file/internal/domain/interfaces"

	"gorm.io/gorm"
)

type FileTreeRepository struct {
	db *gorm.DB
}

func NewFileTreeRepository(db *Database) *FileTreeRepository {
	return &FileTreeRepository{db: db.db}
}

func (r *FileTreeRepository) WithTx(tx *gorm.DB) interfaces.FileTreeRepository {
	return &FileTreeRepository{db: tx}
}

func (r *FileTreeRepository) GetDB() *gorm.DB {
	return r.db
}

// GetDirectoriesWithFiles отдает список всех директорий и файлов, соответствующих входным условиям
func (r *FileTreeRepository) GetDirectoriesWithFiles(ctx context.Context, isArchive bool, userID uint) ([]*domain.Directory, error) {
	panic("implement me!")
}

func (r *FileTreeRepository) GetFileInfo(ctx context.Context, fileID uint) (*domain.File, error) {
	panic("implement me!")
}

// CreateDirectory записывает новую сущность Directory в БД
func (r *FileTreeRepository) CreateDirectory(ctx context.Context, parentPathID *uint, name string, status string, userID uint) error {
	panic("implement me!")
}

// CreateFile записывает новую сущность File в БД
func (r *FileTreeRepository) CreateFile(ctx context.Context, directoryID uint, name string, status string, userID uint) error {
	panic("implement me!")
}

// DeleteDirectory рекурсивно удаляет все файлы и директории внутри указанной
func (r *FileTreeRepository) DeleteDirectory(ctx context.Context, directoryID uint, userID uint) error {
	panic("implement me!")
}

// DeleteFile удаляет файл из БД
func (r *FileTreeRepository) DeleteFile(ctx context.Context, fileID uint, userID uint) error {
	panic("implement me!")
}

// CheckUserDirectoryAccess проверяет доступ пользователя к директории
func (r *FileTreeRepository) CheckUserDirectoryAccess(ctx context.Context, userID, directoryID uint) (bool, error) {
	panic("implement me!")
}

// CheckUserFileAccess проверяет доступ пользователя к файлу
func (r *FileTreeRepository) CheckUserFileAccess(ctx context.Context, userID, fileID uint) (bool, error) {
	panic("implement me!")
}

// GetFileWithDirectory получает файл с директорией
// Кастомные ошибки: ErrFileNotFound
func (r *FileTreeRepository) GetFileWithDirectory(ctx context.Context, fileID uint, tx *gorm.DB) (*domain.File, error) {
	panic("implement me!")
}

// UpdateFileStatus обновляет статус файла
func (r *FileTreeRepository) UpdateFileStatus(ctx context.Context, file *domain.File, tx *gorm.DB) error {
	panic("implement me!")
}
