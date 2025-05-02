package interfaces

import (
	"context"
	"service-core/internal/domain"
)

type AuthUsecase interface {
	Login(ctx context.Context, login, password string) (token string, err error)
	GetCurrentUser(ctx context.Context, userID uint) (userInfo domain.GetCurrentUserResponse, err error)
}

type ApprovalUsecase interface {
	ApproveFile(ctx context.Context, fileID uint) error
	GetApprovalsByUserID(ctx context.Context, userID uint) (approvals []domain.ApprovalResponse, err error)
	SignApproval(ctx context.Context, approvalID, userID uint) error
	AnnotateApproval(ctx context.Context, approvalID, userID uint, message string) error
	FinalizeApproval(ctx context.Context, approvalID, userID uint) error
}

type WorkflowUsecase interface {
	GetWorkflows(ctx context.Context, userID uint) (workflows []domain.WorkflowResponse, err error)
	GetWorkflowByID(ctx context.Context, workflowID, userID uint) (domain.ExtendedWorkflowResponse, error)

	CreateWorkflow(ctx context.Context, name string, stages []domain.WorkflowStage, userID uint) error
	UpdateWorkflow(ctx context.Context, workflowID uint, name string, stages []domain.WorkflowStage, userID uint) error
	DeleteWorkflow(ctx context.Context, workflowID uint, userID uint) error

	AssignWorkflow(ctx context.Context, workflowID uint, directoryIDs []uint, userID uint) error
}

type RoleUsecase interface {
	GetRoles(ctx context.Context, userID uint) (workflows []domain.RoleResponse, err error)
	GetRoleByID(ctx context.Context, roleID uint, userID uint) (role string, err error)
	RegisterRole(ctx context.Context, roleName string, userID uint) error
	UpdateRole(ctx context.Context, roleID uint, roleName string, userID uint) error
	DeleteRole(ctx context.Context, roleID uint, userID uint) error
}

type UserUsecase interface {
	GetUsersGrouped(ctx context.Context, userID uint) (users []domain.RoleData, err error)
	RegisterUser(ctx context.Context, login, password string, roleID, userID uint) error
	UpdateUser(ctx context.Context, login, password string, roleID, userID, actorID uint) error
	DeleteUser(ctx context.Context, userID, actorID uint) error

	AssignUser(ctx context.Context, userID uint, directoryIDs []uint, fileIDs []uint, actorID uint) error
}
