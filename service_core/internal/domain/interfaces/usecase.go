package interfaces

import (
	"context"
	"service-core/internal/domain"
)

type AuthUsecase interface {
	Login(ctx context.Context, login, password string) (token string, err error)
	GetCurrentUser(ctx context.Context, userID uint) (userInfo domain.GetCurrentUserResponse, err error)

	// для админа и локального тестирования
	RegisterUser(ctx context.Context, login, password string, roleID uint) (err error)
	RegisterRole(ctx context.Context, roleName string) (err error)
}

type ApprovalUsecase interface {
	ApproveFile(ctx context.Context, fileID uint) (err error)
	GetApprovalsByUserID(ctx context.Context, userID uint) (approvals []domain.ApprovalResponse, err error)
	SignApproval(ctx context.Context, approvalID, userID uint) error
	AnnotateApproval(ctx context.Context, approvalID, userID uint, message string) error
	FinalizeApproval(ctx context.Context, approvalID, userID uint) error
}

type WorkflowUsecase interface {
	GetWorkflows(ctx context.Context, userID uint) (workflows []domain.WorkflowResponse, err error)
	CreateWorkflow(ctx context.Context, name string, stages []domain.WorkflowStage, userID uint) error
	UpdateWorkflow(ctx context.Context, workflowInfo []domain.Workflow, userID uint) error
	DeleteWorkflow(ctx context.Context, workflowID uint, userID uint) error
}
