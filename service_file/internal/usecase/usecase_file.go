package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"service-file/internal/domain"
	"service-file/internal/domain/interfaces"
	"service-file/pkg/logger/slogger"
)

type FileUsecase struct {
	fileMetadataRepo interfaces.FileMetadataRepository
	directoryRepo    interfaces.DirectoryRepository
	log              *slog.Logger
}

func NewFileUsecase(directoryRepo interfaces.DirectoryRepository, fileMetadataRepo interfaces.FileMetadataRepository, log *slog.Logger) *FileUsecase {
	return &FileUsecase{
		directoryRepo:    directoryRepo,
		fileMetadataRepo: fileMetadataRepo,
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

func (u *FileUsecase) CreateFile(ctx context.Context, directoryID uint, name string, userID uint) (err error) {
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

	log.Debug("creating file")
	status := "draft"
	err = u.fileMetadataRepo.CreateFile(ctx, directoryID, name, status, userID)
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
