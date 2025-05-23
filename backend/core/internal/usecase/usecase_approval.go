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
	fileService  interfaces.FileService
	log          *slog.Logger
}

func NewApprovalUsecase(approvalRepo interfaces.ApprovalRepository, fileService interfaces.FileService, log *slog.Logger) *ApprovalUsecase {
	return &ApprovalUsecase{
		approvalRepo: approvalRepo,
		fileService:  fileService,
		log:          log,
	}
}

// GetApprovalsByUserID получает все Approvals для пользователя через Workflow
func (u *ApprovalUsecase) GetApprovalsByUserID(ctx context.Context, userID uint) ([]domain.ApprovalResponse, error) {
	const op = "usecase.approval.GetApprovalsByUserID"

	log := u.log.With(slog.String("op", op), slog.Any("user_id", userID))
	log.Info("getting user approvals")

	approvals, err := u.approvalRepo.FindApprovalsByUser(ctx, userID)
	if err != nil {
		log.Error("failed to get approvals", slogger.Err(err))
		return []domain.ApprovalResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	var fileIDs []uint32
	fileIDSet := make(map[uint32]struct{})
	for _, a := range approvals {
		fileID := uint32(a.FileID)
		if _, exists := fileIDSet[fileID]; !exists {
			fileIDSet[fileID] = struct{}{}
			fileIDs = append(fileIDs, fileID)
		}
	}

	// Получаем имена файлов из файлового сервиса
	log.Debug("getting file names")
	fileNames, err := u.fileService.GetFilesInfo(ctx, fileIDs)
	if err != nil {
		log.Error("failed to get file names from file service", slogger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Обновляем fileName в каждом approval
	for i := range approvals {
		fileID := uint32(approvals[i].FileID)
		if name, exists := fileNames[fileID]; exists {
			approvals[i].FileName = name
		} else {
			log.Warn("file not found in file service", slog.Uint64("file_id", uint64(fileID)))
			approvals[i].FileName = "unknown" // Или пропустить такие записи
		}
	}

	log.Info("approvals got successfully")

	return approvals, nil
}

// ApproveFile создает новую сущность Approval и обновляет статус файла на "approving"
// Кастомные ошибки: ErrInvalidFileStatus, ErrFileNotFound
func (u *ApprovalUsecase) ApproveFile(ctx context.Context, fileID uint) error {
	const op = "usecase.approval.ApproveFile"

	log := u.log.With(slog.String("op", op), slog.Any("file_id", fileID))
	log.Info("starting file approval process")

	log.Debug("fetching file with directory from file service")
	file, err := u.fileService.GetFileWithDirectory(ctx, fileID)
	if err != nil {
		if errors.Is(err, domain.ErrFileNotFound) {
			log.Error("file not found", slogger.Err(err))
			return fmt.Errorf("%s: %w", op, domain.ErrFileNotFound)
		}
		log.Error("failed to get file", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("checking file status", slog.String("current_status", file.Status))
	if file.Status != "draft" {
		log.Warn("invalid file status for approval", slog.String("status", file.Status))
		return fmt.Errorf("%s: %w", op, domain.ErrInvalidFileStatus)
	}

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

	log.Debug("updating file status")
	err = u.fileService.UpdateFileStatus(ctx, fileID, "approving")
	if err != nil {
		if errors.Is(err, domain.ErrFileNotFound) {
			log.Error("file not found", slogger.Err(err))
			return fmt.Errorf("%s: %w", op, domain.ErrFileNotFound)
		}
		log.Error("failed to get file", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("file approval process completed successfully")
	return nil
}

// SignApproval изменяет order в сущности Approval на +1
// Кастомные ошибки: ErrNoPermission, ErrApprovalNotFound
func (u *ApprovalUsecase) SignApproval(ctx context.Context, approvalID, userID uint) error {
	const op = "usecase.approval.SignApproval"
	log := u.log.With(slog.String("op", op), slog.Any("approval_id", approvalID), slog.Any("user_id", userID))
	log.Info("starting approval signing")

	log.Debug("checking if user is last in workflow")
	approval, err := u.approvalRepo.IsLastUserInWorkflow(ctx, approvalID, userID)
	if err != nil {
		log.Error("failed to check workflow position", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	if approval.ID != 0 {
		log.Error("expected finalize but got sign", slogger.Err(domain.ErrNoPermission))
		return fmt.Errorf("%s: %w", op, domain.ErrNoPermission)
	}

	log.Debug("checking user permissions")
	_, err = u.approvalRepo.CheckUserPermission(ctx, approvalID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrApprovalNotFound) {
			log.Error("permission check failed", slogger.Err(domain.ErrApprovalNotFound))
			return fmt.Errorf("%s: %w", op, domain.ErrApprovalNotFound)
		}
		log.Error("permission check failed", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

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

	log.Debug("checking user permissions")
	approval, err := u.approvalRepo.CheckUserPermission(ctx, approvalID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrApprovalNotFound) {
			log.Error("permission check failed", slogger.Err(domain.ErrApprovalNotFound))
			return fmt.Errorf("%s: %w", op, domain.ErrApprovalNotFound)
		}
		log.Error("permission check failed", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("updating approval")
	if err := u.approvalRepo.AnnotateApproval(ctx, approvalID, message); err != nil {
		log.Error("annotation update failed", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("updating file status")
	err = u.fileService.UpdateFileStatus(ctx, approval.FileID, "draft")
	if err != nil {
		if errors.Is(err, domain.ErrFileNotFound) {
			log.Error("file not found", slogger.Err(err))
			return fmt.Errorf("%s: %w", op, domain.ErrFileNotFound)
		}
		log.Error("failed to get file", slogger.Err(err))
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

	log.Debug("checking if user is last in workflow")
	approval, err := u.approvalRepo.IsLastUserInWorkflow(ctx, approvalID, userID)
	if err != nil {
		log.Error("workflow position check failed", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	if approval.ID == 0 {
		log.Error("cannot finalize this approval", slogger.Err(domain.ErrNoPermission))
		return fmt.Errorf("%s: %w", op, domain.ErrNoPermission)
	}

	log.Debug("finalizing approval status")
	if err := u.approvalRepo.FinalizeApproval(ctx, approvalID); err != nil {
		log.Error("finalization failed", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("updating file status")
	err = u.fileService.UpdateFileStatus(ctx, approval.FileID, "approved")
	if err != nil {
		if errors.Is(err, domain.ErrFileNotFound) {
			log.Error("file not found", slogger.Err(err))
			return fmt.Errorf("%s: %w", op, domain.ErrFileNotFound)
		}
		log.Error("failed to get file", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("approval finalized successfully")
	return nil
}
