package usecase

import (
	"backend/internal/domain"
	"backend/internal/domain/interfaces"
	"backend/pkg/logger/slogger"
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
		FileID:        file.ID,
		Status:        "on approval",
		WorkflowID:    file.Directory.WorkflowID,
		WorkflowOrder: 1,
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

// GetApprovalsByUserID получает все Approvals для пользователя через Workflow
func (u *ApprovalUsecase) GetApprovalsByUserID(ctx context.Context, userID uint) ([]domain.ApprovalResponse, error) {
	return u.approvalRepo.FindApprovalsByUser(ctx, userID)
}

func (u *ApprovalUsecase) SignApproval(ctx context.Context, approvalID, userID uint) error {
	const op = "usecase.approval.SignApproval"
	log := u.log.With(slog.String("op", op), slog.Any("login", approvalID))
	log.Info("signing approval")
	// Проверяем, что пользователь имеет право подписывать Approval
	isLastUser, err := u.approvalRepo.IsLastUserInWorkflow(ctx, approvalID, userID)
	if err != nil {
		return err
	}
	if isLastUser {
		return domain.ErrNoPermission
	}

	hasPermission, err := u.approvalRepo.CheckUserPermission(ctx, approvalID, userID)
	if err != nil {
		log.Error("failed to check user permission", slogger.Err(err))
		return err
	}
	if !hasPermission {
		log.Error("user has no permission to approval", slogger.Err(domain.ErrNoPermission))
		return domain.ErrNoPermission
	}

	// Обновляем order
	err = u.approvalRepo.IncrementApprovalOrder(ctx, approvalID)
	if err != nil {
		log.Error("failed to increment order", slogger.Err(err))
	}
	return nil
}

func (u *ApprovalUsecase) AnnotateApproval(ctx context.Context, approvalID, userID uint, message string) error {
	// Проверяем права пользователя
	hasPermission, err := u.approvalRepo.CheckUserPermission(ctx, approvalID, userID)
	if err != nil || !hasPermission {
		return domain.ErrNoPermission
	}

	// Обновляем Approval и File
	return u.approvalRepo.AnnotateApproval(ctx, approvalID, message)
}

func (u *ApprovalUsecase) FinalizeApproval(ctx context.Context, approvalID, userID uint) error {
	// Проверяем права пользователя (последний ли он в цепочке)
	isLastUser, err := u.approvalRepo.IsLastUserInWorkflow(ctx, approvalID, userID)
	if err != nil || !isLastUser {
		return domain.ErrNoPermission
	}

	// Обновляем Approval и File
	return u.approvalRepo.FinalizeApproval(ctx, approvalID)
}
