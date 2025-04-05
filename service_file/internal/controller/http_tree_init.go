package http

import "service-file/internal/domain/interfaces"

type TreeHandler struct {
	directoryUsecase interfaces.DirectoryUsecase
	fileUsecase      interfaces.FileUsecase
}

func NewTreeHandler(directoryUsecase interfaces.DirectoryUsecase, fileUsecase interfaces.FileUsecase) *TreeHandler {
	return &TreeHandler{
		directoryUsecase: directoryUsecase,
		fileUsecase:      fileUsecase,
	}
}
