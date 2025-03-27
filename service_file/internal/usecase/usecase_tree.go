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

func (u *FileTreeUsecase) GetFileByID(ctx context.Context, fileID uint) (domain.FileResponse, error) {
	const op = "usecases.tree.GetFileByID"

	log := u.log.With(slog.String("op", op))
	log.Info("getting file by id")

	file, err := u.fileMetadataRepo.GetFileByID(ctx, fileID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrFileNotFound):
			log.Error("file not found", slogger.Err(domain.ErrFileNotFound))
			return domain.FileResponse{}, fmt.Errorf("%s: %w", op, domain.ErrFileNotFound)
		default:
			log.Error("failed to get file info", slogger.Err(err))
			return domain.FileResponse{}, fmt.Errorf("%s: %w", op, err)
		}
	}

	response := domain.FileResponse{
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
	panic("implement me!")
}

func (u *FileTreeUsecase) CreateDirectory(ctx context.Context, directoryID *uint, name string, userID uint) (err error) {
	panic("implement me!")
}

func (u *FileTreeUsecase) DeleteFile(ctx context.Context, fileID uint, userID uint) (err error) {
	panic("implement me!")
}

func (u *FileTreeUsecase) DeleteDirectory(ctx context.Context, directoryID uint, userID uint) (err error) {
	panic("implement me!")
}

func (u *FileTreeUsecase) UpdateFileStatus(ctx context.Context, fileID uint, status string) error {
	panic("implement me!")
}

func (u *FileTreeUsecase) CheckAccessToFile(ctx context.Context, fileID, userID uint) (bool, error) {
	panic("implement me!")
}

func (u *FileTreeUsecase) CheckAccessToDirectory(ctx context.Context, directoryID, userID uint) (bool, error) {
	panic("implement me!")
}
