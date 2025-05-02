package http

import (
	"errors"
	"net/http"
	"service-core/internal/domain"
	"service-core/internal/domain/interfaces"
	"service-core/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	usecase interfaces.ApprovalUsecase
}

func NewFileHandler(usecase interfaces.ApprovalUsecase) *FileHandler {
	return &FileHandler{usecase: usecase}
}

// ApproveFile godoc
// @Summary Отправить файл на согласование
// @Description Переводит файл в статус "на согласовании". Файл должен находиться в состоянии черновика.
// @Tags approval
// @Security ApiKeyAuth
// @Param file_id path string true "ID файла (числовой формат, например: 123)"
// @Accept json
// @Produce json
// @Success 201 {object} nil "Файл успешно отправлен на согласование"
// @Failure 400 {object} domain.ErrorResponse "Невалидный ID файла или файл не в статусе 'черновик'"
// @Failure 404 {object} domain.ErrorResponse "Файл с указанным ID не найден"
// @Failure 500 {object} domain.ErrorResponse "Ошибка при изменении статуса файла"
// @Router /files/{file_id}/approve [put]
func (h *FileHandler) ApproveFile(c *gin.Context) {
	_, err := utils.ExtractUserID(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	fileIDStr := c.Param("file_id")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 64)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "INVALID_FILE_ID", "Invalid file ID")
		return
	}

	err = h.usecase.ApproveFile(c.Request.Context(), uint(fileID))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrFileNotFound):
			utils.SendErrorResponse(c, http.StatusNotFound, "FILE_NOT_FOUND", "File not found")
		case errors.Is(err, domain.ErrInvalidFileStatus):
			utils.SendErrorResponse(c, http.StatusBadRequest, "INVALID_STATUS", "File is not in a draft state")
		default:
			utils.SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to approve file")
		}
		return
	}

	c.Status(http.StatusCreated)
}
