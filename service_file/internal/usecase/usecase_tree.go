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

type FileTreeUsecase struct {
	fileMetadataRepo interfaces.FileMetadataRepository
	log              *slog.Logger
}

func NewFileTreeUsecase(fileMetadataRepo interfaces.FileMetadataRepository, log *slog.Logger) *FileTreeUsecase {
	return &FileTreeUsecase{
		fileMetadataRepo: fileMetadataRepo,
		log:              log,
	}
}

func (u *FileTreeUsecase) GetFileTree(ctx context.Context, isArchive bool, userID uint) (domain.GetFileTreeResponse, error) {
	const op = "usecases.tree.GetFileTree"

	log := u.log.With(slog.String("op", op))
	log.Info("getting file tree")

	directories, err := u.fileMetadataRepo.GetFileTree(ctx, isArchive, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrAccessDenied):
			log.Error("user has no acces to this repository", slogger.Err(domain.ErrAccessDenied))
			return domain.GetFileTreeResponse{}, fmt.Errorf("%s: %w", op, domain.ErrAccessDenied)
		default:
			log.Error("failed to get file tree", slogger.Err(err))
			return domain.GetFileTreeResponse{}, fmt.Errorf("%s: %w", op, err)
		}
	}

	log.Debug("directories retrieved successfully")

	response := domain.GetFileTreeResponse{
		Data: make([]domain.DirectoryResponse, 0, len(directories)),
	}

	for _, dir := range directories {
		dirData := domain.DirectoryResponse{
			ID:           dir.ID,
			NameFolder:   dir.Name,
			Status:       dir.Status,
			ParentPathID: dir.ParentPathID,
			Files:        make([]domain.FileResponse, 0, len(dir.Files)),
		}

		for _, file := range dir.Files {
			dirData.Files = append(dirData.Files, domain.FileResponse{
				ID:          file.ID,
				NameFile:    file.Name,
				Status:      file.Status,
				DirectoryID: file.DirectoryID,
			})
		}

		response.Data = append(response.Data, dirData)
	}

	log.Info("file tree response prepared successfully")
	return response, nil
}

func (u *FileTreeUsecase) GetFileByID(ctx context.Context, fileID uint) (domain.File, error) {
	const op = "usecases.tree.GetFileByID"

	log := u.log.With(slog.String("op", op))
	log.Info("getting file by id")

	file, err := u.fileMetadataRepo.GetFileByID(ctx, fileID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrFileNotFound):
			log.Error("file not found", slogger.Err(domain.ErrFileNotFound))
			return domain.File{}, fmt.Errorf("%s: %w", op, domain.ErrFileNotFound)
		default:
			log.Error("failed to get file info", slogger.Err(err))
			return domain.File{}, fmt.Errorf("%s: %w", op, err)
		}
	}

	log.Info("file by id got successfully")
	return *file, nil
}

func (u *FileTreeUsecase) GetFileInfo(ctx context.Context, fileID uint, userID uint) (*domain.FileResponse, error) {
	const op = "usecase.file_tree.GetFileInfo"
	log := u.log.With(slog.String("op", op), slog.Any("user_id", userID), slog.Any("file_id", fileID))

	log.Info("getting file info")

	// Проверка доступа пользователя к файлу
	hasAccess, err := u.fileMetadataRepo.CheckUserFileAccess(ctx, userID, fileID)
	if err != nil {
		if errors.Is(err, domain.ErrFileNotFound) {
			log.Error("failed to get file", slogger.Err(domain.ErrFileNotFound))
			return nil, domain.ErrFileNotFound
		}
		log.Error("failed to check user file access", slogger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if !hasAccess {
		log.Warn("user has no access to the file")
		return nil, domain.ErrAccessDenied
	}

	// Получение информации о файле
	file, err := u.fileMetadataRepo.GetFileInfo(ctx, fileID)
	if err != nil {
		if errors.Is(err, domain.ErrFileNotFound) {
			log.Warn("file not found in the database")
			return nil, domain.ErrFileNotFound
		}
		log.Error("failed to get file info", slogger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Формирование ответа
	response := &domain.FileResponse{
		ID:          file.ID,
		NameFile:    file.Name,
		Status:      file.Status,
		DirectoryID: file.DirectoryID,
	}

	log.Info("file info response prepared successfully")
	return response, nil
}

func (u *FileTreeUsecase) GetDirectoryByID(ctx context.Context, directoryID uint) (*domain.DirectoryResponse, error) {
	panic("implement me!")
}

func (u *FileTreeUsecase) CreateFile(ctx context.Context, directoryID uint, name string, userID uint) (err error) {
	const op = "usecase.file.UploadFile"
	log := u.log.With(slog.String("op", op), slog.Any("user_id", userID))

	log.Info("uploading file")
	// Проверка доступа к директории
	hasAccess, err := u.fileMetadataRepo.CheckUserDirectoryAccess(ctx, userID, directoryID)
	if err != nil {
		if errors.Is(err, domain.ErrDirectoryNotFound) {
			log.Error("directory not found", slogger.Err(domain.ErrDirectoryNotFound))
			return domain.ErrDirectoryNotFound
		}
		log.Error("failed to check directory access", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	if !hasAccess {
		log.Warn("user has no access to directory", slog.Any("directory_id", directoryID))
		return domain.ErrAccessDenied
	}

	status := "draft"
	err = u.fileMetadataRepo.CreateFile(ctx, directoryID, name, status, userID)
	if err != nil {
		log.Error("failed to create file", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("file created successfully")
	return nil
}

func (u *FileTreeUsecase) CreateDirectory(ctx context.Context, parentPathID *uint, name string, userID uint) (err error) {
	const op = "usecase.file.UploadDirectory"
	log := u.log.With(slog.String("op", op), slog.Any("user_id", userID))

	log.Info("uploading directory")
	// Проверка доступа к родительской директории
	if parentPathID != nil {
		hasAccess, err := u.fileMetadataRepo.CheckUserDirectoryAccess(ctx, userID, *parentPathID)
		if err != nil {
			if errors.Is(err, domain.ErrDirectoryNotFound) {
				log.Error("file not found", slogger.Err(domain.ErrDirectoryNotFound))
				return domain.ErrDirectoryNotFound
			}
			log.Error("failed to check directory access", slogger.Err(err))
			return fmt.Errorf("%s: %w", op, err)
		}
		if !hasAccess {
			log.Warn("user has no access to parent directory", slog.Any("parent_id", *parentPathID))
			return domain.ErrAccessDenied
		}
	}

	status := "draft"
	err = u.fileMetadataRepo.CreateDirectory(ctx, parentPathID, name, status, userID)
	if err != nil {
		log.Error("failed to create directory", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("directory created successfully")
	return nil
}

func (u *FileTreeUsecase) DeleteFile(ctx context.Context, fileID uint, userID uint) (err error) {
	const op = "usecase.file_tree.DeleteFile"
	log := u.log.With(slog.String("op", op), slog.Any("user_id", userID), slog.Any("file_id", fileID))

	log.Info("deleting file")
	err = u.fileMetadataRepo.DeleteFile(ctx, fileID, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrDirectoryNotFound):
			log.Error("directory not found", slogger.Err(domain.ErrDirectoryNotFound))
			return domain.ErrDirectoryNotFound
		case errors.Is(err, domain.ErrAccessDenied):
			log.Error("user has no access to this directory", slogger.Err(domain.ErrAccessDenied))
			return domain.ErrAccessDenied
		case errors.Is(err, domain.ErrCannotDeleteNonDraftFile):
			log.Error("cannot delete non draft file", slogger.Err(domain.ErrCannotDeleteNonDraftFile))
			return domain.ErrCannotDeleteNonDraftFile
		}
		log.Error("failed to delete file", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("file deleted successfully")
	return nil
}

func (u *FileTreeUsecase) DeleteDirectory(ctx context.Context, directoryID uint, userID uint) (err error) {
	const op = "usecase.file_tree.DeleteDirectory"
	log := u.log.With(slog.String("op", op), slog.Any("user_id", userID), slog.Any("directory_id", directoryID))

	log.Info("deleting directory")
	err = u.fileMetadataRepo.DeleteDirectory(ctx, directoryID, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrDirectoryNotFound):
			log.Error("directory not found", slogger.Err(domain.ErrDirectoryNotFound))
			return domain.ErrDirectoryNotFound
		case errors.Is(err, domain.ErrAccessDenied):
			log.Error("user has no access to this directory", slogger.Err(domain.ErrAccessDenied))
			return domain.ErrAccessDenied
		case errors.Is(err, domain.ErrDirectoryContainsNonDraftFiles):
			log.Error("directory contains non draft files", slogger.Err(domain.ErrDirectoryContainsNonDraftFiles))
			return domain.ErrDirectoryContainsNonDraftFiles
		}
		log.Error("failed to delete directory", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("directory deleted successfully")
	return nil
}

func (u *FileTreeUsecase) UpdateFileStatus(ctx context.Context, fileID uint, status string) error {
	const op = "usecases.tree.UpdateFileStatus"

	log := u.log.With(slog.String("op", op))
	log.Info("updating file status")

	tx := u.fileMetadataRepo.GetDB().Begin()
	defer tx.Rollback()

	if err := u.fileMetadataRepo.UpdateFileStatus(ctx, fileID, status, tx); err != nil {
		switch {
		case errors.Is(err, domain.ErrFileNotFound):
			log.Error("file not found", slogger.Err(domain.ErrFileNotFound))
			return domain.ErrFileNotFound
		default:
			log.Error("failed to update file status", slogger.Err(err))
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Error("failed to commit transaction", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("file status updated successfully")
	return nil
}

func (u *FileTreeUsecase) CheckAccessToFile(ctx context.Context, fileID, userID uint) (bool, error) {
	panic("implement me!")
}

func (u *FileTreeUsecase) CheckAccessToDirectory(ctx context.Context, directoryID, userID uint) (bool, error) {
	panic("implement me!")
}
