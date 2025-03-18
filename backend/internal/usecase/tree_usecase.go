package usecase

import (
	"backend/internal/domain"
	"backend/internal/domain/interfaces"
	"context"
)

type FileTreeUsecase struct {
	fileTreeRepo interfaces.FileTreeRepository
}

func NewFileTreeUsecase(fileTreeRepo interfaces.FileTreeRepository) interfaces.FileTreeUsecase {
	return &FileTreeUsecase{fileTreeRepo: fileTreeRepo}
}

func (u *FileTreeUsecase) GetFileTree(ctx context.Context, isArchive bool, userID uint) (domain.GetFileTreeResponse, error) {
	directories, err := u.fileTreeRepo.GetDirectoriesWithFiles(ctx, isArchive, userID)
	if err != nil {
		return domain.GetFileTreeResponse{}, err
	}

	// Преобразование данных в формат ответа
	response := domain.GetFileTreeResponse{
		Data: make([]domain.DirectoryResponse, 0),
	}

	for _, dir := range directories {
		dirData := domain.DirectoryResponse{
			NameFolder: dir.Name,
			Status:     dir.Status,
		}
		parentPathID := dir.ParentPathID
		if parentPathID != nil {
			dirData.ParentPathID = parentPathID
		}
		dirData.Files = make([]domain.FileResponse, 0)

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

	return response, nil
}

func (u *FileTreeUsecase) UploadDirectory(ctx context.Context, parentPathID *uint, name string, userID uint) (uint, error) {
	status := "draft"
	return u.fileTreeRepo.CreateDirectory(ctx, parentPathID, name, status)
}

func (u *FileTreeUsecase) UploadFile(ctx context.Context, directoryID uint, name string, userID uint) (uint, error) {
	status := "draft"
	return u.fileTreeRepo.CreateFile(ctx, directoryID, name, status)
}
