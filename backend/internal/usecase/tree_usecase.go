package usecase

import (
	"backend/internal/domain"
	"backend/internal/domain/interfaces"
	"backend/pkg/logger/slogger"
	"context"
	"errors"
	"fmt"
	"log/slog"
)

type FileTreeUsecase struct {
	fileTreeRepo interfaces.FileTreeRepository
	log          *slog.Logger
}

func NewFileTreeUsecase(fileTreeRepo interfaces.FileTreeRepository, log *slog.Logger) interfaces.FileTreeUsecase {
	return &FileTreeUsecase{
		fileTreeRepo: fileTreeRepo,
		log:          log,
	}
}

func (u *FileTreeUsecase) GetFileTree(ctx context.Context, isArchive bool, userID uint) (domain.GetFileTreeResponse, error) {
	const op = "usecase.tree.GetFileTree"
	log := u.log.With(slog.String("op", op), slog.Any("user_id", userID), slog.Bool("archive", isArchive))

	log.Info("starting file tree retrieval")

	// Получение директорий с файлами из репозитория
	directories, err := u.fileTreeRepo.GetDirectoriesWithFiles(ctx, isArchive, userID)
	if err != nil {
		log.Error("failed to get directories from repository", slogger.Err(err))
		return domain.GetFileTreeResponse{}, fmt.Errorf("%s: %w", op, domain.ErrInternal)
	}

	log.Debug("directories retrieved successfully", slog.Int("count", len(directories)))

	// Формирование ответа
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

func (u *FileTreeUsecase) GetFileInfo(ctx context.Context, fileID uint, userID uint) (*domain.FileResponse, error) {
	const op = "usecase.file_tree.GetFileInfo"
	log := u.log.With(slog.String("op", op), slog.Any("user_id", userID), slog.Any("file_id", fileID))

	log.Info("getting file info")

	// Проверка доступа пользователя к файлу
	hasAccess, err := u.fileTreeRepo.CheckUserFileAccess(ctx, userID, fileID)
	if err != nil {
		if errors.Is(err, domain.ErrFileNotFound) {
			log.Error("failed to get file", slogger.Err(domain.ErrFileNotFound))
		}
		log.Error("failed to check user file access", slogger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if !hasAccess {
		log.Warn("user has no access to the file")
		return nil, domain.ErrAccessDenied
	}

	// Получение информации о файле
	file, err := u.fileTreeRepo.GetFileInfo(ctx, fileID)
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

func (u *FileTreeUsecase) UploadDirectory(ctx context.Context, parentPathID *uint, name string, userID uint) (uint, error) {
	const op = "usecase.file.UploadDirectory"
	log := u.log.With(slog.String("op", op), slog.Any("user_id", userID))

	log.Info("uploading directory")
	// Проверка доступа к родительской директории
	if parentPathID != nil {
		hasAccess, err := u.fileTreeRepo.CheckUserDirectoryAccess(ctx, userID, *parentPathID)
		if err != nil {
			log.Error("failed to check directory access", slogger.Err(err))
			return 0, fmt.Errorf("%s: %w", op, domain.ErrInternal)
		}
		if !hasAccess {
			log.Warn("user has no access to parent directory", slog.Any("parent_id", *parentPathID))
			return 0, domain.ErrAccessDenied
		}
	}

	status := "wip"
	dirID, err := u.fileTreeRepo.CreateDirectory(ctx, parentPathID, name, status, userID)
	if err != nil {
		log.Error("failed to create directory", slogger.Err(err))
		return 0, fmt.Errorf("%s: %w", op, domain.ErrInternal)
	}

	log.Info("directory created successfully", slog.Any("directory_id", dirID))
	return dirID, nil
}

func (u *FileTreeUsecase) UploadFile(ctx context.Context, directoryID uint, name string, userID uint) (uint, error) {
	const op = "usecase.file.UploadFile"
	log := u.log.With(slog.String("op", op), slog.Any("user_id", userID))

	log.Info("uploading file")
	// Проверка доступа к директории
	hasAccess, err := u.fileTreeRepo.CheckUserDirectoryAccess(ctx, userID, directoryID)
	if err != nil {
		log.Error("failed to check directory access", slogger.Err(err))
		return 0, fmt.Errorf("%s: %w", op, domain.ErrInternal)
	}
	if !hasAccess {
		log.Warn("user has no access to directory", slog.Any("directory_id", directoryID))
		return 0, domain.ErrAccessDenied
	}

	status := "wip"
	fileID, err := u.fileTreeRepo.CreateFile(ctx, directoryID, name, status, userID)
	if err != nil {
		log.Error("failed to create file", slogger.Err(err))
		return 0, fmt.Errorf("%s: %w", op, domain.ErrInternal)
	}

	log.Info("file created successfully", slog.Any("file_id", fileID))
	return fileID, nil
}

func (u *FileTreeUsecase) DeleteDirectory(ctx context.Context, directoryID uint, userID uint) error {
	const op = "usecase.file_tree.DeleteDirectory"
	log := u.log.With(slog.String("op", op), slog.Any("user_id", userID), slog.Any("directory_id", directoryID))

	log.Info("deleting directory")
	err := u.fileTreeRepo.DeleteDirectory(ctx, directoryID, userID)
	if err != nil {
		log.Error("failed to delete directory", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("directory deleted successfully")
	return nil
}

func (u *FileTreeUsecase) DeleteFile(ctx context.Context, fileID uint, userID uint) error {
	const op = "usecase.file_tree.DeleteFile"
	log := u.log.With(slog.String("op", op), slog.Any("user_id", userID), slog.Any("file_id", fileID))

	log.Info("deleting file")
	err := u.fileTreeRepo.DeleteFile(ctx, fileID, userID)
	if err != nil {
		log.Error("failed to delete file", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("file deleted successfully")
	return nil
}
