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

type DirectoryUsecase struct {
	directoryRepo interfaces.DirectoryRepository
	log           *slog.Logger
}

func NewDirectoryUsecase(directoryRepo interfaces.DirectoryRepository, log *slog.Logger) *DirectoryUsecase {
	return &DirectoryUsecase{
		directoryRepo: directoryRepo,
		log:           log,
	}
}

func (u *DirectoryUsecase) GetFileTree(ctx context.Context, isArchive bool, userID uint) (response domain.GetFileTreeResponse, err error) {
	const op = "usecases.directory.GetFileTree"

	log := u.log.With(slog.String("op", op))
	log.Info("getting file tree")

	var directories []domain.Directory

	switch isArchive {
	case false:
		directories, err = u.directoryRepo.GetFileTreeWorking(ctx, userID)
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

	case true:
		directories, err = u.directoryRepo.GetFileTreeArchive(ctx, userID)
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
	}

	log.Debug("directories retrieved successfully")

	response.Data = make([]domain.DirectoryResponse, 0, len(directories))

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

func (u *DirectoryUsecase) CreateDirectory(ctx context.Context, parentPathID *uint, name string, userID uint) (err error) {
	const op = "usecase.directory.CreateDirectory"
	log := u.log.With(slog.String("op", op), slog.Any("user_id", userID))

	log.Info("uploading directory")

	log.Debug("checking acces to directory", slog.Any("parent_path_id", parentPathID))
	if parentPathID != nil {
		hasAccess, err := u.directoryRepo.CheckUserDirectoryAccess(ctx, userID, *parentPathID)
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
			log.Warn("user has no access to directory", slog.Any("directory_id", parentPathID))
			return fmt.Errorf("%s: %w", op, domain.ErrAccessDenied)
		}
	}

	log.Debug("creating directory")
	status := "draft"
	err = u.directoryRepo.CreateDirectory(ctx, parentPathID, name, status, userID)
	if err != nil {
		if errors.Is(err, domain.ErrDirectoryAlreadyExists) {
			log.Error("directory already exists", slogger.Err(domain.ErrDirectoryAlreadyExists))
			return fmt.Errorf("%s: %w", op, domain.ErrDirectoryAlreadyExists)
		}
		log.Error("failed to create directory", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("directory created successfully")
	return nil
}

func (u *DirectoryUsecase) DeleteDirectory(ctx context.Context, directoryID uint, userID uint) (err error) {
	const op = "usecase.directory.DeleteDirectory"

	log := u.log.With(slog.String("op", op), slog.Any("user_id", userID), slog.Any("directory_id", directoryID))
	log.Info("deleting directory")

	err = u.directoryRepo.DeleteDirectory(ctx, directoryID, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrDirectoryNotFound):
			log.Error("directory not found", slogger.Err(domain.ErrDirectoryNotFound))
			return fmt.Errorf("%s: %w", op, domain.ErrDirectoryNotFound)
		case errors.Is(err, domain.ErrAccessDenied):
			log.Error("user has no access to this directory", slogger.Err(domain.ErrAccessDenied))
			return fmt.Errorf("%s: %w", op, domain.ErrAccessDenied)
		case errors.Is(err, domain.ErrDirectoryContainsNonDraftFiles):
			log.Error("directory contains non draft files", slogger.Err(domain.ErrDirectoryContainsNonDraftFiles))
			return fmt.Errorf("%s: %w", op, domain.ErrDirectoryContainsNonDraftFiles)
		}
		log.Error("failed to delete directory", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("directory deleted successfully")
	return nil
}
