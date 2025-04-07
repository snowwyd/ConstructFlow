package postgresrepo

import (
	"context"
	"fmt"
	"service-core/internal/domain"

	"gorm.io/gorm"
)

type WorkflowRepository struct {
	db *gorm.DB
}

func NewWorkflowRepository(db *Database) *WorkflowRepository {
	return &WorkflowRepository{db: db.db}
}

func (workflowRepo *WorkflowRepository) GetWorkflows(ctx context.Context) ([]domain.WorkflowResponse, error) {
	const op = "infrastructure.postgresrepo.workflow.GetWorkflows"

	var workflows []domain.WorkflowResponse
	err := workflowRepo.db.WithContext(ctx).
		Table("workflows").
		Select(`
        workflow_id, 
        MAX(workflow_name) AS workflow_name, 
        MAX(workflow_order) AS workflow_length
    `).
		Where("deleted_at IS NULL").
		Group("workflow_id").
		Having("COUNT(CASE WHEN deleted_at IS NOT NULL THEN 1 END) = 0").
		Order("workflow_name ASC").
		Scan(&workflows).Error

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return workflows, nil
}

func (workflowRepo *WorkflowRepository) CreateWorkflow(ctx context.Context, name string, stages []domain.WorkflowStage) error {
	const op = "infrastructure.postgresrepo.workflow.CreateWorkflow"

	tx := workflowRepo.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var newWorkflowID uint
	err := tx.Raw("SELECT nextval('workflows_workflow_id_seq')").Scan(&newWorkflowID).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	var workflows []domain.Workflow
	for _, stage := range stages {
		workflows = append(workflows, domain.Workflow{
			WorkflowID:    newWorkflowID,
			WorkflowName:  name,
			UserID:        stage.UserID,
			WorkflowOrder: stage.Order,
		})
	}

	if err := tx.Create(&workflows).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (workflowRepo *WorkflowRepository) DeleteWorkflow(ctx context.Context, workflowID uint) error {
	const op = "infrastructure.postgresrepo.workflow.DeleteWorkflow"

	tx := workflowRepo.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	result := tx.Where("workflow_id = ?", workflowID).Delete(&domain.Workflow{})
	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, result.Error)
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, domain.ErrWorkflowNotFound)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (workflowRepository *WorkflowRepository) UpdateWorkflow(ctx context.Context, workflowID uint, name string, stages []domain.WorkflowStage) error {
	const op = "infrastructure.postgresrepo.workflow.UpdateWorkflow"

	tx := workflowRepository.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Where("workflow_id = ?", workflowID).Delete(&domain.Workflow{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	var workflows []domain.Workflow
	for _, stage := range stages {
		workflows = append(workflows, domain.Workflow{
			WorkflowID:    workflowID,
			WorkflowName:  name,
			UserID:        stage.UserID,
			WorkflowOrder: stage.Order,
		})
	}

	if err := tx.Create(&workflows).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *WorkflowRepository) CheckWorkflow(ctx context.Context, workflowID uint) (bool, error) {
	const op = "infrastructure.postgresrepo.workflow.CheckWorkflow"

	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.Workflow{}).
		Where("workflow_id = ?", workflowID).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return count > 0, nil
}
