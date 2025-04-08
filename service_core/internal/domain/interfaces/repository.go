package interfaces

import (
	"context"
	"service-core/internal/domain"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, userID uint) (user domain.User, err error)
	GetUserByLogin(ctx context.Context, login string) (user domain.User, err error)
	GetUserRole(ctx context.Context, userID uint) (role string, err error)

	SaveUser(ctx context.Context, login string, passHash []byte, roleID uint) error
	CheckUsersExist(ctx context.Context, userIDs []uint) (bool, error)
	CheckUsersWithRole(ctx context.Context, roleID uint) (bool, error)
	GetUsersGroupedByRoles(ctx context.Context) ([]domain.RoleData, error)
}

type ApprovalRepository interface {
	CreateApproval(ctx context.Context, approval *domain.Approval) error
	FindApprovalsByUser(ctx context.Context, userID uint) ([]domain.ApprovalResponse, error)

	IsLastUserInWorkflow(ctx context.Context, approvalID, userID uint) (*domain.Approval, error)
	CheckUserPermission(ctx context.Context, approvalID, userID uint) (*domain.Approval, error)

	IncrementApprovalOrder(ctx context.Context, approvalID uint) error
	AnnotateApproval(ctx context.Context, approvalID uint, message string) error
	FinalizeApproval(ctx context.Context, approvalID uint) error
}

type WorkflowRepository interface {
	GetWorkflows(ctx context.Context) (workflows []domain.WorkflowResponse, err error)
	CreateWorkflow(ctx context.Context, name string, stages []domain.WorkflowStage) error
	UpdateWorkflow(ctx context.Context, workflowID uint, name string, stages []domain.WorkflowStage) error
	DeleteWorkflow(ctx context.Context, workflowID uint) error
	CheckWorkflow(ctx context.Context, workflowID uint) (bool, error)
}

type RoleRepository interface {
	GetRoles(ctx context.Context) (roles []domain.RoleResponse, err error)
	GetRoleByID(ctx context.Context, roleID uint) (roleName string, err error)
	CreateRole(ctx context.Context, roleName string) error
	UpdateRole(ctx context.Context, roleID uint, roleName string) error
	DeleteRole(ctx context.Context, roleID uint) error

	CheckRole(ctx context.Context, roleID uint) (bool, error)
	CheckRoleByName(ctx context.Context, roleName string) (bool, error)
}
