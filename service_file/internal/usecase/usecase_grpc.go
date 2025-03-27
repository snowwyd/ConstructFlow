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
	log              *slog.Logger
}

func NewGRPCUsecase(fileMetadataRepo interfaces.FileMetadataRepository, log *slog.Logger) *GRPCUsecase {
	return &GRPCUsecase{
		fileMetadataRepo: fileMetadataRepo,
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
