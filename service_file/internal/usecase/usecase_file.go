package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"service-file/internal/domain"
	"service-file/internal/domain/interfaces"
	"service-file/pkg/logger/slogger"
	"time"
)

type FileUsecase struct {
	fileMetadataRepo interfaces.FileMetadataRepository
	directoryRepo    interfaces.DirectoryRepository
	fileStorage      interfaces.FileStorage
	log              *slog.Logger
}

func NewFileUsecase(directoryRepo interfaces.DirectoryRepository, fileMetadataRepo interfaces.FileMetadataRepository, fileStorage interfaces.FileStorage, log *slog.Logger) *FileUsecase {
	return &FileUsecase{
		directoryRepo:    directoryRepo,
		fileMetadataRepo: fileMetadataRepo,
		fileStorage:      fileStorage,
		log:              log,
	}
}

func (u *FileUsecase) GetFileInfo(ctx context.Context, fileID uint, userID uint) (*domain.FileResponse, error) {
	const op = "usecase.file.GetFileInfo"

	log := u.log.With(slog.String("op", op), slog.Any("user_id", userID), slog.Any("file_id", fileID))
	log.Info("getting file info")

	file, err := u.fileMetadataRepo.GetFileInfo(ctx, fileID, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrFileNotFound):
			log.Error("file not found", slogger.Err(err))
			return nil, fmt.Errorf("%s: %w", op, domain.ErrFileNotFound)
		case errors.Is(err, domain.ErrAccessDenied):
			log.Warn("access denied", slogger.Err(err))
			return nil, fmt.Errorf("%s: %w", op, domain.ErrAccessDenied)
		default:
			log.Error("failed to get file info", slogger.Err(err))
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	response := &domain.FileResponse{
		ID:          file.ID,
		NameFile:    file.Name,
		Status:      file.Status,
		DirectoryID: file.DirectoryID,
	}

	log.Info("file info response prepared successfully")
	return response, nil
}

func (u *FileUsecase) CreateFile(ctx context.Context, directoryID uint, name string, fileData []byte, userID uint) (err error) {
	const op = "usecase.file.CreateFile"

	log := u.log.With(slog.String("op", op), slog.Any("user_id", userID))
	log.Info("uploading file")

	log.Debug("checking acces to directory", slog.Any("directory_id", directoryID))
	hasAccess, err := u.directoryRepo.CheckUserDirectoryAccess(ctx, userID, directoryID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrDirectoryNotFound):
			log.Error("directory not found", slogger.Err(domain.ErrDirectoryNotFound))
			return fmt.Errorf("%s: %w", op, domain.ErrDirectoryNotFound)
		default:
			log.Error("failed to check directory access", slogger.Err(err))
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	if !hasAccess {
		log.Warn("user has no access to directory", slog.Any("directory_id", directoryID))
		return fmt.Errorf("%s: %w", op, domain.ErrAccessDenied)
	}

	log.Debug("uploading file into minio")
	objectKey := fmt.Sprintf("files/%d/%s", directoryID, name)
	if err := u.fileStorage.UploadFile(ctx, "files", objectKey, fileData); err != nil {
		log.Error("failed to upload file into minio")
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("creating file")
	status := "draft"
	err = u.fileMetadataRepo.CreateFile(ctx, directoryID, name, status, userID, objectKey)
	if err != nil {
		if errors.Is(err, domain.ErrFileAlreadyExists) {
			log.Error("file already exists", slogger.Err(domain.ErrFileAlreadyExists))
			return fmt.Errorf("%s: %w", op, domain.ErrFileAlreadyExists)
		}
		log.Error("failed to create file", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("file created successfully")
	return nil
}

func (u *FileUsecase) UpdateFile(
	ctx context.Context,
	fileID uint,
	newData []byte,
	userID uint,
) error {
	const op = "usecase.file.UpdateFile"
	log := u.log.With(slog.String("op", op), slog.Any("file_id", fileID))

	// 1. Получаем текущий файл
	file, err := u.fileMetadataRepo.GetFileByID(ctx, fileID)
	if err != nil {
		if errors.Is(err, domain.ErrFileNotFound) {
			log.Error("file not found", slogger.Err(domain.ErrFileNotFound))
			return fmt.Errorf("%s: %w", op, domain.ErrFileNotFound)
		}
		log.Error("failed to get file", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	// 2. Проверяем доступ пользователя к директории
	hasAccess, err := u.directoryRepo.CheckUserDirectoryAccess(ctx, userID, file.DirectoryID)
	if err != nil {
		log.Error("failed to check directory access", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	if !hasAccess {
		log.Warn("user has no access to directory", slogger.Err(domain.ErrAccessDenied))
		return fmt.Errorf("%s: %w", op, domain.ErrAccessDenied)
	}

	// 4. Загружаем новую версию в MinIO
	newKey, err := u.fileStorage.UploadNewVersion(ctx, "files", file.MinioObjectKey, newData)
	if err != nil {
		log.Error("failed to upload new version", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	// 5. Обновляем метаданные файла
	file.MinioObjectKey = newKey
	file.Version++
	file.UpdatedAt = time.Now()

	// 6. Сохраняем изменения в БД
	if err := u.fileMetadataRepo.UpdateFile(ctx, file); err != nil {
		log.Error("failed to update file metadata", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("file updated successfully", slog.Int("new_version", file.Version))
	return nil
}

func (u *FileUsecase) DeleteFile(ctx context.Context, fileID uint, userID uint) (err error) {
	const op = "usecase.file.DeleteFile"
	log := u.log.With(slog.String("op", op), slog.Any("user_id", userID), slog.Any("file_id", fileID))

	log.Info("deleting file")
	err = u.fileMetadataRepo.DeleteFile(ctx, fileID, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrFileNotFound):
			log.Error("file not found", slogger.Err(domain.ErrFileNotFound))
			return fmt.Errorf("%s: %w", op, domain.ErrFileNotFound)
		case errors.Is(err, domain.ErrAccessDenied):
			log.Error("user has no access to this directory", slogger.Err(domain.ErrAccessDenied))
			return fmt.Errorf("%s: %w", op, domain.ErrAccessDenied)
		case errors.Is(err, domain.ErrCannotDeleteNonDraftFile):
			log.Error("cannot delete non draft file", slogger.Err(domain.ErrCannotDeleteNonDraftFile))
			return fmt.Errorf("%s: %w", op, domain.ErrCannotDeleteNonDraftFile)
		}
		log.Error("failed to delete file", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("file deleted successfully")
	return nil
}
