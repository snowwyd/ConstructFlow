package interfaces

import (
	"context"
	"service-core/internal/domain"
)

// Для вызова методов слоя БД для работы с пользователями
type UserRepository interface {
	GetUserByID(ctx context.Context, userID uint) (user domain.User, err error)
	GetUserByLogin(ctx context.Context, login string) (user domain.User, err error)

	// Админские ручки для создания/изменения пользователей
	SaveUser(ctx context.Context, login string, passHash []byte, roleID uint) (err error)
}

// Для вызова методов слоя БД для работы с ролями
type RoleRepository interface {
	CreateRole(ctx context.Context, roleName string) (err error)
	GetRoleByID(ctx context.Context, roleID uint) (roleName string, err error)
}

// Для вызова методов слоя БД для работы с файлами на согласовании
type ApprovalRepository interface {
	CreateApproval(ctx context.Context, approval *domain.Approval) error
	FindApprovalsByUser(ctx context.Context, userID uint) ([]domain.ApprovalResponse, error)

	IsLastUserInWorkflow(ctx context.Context, approvalID, userID uint) (*domain.Approval, error)
	CheckUserPermission(ctx context.Context, approvalID, userID uint) (*domain.Approval, error)

	IncrementApprovalOrder(ctx context.Context, approvalID uint) error
	AnnotateApproval(ctx context.Context, approvalID uint, message string) error
	FinalizeApproval(ctx context.Context, approvalID uint) error
}
