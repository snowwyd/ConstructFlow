package interfaces

import (
	"context"
	"service-core/internal/domain"

	"gorm.io/gorm"
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

type FileTreeRepository interface {
	GetDirectoriesWithFiles(ctx context.Context, isArchive bool, userID uint) ([]*domain.Directory, error)
	GetFileInfo(ctx context.Context, fileID uint) (*domain.File, error)

	CreateDirectory(ctx context.Context, parentPathID *uint, name string, status string, userID uint) error
	CreateFile(ctx context.Context, directoryID uint, name string, status string, userID uint) error

	DeleteDirectory(ctx context.Context, directoryID uint, userID uint) error
	DeleteFile(ctx context.Context, fileID uint, userID uint) error

	CheckUserDirectoryAccess(ctx context.Context, userID, directoryID uint) (bool, error)
	CheckUserFileAccess(ctx context.Context, userID, fileID uint) (bool, error)

	WithTx(tx *gorm.DB) FileTreeRepository // Метод для передачи транзакции
	GetDB() *gorm.DB

	GetFileWithDirectory(ctx context.Context, fileID uint, tx *gorm.DB) (*domain.File, error)
	UpdateFileStatus(ctx context.Context, file *domain.File, tx *gorm.DB) error
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
