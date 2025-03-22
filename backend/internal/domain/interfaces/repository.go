package interfaces

import (
	"backend/internal/domain"
	"context"

	"gorm.io/gorm"
)

// Для вызова методов слоя БД для работы с пользователями
type UserRepository interface {
	GetUserByID(ctx context.Context, userID uint) (user domain.User, err error)
	GetUserByLogin(ctx context.Context, login string) (user domain.User, err error)

	// Админские ручки для создания/изменения пользователей
	SaveUser(ctx context.Context, login string, passHash []byte, roleID uint) (userID uint, err error)
}

// Для вызова методов слоя БД для работы с ролями
type RoleRepository interface {
	CreateRole(ctx context.Context, roleName string) (roleID uint, err error)
	GetRoleByID(ctx context.Context, roleID uint) (roleName string, err error)
}

// Для вызова методов слоя БД для работы с папками
// type DirectoryRepository interface {
// 	SaveDirectory(ctx context.Context, parentID int, name string) (directoryID int, err error)
// 	DeleteDirectory(ctx context.Context, directoryID int) (success bool, err error)

// 	GetDirectory(ctx context.Context, directoryID int) (directory domain.Directory, err error)
// 	GetDirectorys(ctx context.Context, directoryIDs []int) (directorys []domain.Directory, err error)
// 	GetDirectoryByName(ctx context.Context, parentID int, name string) (directory domain.Directory, err error)
// }

// // Для вызова методов слоя БД для работы с таблицей users_directorys
// type UsersDirectorysRepository interface {
// 	GetUserDirectoryIDs(ctx context.Context, userID int) (directoryIDs []int, err error)

// 	CheckDirectoryAccess(ctx context.Context, userID, directoryID int) (access bool, err error)
// }

// // Для вызова методов слоя БД для работы с файлами
// type FileRepository interface {
// 	SaveFile(ctx context.Context, directoryID int) (fileID int, err error)
// 	UpdateFile(ctx context.Context, fileID int) (success bool, err error)
// 	DeleteFile(ctx context.Context, fileID int) (success bool, err error)

// 	GetFile(ctx context.Context, fileID int) (file domain.File, err error)
// 	GetFiles(ctx context.Context, fileIDs []int) (files []domain.File, err error)
// 	GetDirectoryByName(ctx context.Context, directoryID int, name string) (file domain.File, err error)
// }

// type UsersFilesRepository interface {
// 	GetUserFileIDs(ctx context.Context, userID int) (fileIDs []int, err error)

// 	CheckFileAccess(ctx context.Context, userID, fileID int) (access bool, err error)
// }

// // Для вызова методов слоя БД для работы с процедурами согласования
// type WorkflowRepository interface {
// 	GetUserFromWorkflow(ctx context.Context, workflowID, order int) (userID int, err error)
// }

// // Для вызова методов слоя БД для работы с аннотациями
// type AnnotationRepository interface {
// 	SaveAnnotation(ctx context.Context, approvalID int, text string) (annotationID int, err error)

// 	GetAnnotation(ctx context.Context, annotationID int) (annotation domain.Annotation, err error)
// }

type FileTreeRepository interface {
	GetDirectoriesWithFiles(ctx context.Context, isArchive bool, userID uint) ([]*domain.Directory, error)
	GetFileInfo(ctx context.Context, fileID uint) (*domain.File, error)

	CreateDirectory(ctx context.Context, parentPathID *uint, name string, status string, userID uint) (uint, error)
	CreateFile(ctx context.Context, directoryID uint, name string, status string, userID uint) (uint, error)

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
	CreateApproval(ctx context.Context, approval *domain.Approval, tx *gorm.DB) error
	FindApprovalsByUser(ctx context.Context, userID uint) ([]domain.ApprovalResponse, error)

	IsLastUserInWorkflow(ctx context.Context, approvalID, userID uint) (bool, error)
	CheckUserPermission(ctx context.Context, approvalID, userID uint) (bool, error)

	IncrementApprovalOrder(ctx context.Context, approvalID uint) error
	AnnotateApproval(ctx context.Context, approvalID uint, message string) error
	FinalizeApproval(ctx context.Context, approvalID uint) error
}
