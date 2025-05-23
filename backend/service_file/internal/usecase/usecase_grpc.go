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

type GRPCUsecase struct {
	fileMetadataRepo interfaces.FileMetadataRepository
	directoryRepo    interfaces.DirectoryRepository
	log              *slog.Logger
}

func NewGRPCUsecase(fileMetadataRepo interfaces.FileMetadataRepository, directoryRepo interfaces.DirectoryRepository, log *slog.Logger) *GRPCUsecase {
	return &GRPCUsecase{
		fileMetadataRepo: fileMetadataRepo,
		directoryRepo:    directoryRepo,
		log:              log,
	}
}

func (u *GRPCUsecase) GetFileByID(ctx context.Context, fileID uint) (*domain.File, error) {
	const op = "usecases.grpc.GetFileByID"

	log := u.log.With(slog.String("op", op))
	log.Info("getting file by id")

	file, err := u.fileMetadataRepo.GetFileByID(ctx, fileID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrFileNotFound):
			log.Error("file not found", slogger.Err(domain.ErrFileNotFound))
			return nil, fmt.Errorf("%s: %w", op, domain.ErrFileNotFound)
		default:
			log.Error("failed to get file info", slogger.Err(err))
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	log.Info("file by id got successfully")
	return file, nil
}

func (u *GRPCUsecase) UpdateFileStatus(ctx context.Context, fileID uint, status string) error {
	const op = "usecases.grpc.UpdateFileStatus"

	log := u.log.With(slog.String("op", op))
	log.Info("updating file status")

	tx := u.fileMetadataRepo.GetDB().Begin()
	defer tx.Rollback()

	if err := u.fileMetadataRepo.UpdateFileStatus(ctx, fileID, status, tx); err != nil {
		switch {
		case errors.Is(err, domain.ErrFileNotFound):
			log.Error("file not found", slogger.Err(domain.ErrFileNotFound))
			return fmt.Errorf("%s: %w", op, domain.ErrFileNotFound)
		default:
			log.Error("failed to update file status", slogger.Err(err))
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Error("failed to commit transaction", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("file status updated successfully")
	return nil
}

func (u *GRPCUsecase) GetFilesByID(ctx context.Context, fileIDs []uint32) ([]domain.File, error) {
	const op = "usecases.grpc.GetFilesByID"

	log := u.log.With(slog.String("op", op))
	log.Info("getting files by id")

	var files []domain.File
	if err := u.fileMetadataRepo.GetFilesByID(ctx, fileIDs, &files); err != nil {
		log.Error("failed to get files info", slogger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("files by id got successfully")
	return files, nil
}

func (grpcUsecase *GRPCUsecase) CheckWorkflow(ctx context.Context, workflowID uint) (bool, error) {
	const op = "usecases.grpc.CheckWorkflow"

	log := grpcUsecase.log.With(slog.String("op", op))
	log.Info("checking workflow existence")

	log.Debug("checking database")
	exists, err := grpcUsecase.directoryRepo.CheckWorkflow(ctx, workflowID)
	if err != nil {
		log.Error("failed to check workflow existence", slogger.Err(err))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("workflow existence checked successfully")
	return exists, nil
}

func (grpcUsecase *GRPCUsecase) AssignWorkflow(ctx context.Context, workflowID uint32, directoryIDs []uint32) error {
	const op = "usecases.grpc.AssignWorkflow"

	log := grpcUsecase.log.With(slog.String("op", op))
	log.Info("assigning workflow")

	log.Debug("checking directories")
	directories := make([]uint, len(directoryIDs))
	for i, value := range directoryIDs {
		directories[i] = uint(value)
	}

	exists, err := grpcUsecase.directoryRepo.CheckDirectoriesExist(ctx, directories)
	if err != nil {
		log.Error("failed to check workflow existence", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	if !exists {
		return domain.ErrDirectoryNotFound
	}

	log.Debug("updating directory data")
	if err := grpcUsecase.directoryRepo.UpdateDirectories(ctx, uint(workflowID), directories); err != nil {
		log.Error("failed to update directories", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("workflow assigned successfully")
	return nil
}

func (grpcUsecase *GRPCUsecase) DeleteUserRelations(ctx context.Context, userID uint) error {
	const op = "usecases.grpc.DeleteUserRelations"

	log := grpcUsecase.log.With(slog.String("op", op))
	log.Info("deleting user relations")

	log.Debug("deleting user_directories relations")
	if err := grpcUsecase.directoryRepo.DeleteUserRelations(ctx, userID); err != nil {
		if errors.Is(err, domain.ErrNoRelationsFound) {
			log.Warn("no user relations found")
			return nil
		}
		log.Error("failed to delete user_directories relations", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("deleting user_files relations")
	if err := grpcUsecase.fileMetadataRepo.DeleteUserRelations(ctx, userID); err != nil {
		if errors.Is(err, domain.ErrNoRelationsFound) {
			log.Warn("no user relations found")
			return nil
		}
		log.Error("failed to delete user_files relations", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user relations deleted successfully")
	return nil
}

func (grpcUsecase *GRPCUsecase) AssignUser(ctx context.Context, userID uint, directoryIDs []uint32, fileIDs []uint32) error {
	const op = "usecases.grpc.AssignUser"

	log := grpcUsecase.log.With(slog.String("op", op))
	log.Info("assigning user")

	log.Debug("preparing directories")
	directories := make([]uint, len(directoryIDs))
	for i, value := range directoryIDs {
		directories[i] = uint(value)
	}

	log.Debug("preparing files")
	files := make([]uint, len(fileIDs))
	for i, value := range fileIDs {
		files[i] = uint(value)
	}

	exists, err := grpcUsecase.directoryRepo.CheckDirectoriesExist(ctx, directories)
	if err != nil {
		log.Error("failed to check directories existence", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	if !exists {
		return domain.ErrDirectoryNotFound
	}

	exists, err = grpcUsecase.fileMetadataRepo.CheckFilesExist(ctx, files)
	if err != nil {
		log.Error("failed to check files existence", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	if !exists {
		return domain.ErrFileNotFound
	}

	log.Debug("updating directory data")
	if err := grpcUsecase.directoryRepo.UpdateUserDirectoryRelations(ctx, userID, directories); err != nil {
		log.Error("failed to update directories", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("updating file data")
	if err := grpcUsecase.fileMetadataRepo.UpdateUserFileRelations(ctx, userID, files); err != nil {
		log.Error("failed to update files", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user assigned successfully")
	return nil
}
