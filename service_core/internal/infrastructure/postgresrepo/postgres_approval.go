package postgresrepo

import (
	"context"
	"errors"
	"fmt"
	"service-core/internal/domain"

	"gorm.io/gorm"
)

type ApprovalRepository struct {
	db *gorm.DB
}

func NewApprovalRepository(db *Database) *ApprovalRepository {
	return &ApprovalRepository{db: db.db}
}

func (r *ApprovalRepository) CreateApproval(ctx context.Context, approval *domain.Approval) error {
	return r.db.Create(approval).Error
}

// FindApprovalsByUser находит Approvals через связь с Workflow
func (r *ApprovalRepository) FindApprovalsByUser(ctx context.Context, userID uint) ([]domain.ApprovalResponse, error) {
	const op = "infrastructure.postgresrepo.approval.FindApprovalsByUser"

	var approvals []domain.ApprovalResponse

	err := r.db.WithContext(ctx).
		Table("approvals").
		Select(`
            approvals.id,
            approvals.file_id,
            approvals.status,
            approvals.workflow_order
        `).
		Joins("JOIN workflows ON workflows.workflow_id = approvals.workflow_id AND workflows.workflow_order = approvals.workflow_order").
		Where("workflows.user_id = ?", userID).
		Where("approvals.status = ?", "on approval").
		Scan(&approvals).Error

	return approvals, fmt.Errorf("%s: %w", op, err)
}

// CheckUserPermission проверяет, имеет ли пользователь право подписывать Approval
// Кастомные ошибки: ErrApprovalNotFound
func (r *ApprovalRepository) CheckUserPermission(ctx context.Context, approvalID, userID uint) (*domain.Approval, error) {
	const op = "infrastructure.postgresrepo.approval.CheckUserPermission"

	var approval domain.Approval

	err := r.db.WithContext(ctx).
		Joins("JOIN workflows ON workflows.workflow_id = approvals.workflow_id AND workflows.workflow_order = approvals.workflow_order").
		Where("approvals.id = ? AND workflows.user_id = ?", approvalID, userID).
		Where("approvals.status = ?", "on approval").
		First(&approval).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &domain.Approval{}, fmt.Errorf("%s: %w", op, domain.ErrApprovalNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &approval, nil
}

// IsLastUserInWorkflow - проверка на то, является ли пользователь крайним в цепочке согласования
func (r *ApprovalRepository) IsLastUserInWorkflow(ctx context.Context, approvalID, userID uint) (*domain.Approval, error) {
	const op = "infrastructure.postgresrepo.approval.IsLastUserInWorkflow"

	var approval domain.Approval

	err := r.db.WithContext(ctx).
		Joins("JOIN workflows ON workflows.workflow_id = approvals.workflow_id AND workflows.workflow_order = approvals.workflow_order").
		Where("approvals.id = ?", approvalID).
		Where("workflows.user_id = ?", userID).
		Where("workflows.workflow_order = (SELECT MAX(workflow_order) FROM workflows WHERE workflow_id = approvals.workflow_id)").
		Where("approvals.workflow_order = (SELECT MAX(workflow_order) FROM workflows WHERE workflow_id = approvals.workflow_id)").
		First(&approval).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &domain.Approval{}, nil
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &approval, nil
}

// IncrementApprovalOrder увеличивает поле order у Approval
func (r *ApprovalRepository) IncrementApprovalOrder(ctx context.Context, approvalID uint) error {
	const op = "infrastructure.postgresrepo.approval.IncrementApprovalOrder"

	return fmt.Errorf("%s: %w", op, r.db.WithContext(ctx).
		Model(&domain.Approval{}).
		Where("id = ?", approvalID).
		Update("workflow_order", gorm.Expr("workflow_order + 1")).
		Error)
}

// AnnotateApproval обновляет Approval
func (r *ApprovalRepository) AnnotateApproval(ctx context.Context, approvalID uint, message string) error {
	const op = "infrastructure.postgresrepo.approval.AnnotateApproval"

	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var approval domain.Approval
	if err := tx.First(&approval, approvalID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%s: %w", op, domain.ErrApprovalNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	// Обновляем Approval
	approval.Status = "annotated"
	approval.AnnotationText = message
	if err := tx.Save(&approval).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	return fmt.Errorf("%s: %w", op, tx.Commit().Error)
}

// FinalizeApproval обновляет Approval и связанный File
func (r *ApprovalRepository) FinalizeApproval(ctx context.Context, approvalID uint) error {
	const op = "infrastructure.postgresrepo.approval.FinalizeApproval"

	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var approval domain.Approval
	if err := tx.First(&approval, approvalID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%s: %w", op, domain.ErrApprovalNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	// Обновляем Approval
	approval.Status = "approved"
	if err := tx.Save(&approval).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	return fmt.Errorf("%s: %w", op, tx.Commit().Error)
}
