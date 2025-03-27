package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"service-core/internal/domain"
	"service-core/internal/domain/interfaces"
	"service-core/pkg/logger/slogger"
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

// GetApprovalsByUserID получает все Approvals для пользователя через Workflow
func (u *ApprovalUsecase) GetApprovalsByUserID(ctx context.Context, userID uint) ([]domain.ApprovalResponse, error) {
	return u.approvalRepo.FindApprovalsByUser(ctx, userID)
}

// ApproveFile создает новую сущность Approval и обновляет статус файла на "approving"
// Кастомные ошибки: ErrInvalidFileStatus, ErrFileNotFound
func (u *ApprovalUsecase) ApproveFile(ctx context.Context, fileID uint) error {
	const op = "usecase.approval.ApproveFile"
	log := u.log.With(slog.String("op", op), slog.Any("file_id", fileID))

	log.Info("starting file approval process")

	// 1. Получить файл с директорией
	// log.Debug("fetching file with directory")
	// file, err := u.fileTreeRepo.GetFileWithDirectory(ctx, fileID, tx)
	// if err != nil {
	// 	if errors.Is(err, domain.ErrFileNotFound) {
	// 		log.Error("file not found", slogger.Err(err))
	// 		return domain.ErrFileNotFound
	// 	}
	// 	log.Error("failed to get file", slogger.Err(err))
	// 	return fmt.Errorf("%s: %w", op, err)
	// }
	// 1. Получить файл с директорией через gRPC
	log.Debug("fetching file with directory from file service")
	file, err := u.fileTreeRepo.GetFileWithDirectory(ctx, fileID)
	if err != nil {
		if errors.Is(err, domain.ErrFileNotFound) {
			log.Error("file not found", slogger.Err(err))
			return domain.ErrFileNotFound
		}
		log.Error("failed to get file", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	// 2. Проверить статус файла
	log.Debug("checking file status", slog.String("current_status", file.Status))
	if file.Status != "draft" {
		log.Warn("invalid file status for approval", slog.String("status", file.Status))
		return domain.ErrInvalidFileStatus
	}

	// 3. Создать approval
	log.Debug("creating approval entity")
	approval := &domain.Approval{
		FileID:        file.ID,
		Status:        "on approval",
		WorkflowID:    file.Directory.WorkflowID,
		WorkflowOrder: 1,
	}
	if err := u.approvalRepo.CreateApproval(ctx, approval); err != nil {
		log.Error("failed to create approval", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	// 4. Обновить статус файла
	// log.Debug("updating file status", slog.String("new_status", "approving"))
	// file.Status = "approving"
	// if err := u.fileTreeRepo.UpdateFileStatus(ctx, file, tx); err != nil {
	// 	log.Error("failed to update file status", slogger.Err(err))
	// 	return fmt.Errorf("%s: %w", op, err)
	// }

	// if err := tx.Commit().Error; err != nil {
	// 	log.Error("transaction commit failed", slogger.Err(err))
	// 	return fmt.Errorf("%s: %w", op, err)
	// }
	fmt.Println(file)
	fmt.Println(approval)
	log.Info("file approval process completed successfully")
	return nil
}

// SignApproval изменяет order в сущности Approval на +1
// Кастомные ошибки: ErrNoPermission, ErrApprovalNotFound
func (u *ApprovalUsecase) SignApproval(ctx context.Context, approvalID, userID uint) error {
	const op = "usecase.approval.SignApproval"
	log := u.log.With(slog.String("op", op), slog.Any("approval_id", approvalID), slog.Any("user_id", userID))
	log.Info("starting approval signing")

	// Проверка последнего пользователя в workflow
	log.Debug("checking if user is last in workflow")
	isLastUser, err := u.approvalRepo.IsLastUserInWorkflow(ctx, approvalID, userID)
	if err != nil {
		log.Error("failed to check workflow position", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	if isLastUser {
		log.Warn("user is last in workflow and cannot sign", slogger.Err(domain.ErrNoPermission))
		return domain.ErrNoPermission
	}

	// Проверка прав доступа
	log.Debug("checking user permissions")
	hasPermission, err := u.approvalRepo.CheckUserPermission(ctx, approvalID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrApprovalNotFound) {
			log.Error("permission check failed", slogger.Err(domain.ErrApprovalNotFound))
			return domain.ErrApprovalNotFound
		}
		log.Error("permission check failed", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	if !hasPermission {
		log.Warn("user has no permission to sign approval")
		return domain.ErrNoPermission
	}

	// Обновление порядка
	log.Debug("incrementing approval order")
	if err := u.approvalRepo.IncrementApprovalOrder(ctx, approvalID); err != nil {
		log.Error("failed to increment approval order", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("approval signed successfully")
	return nil
}

// AnnotateApproval меняет статус и сообщение в сущности Approval и возвращает файлу статус "draft"
// Кастомные ошибки: ErrNoPermission, ErrApprovalNotFound
func (u *ApprovalUsecase) AnnotateApproval(ctx context.Context, approvalID, userID uint, message string) error {
	const op = "usecase.approval.AnnotateApproval"
	log := u.log.With(slog.String("op", op), slog.Any("approval_id", approvalID), slog.Any("user_id", userID))
	log.Info("adding annotation to approval")

	// Проверка прав
	log.Debug("checking user permissions")
	hasPermission, err := u.approvalRepo.CheckUserPermission(ctx, approvalID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrApprovalNotFound) {
			log.Error("permission check failed", slogger.Err(domain.ErrApprovalNotFound))
			return domain.ErrApprovalNotFound
		}
		log.Error("permission check failed", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	if !hasPermission {
		log.Warn("user has no permission to annotate")
		return domain.ErrNoPermission
	}

	// Обновление аннотации
	log.Debug("updating approval annotation")
	if err := u.approvalRepo.AnnotateApproval(ctx, approvalID, message); err != nil {
		log.Error("annotation update failed", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("annotation added successfully")
	return nil
}

// FinalizeApproval меняет статус сущностей Approval и File на approved
// Кастомные ошибки: ErrNoPermission
func (u *ApprovalUsecase) FinalizeApproval(ctx context.Context, approvalID, userID uint) error {
	const op = "usecase.approval.FinalizeApproval"
	log := u.log.With(slog.String("op", op), slog.Any("approval_id", approvalID), slog.Any("user_id", userID))
	log.Info("finalizing approval")

	// Проверка последнего пользователя
	log.Debug("checking if user is last in workflow")
	isLastUser, err := u.approvalRepo.IsLastUserInWorkflow(ctx, approvalID, userID)
	if err != nil {
		log.Error("workflow position check failed", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	if !isLastUser {
		log.Warn("only last user can finalize approval")
		return domain.ErrNoPermission
	}

	log.Debug("finalizing approval status")
	if err := u.approvalRepo.FinalizeApproval(ctx, approvalID); err != nil {
		log.Error("finalization failed", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("approval finalized successfully")
	return nil
}
