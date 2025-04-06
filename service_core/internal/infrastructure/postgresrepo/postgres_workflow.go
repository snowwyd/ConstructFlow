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

	// Генерация уникального workflow_id
	var newWorkflowID uint
	err := tx.Raw("SELECT nextval('workflows_workflow_id_seq')").Scan(&newWorkflowID).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	// Создаем срез для массовой вставки
	var workflows []domain.Workflow
	for _, stage := range stages {
		workflows = append(workflows, domain.Workflow{
			WorkflowID:    newWorkflowID,
			WorkflowName:  name,
			UserID:        stage.UserID,
			WorkflowOrder: stage.Order,
		})
	}

	// Массовая вставка
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

	// Удаляем все этапы workflow
	result := tx.Where("workflow_id = ?", workflowID).Delete(&domain.Workflow{})
	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("%s: deletion failed: %w", op, result.Error)
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, domain.ErrWorkflowNotFound)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("%s: commit failed: %w", op, err)
	}

	return nil
}
