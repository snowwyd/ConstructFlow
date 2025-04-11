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

type AdminUsecase struct {
	fileMetadataRepo interfaces.FileMetadataRepository
	directoryRepo    interfaces.DirectoryRepository
	log              *slog.Logger
}

func NewAdminUsecase(directoryRepo interfaces.DirectoryRepository, fileMetadataRepo interfaces.FileMetadataRepository, log *slog.Logger) *AdminUsecase {
	return &AdminUsecase{
		directoryRepo:    directoryRepo,
		fileMetadataRepo: fileMetadataRepo,
		log:              log,
	}
}

func (adminUsecase *AdminUsecase) GetUserTree(ctx context.Context, userID, actorID uint) (directories []domain.DirectoryUserResponse, err error) {
	const op = "usecases.admin.GetUserTree"

	log := adminUsecase.log.With(slog.String("op", op))
	log.Info("getting file tree for user")

	// TODO: admin check from microservice (low priority)

	directories, err = adminUsecase.directoryRepo.GetUserFileTree(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrAccessDenied):
			log.Error("user has no access to this repository", slogger.Err(domain.ErrAccessDenied))
			return nil, fmt.Errorf("%s: %w", op, domain.ErrAccessDenied)
		default:
			log.Error("failed to get file tree", slogger.Err(err))
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	log.Info("file tree response prepared successfully")
	return directories, nil
}

func (adminUsecase *AdminUsecase) GetWorkflowTree(ctx context.Context, workflowID, actorID uint) (directories []domain.DirectoryWorkflowResponse, err error) {
	const op = "usecases.admin.GetWorkflowTree"

	log := adminUsecase.log.With(slog.String("op", op))
	log.Info("getting file tree for workflow")

	// TODO: admin check from microservice (low priority)

	directories, err = adminUsecase.directoryRepo.GetWorkflowFileTree(ctx, workflowID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrAccessDenied):
			log.Error("user has no access to this repository", slogger.Err(domain.ErrAccessDenied))
			return nil, fmt.Errorf("%s: %w", op, domain.ErrAccessDenied)
		default:
			log.Error("failed to get file tree", slogger.Err(err))
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	log.Info("file tree response prepared successfully")
	return directories, nil
}
