package app

import (
	"backend/internal/domain"
	"context"
)

// Для вызова методов слоя БД для работы с пользователями
type UserRepository interface {
	GetUserByID(ctx context.Context, userID int) (user domain.User, err error)
	GetUserByLogin(ctx context.Context, login string) (user domain.User, err error)

	// Админские ручки для создания/изменения пользователей
	SaveUser(ctx context.Context, login, password, role string) (userID int, err error)
	UpdateUserRole(ctx context.Context, userID int, role string) (success bool, err error)
}

// Для вызова методов слоя БД для работы с папками
type FolderRepository interface {
	SaveFolder(ctx context.Context, parentID int, name string) (folderID int, err error)
	DeleteFolder(ctx context.Context, folderID int) (success bool, err error)

	GetFolder(ctx context.Context, folderID int) (folder domain.Folder, err error)
	GetFolders(ctx context.Context, folderIDs []int) (folders []domain.Folder, err error)
	GetFolderByName(ctx context.Context, parentID int, name string) (folder domain.Folder, err error)
}

// Для вызова методов слоя БД для работы с таблицей users_folders
type UsersFoldersRepository interface {
	GetUserFolderIDs(ctx context.Context, userID int) (folderIDs []int, err error)

	CheckFolderAccess(ctx context.Context, userID, folderID int) (access bool, err error)
}

// Для вызова методов слоя БД для работы с файлами
type FileRepository interface {
	SaveFile(ctx context.Context, folderID int) (fileID int, err error)
	UpdateFile(ctx context.Context, fileID int) (success bool, err error)
	DeleteFile(ctx context.Context, fileID int) (success bool, err error)

	GetFile(ctx context.Context, fileID int) (file domain.File, err error)
	GetFiles(ctx context.Context, fileIDs []int) (files []domain.File, err error)
	GetFolderByName(ctx context.Context, folderID int, name string) (file domain.File, err error)
}

type UsersFilesRepository interface {
	GetUserFileIDs(ctx context.Context, userID int) (fileIDs []int, err error)

	CheckFileAccess(ctx context.Context, userID, fileID int) (access bool, err error)
}

// Для вызова методов слоя БД для работы с файлами на согласовании
type ApprovalRepository interface {
	SaveApproval(ctx context.Context, userID, fileID int, fileName string) (success bool, err error)
	UpdateApproval(ctx context.Context, approvalID, userID int, status string) (success bool, err error)

	GetApprovals(ctx context.Context, userID int) (approvals []domain.Approval, err error)
}

// Для вызова методов слоя БД для работы с процедурами согласования
type WorkflowRepository interface {
	GetUserFromWorkflow(ctx context.Context, workflowID, order int) (userID int, err error)
}

// Для вызова методов слоя БД для работы с аннотациями
type AnnotationRepository interface {
	SaveAnnotation(ctx context.Context, approvalID int, text string) (annotationID int, err error)

	GetAnnotation(ctx context.Context, annotationID int) (annotation domain.Annotation, err error)
}
