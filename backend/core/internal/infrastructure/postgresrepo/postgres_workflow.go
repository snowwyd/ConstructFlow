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

func (workflowRepo *WorkflowRepository) GetWorkflowByID(ctx context.Context, workflowID uint) (domain.ExtendedWorkflowResponse, error) {
	const op = "infrastructure.postgresrepo.workflow.GetWorkflowByID"

	type WorkflowStageRow struct {
		UserID       uint   `json:"user_id"`
		Order        int    `json:"order"`
		WorkflowName string `json:"workflow_name"`
	}

	var rows []WorkflowStageRow

	err := workflowRepo.db.WithContext(ctx).
		Table("workflows").
		Select("workflows.user_id, workflows.workflow_order AS order, workflows.workflow_name").
		Where("workflows.workflow_id = ?", workflowID).
		Order("workflows.workflow_order ASC").
		Scan(&rows).Error

	if err != nil {
		return domain.ExtendedWorkflowResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	if len(rows) == 0 {
		return domain.ExtendedWorkflowResponse{}, fmt.Errorf("%s: %w", op, domain.ErrWorkflowNotFound)
	}

	var workflow domain.ExtendedWorkflowResponse
	workflow.WorkflowName = rows[0].WorkflowName
	workflow.Stages = make([]domain.WorkflowStage, len(rows))

	for i, row := range rows {
		workflow.Stages[i] = domain.WorkflowStage{
			UserID: row.UserID,
			Order:  row.Order,
		}
	}

	return workflow, nil
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

func (workflowRepo *WorkflowRepository) UpdateWorkflow(ctx context.Context, workflowID uint, name string, stages []domain.WorkflowStage) error {
	const op = "infrastructure.postgresrepo.workflow.UpdateWorkflow"

	tx := workflowRepo.db.WithContext(ctx).Begin()
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

func (workflowRepo *WorkflowRepository) CheckUserInWorkflow(ctx context.Context, userID uint) (bool, error) {
	const op = "infrastructure.postgresrepo.workflow.CheckUserInWorkflow"

	var count int64
	err := workflowRepo.db.WithContext(ctx).
		Model(&domain.Workflow{}).
		Where("user_id = ?", userID).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return count > 0, nil
}
