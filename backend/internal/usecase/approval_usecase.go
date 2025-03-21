package usecase

import (
	"backend/internal/domain"
	"backend/internal/domain/interfaces"
	"context"
	"log/slog"
)

type ApprovalUsecase struct {
	approvalRepo interfaces.ApprovalRepository
	fileTreeRepo interfaces.FileTreeRepository
	log          *slog.Logger
}

func NewApprovalUsecase(fileTreeRepo interfaces.FileTreeRepository, approvalRepo interfaces.ApprovalRepository, log *slog.Logger) *ApprovalUsecase {
	return &ApprovalUsecase{
		approvalRepo: approvalRepo,
		fileTreeRepo: fileTreeRepo,
		log:          log,
	}
}

func (u *ApprovalUsecase) ApproveFile(ctx context.Context, fileID uint) error {
	tx := u.fileTreeRepo.GetDB().Begin()
	defer tx.Rollback()

	// 1. Получить файл с директорией
	file, err := u.fileTreeRepo.GetFileWithDirectory(ctx, fileID, tx)
	if err != nil {
		return err
	}

	// 2. Проверить статус файла
	if file.Status != "draft" {
		return domain.ErrInvalidFileStatus
	}

	// 3. Создать approval
	approval := &domain.Approval{
		FileID:     file.ID,
		Status:     "on approval",
		WorkflowID: file.Directory.WorkflowID,
		Order:      1,
	}
	if err := u.approvalRepo.CreateApproval(ctx, approval, tx); err != nil {
		return err
	}

	// 4. Обновить статус файла
	file.Status = "approving"
	if err := u.fileTreeRepo.UpdateFileStatus(ctx, file, tx); err != nil {
		return err
	}

	return tx.Commit().Error
}
