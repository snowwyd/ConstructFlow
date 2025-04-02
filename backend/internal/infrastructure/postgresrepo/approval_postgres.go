package postgresrepo

import (
	"backend/internal/domain"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type ApprovalRepository struct {
	db *gorm.DB
}

func NewApprovalRepository(db *Database) *ApprovalRepository {
	return &ApprovalRepository{db: db.db}
}

func (r *ApprovalRepository) CreateApproval(ctx context.Context, approval *domain.Approval, tx *gorm.DB) error {
	return tx.WithContext(ctx).Create(approval).Error
}

// FindApprovalsByUser находит Approvals через связь с Workflow
func (r *ApprovalRepository) FindApprovalsByUser(ctx context.Context, userID uint) ([]domain.ApprovalResponse, error) {
	var approvals []domain.ApprovalResponse

	err := r.db.WithContext(ctx).
		Table("approvals").
		Select(`
            approvals.id,
            approvals.file_id,
            files.name AS file_name,
            approvals.status,
            approvals.workflow_order,
            (SELECT MAX(workflow_order) FROM workflows WHERE workflows.workflow_id = approvals.workflow_id) AS workflow_user_count
        `).
		Joins("JOIN workflows ON workflows.workflow_id = approvals.workflow_id AND workflows.workflow_order = approvals.workflow_order").
		Joins("JOIN files ON files.id = approvals.file_id").
		Where("workflows.user_id = ?", userID).
		Where("approvals.status = ?", "on approval").
		Scan(&approvals).Error

	return approvals, err
}

// CheckUserPermission проверяет, имеет ли пользователь право подписывать Approval
// Кастомные ошибки: ErrApprovalNotFound
func (r *ApprovalRepository) CheckUserPermission(ctx context.Context, approvalID, userID uint) (bool, error) {
	// Проверка существования Approval
	var approvalExists bool
	err := r.db.WithContext(ctx).
		Model(&domain.Approval{}).
		Select("COUNT(*) > 0").
		Where("id = ?", approvalID).
		Scan(&approvalExists).Error
	if err != nil {
		return false, fmt.Errorf("existence check failed: %w", domain.ErrInternal)
	}
	if !approvalExists {
		return false, domain.ErrApprovalNotFound
	}

	// Проверка прав доступа
	var hasPermission bool
	err = r.db.WithContext(ctx).
		Model(&domain.Approval{}).
		Joins("JOIN workflows ON workflows.workflow_id = approvals.workflow_id AND workflows.workflow_order = approvals.workflow_order").
		Where("approvals.id = ? AND workflows.user_id = ?", approvalID, userID).
		Where("approvals.status = ?", "on approval").
		Select("COUNT(*) > 0").
		Scan(&hasPermission).Error
	if err != nil {
		return false, fmt.Errorf("permission check failed: %w", domain.ErrInternal)
	}

	return hasPermission, nil
}

// IncrementApprovalOrder увеличивает поле order у Approval
func (r *ApprovalRepository) IncrementApprovalOrder(ctx context.Context, approvalID uint) error {
	return r.db.WithContext(ctx).
		Model(&domain.Approval{}).
		Where("id = ?", approvalID).
		Update("workflow_order", gorm.Expr("workflow_order + 1")).
		Error
}

// AnnotateApproval обновляет Approval и связанный File
func (r *ApprovalRepository) AnnotateApproval(ctx context.Context, approvalID uint, message string) error {
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Находим Approval и связанный File
	var approval domain.Approval
	if err := tx.Preload("File").First(&approval, approvalID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Обновляем Approval
	approval.Status = "annotated"
	approval.AnnotationText = message
	if err := tx.Save(&approval).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Обновляем File
	if approval.File.ID != 0 {
		approval.File.Status = "draft"
		if err := tx.Save(approval.File).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// IsLastUserInWorkflow - проверка на то, является ли пользователь крайним в цепочке согласования
func (r *ApprovalRepository) IsLastUserInWorkflow(ctx context.Context, approvalID, userID uint) (bool, error) {
	var isLast bool
	err := r.db.WithContext(ctx).
		Model(&domain.Approval{}).
		Joins("JOIN workflows ON workflows.workflow_id = approvals.workflow_id AND workflows.workflow_order = approvals.workflow_order").
		Where("approvals.id = ?", approvalID).
		Select(`
            workflows.user_id = ? AND 
            workflows.workflow_order = (SELECT MAX(workflow_order) FROM workflows WHERE workflow_id = approvals.workflow_id)
        `, userID).
		Scan(&isLast).Error
	return isLast, err
}

// FinalizeApproval обновляет Approval и связанный File
func (r *ApprovalRepository) FinalizeApproval(ctx context.Context, approvalID uint) error {
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Находим Approval и связанный File
	var approval domain.Approval
	if err := tx.Preload("File").First(&approval, approvalID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Обновляем Approval
	approval.Status = "approved"
	if err := tx.Save(&approval).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Обновляем File
	if approval.File.ID != 0 {
		approval.File.Status = "approved"
		if err := tx.Save(approval.File).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
