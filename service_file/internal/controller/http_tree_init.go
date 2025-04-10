package http

import "service-file/internal/domain/interfaces"

type TreeHandler struct {
	directoryUsecase interfaces.DirectoryUsecase
	fileUsecase      interfaces.FileUsecase
	adminUsecase     interfaces.AdminUsecase
}

func NewTreeHandler(directoryUsecase interfaces.DirectoryUsecase, fileUsecase interfaces.FileUsecase, adminUsecase interfaces.AdminUsecase) *TreeHandler {
	return &TreeHandler{
		directoryUsecase: directoryUsecase,
		fileUsecase:      fileUsecase,
		adminUsecase:     adminUsecase,
	}
}
