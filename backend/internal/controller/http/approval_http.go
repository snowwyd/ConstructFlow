package http

import (
	"errors"
	"net/http"
	"strconv"

	"backend/internal/domain"
	"backend/internal/domain/interfaces"
	"backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

type ApprovalHandler struct {
	usecase interfaces.ApprovalUsecase
}

// конструктор
func NewApprovalHandler(usecase interfaces.ApprovalUsecase) *ApprovalHandler {
	return &ApprovalHandler{usecase: usecase}
}

// ApproveFile godoc
// @Summary Одобрить файл
// @Description Отправляет файл на одобрение
// @Tags approval
// @Param file_id path uint true "ID файла"
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Файл отправлен на одобрение"
// @Failure 400 {object} domain.ErrorResponse "Неверный ID файла или файл не в состоянии черновика"
// @Failure 404 {object} domain.ErrorResponse "Файл не найден"
// @Failure 500 {object} domain.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/approval/{file_id}/approve [post]
func (h *ApprovalHandler) ApproveFile(c *gin.Context) {
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

// GetApprovalsByUser godoc
// @Summary Получить одобрения пользователя
// @Description Возвращает список одобрений для текущего пользователя
// @Tags approval
// @Accept json
// @Produce json
// @Success 200 {array} domain.Approval "Список одобрений"
// @Failure 401 {object} domain.ErrorResponse "Пользователь не аутентифицирован"
// @Failure 500 {object} domain.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/approvals [get]
func (h *ApprovalHandler) GetApprovalsByUser(c *gin.Context) {
	userID, err := utils.ExtractUserID(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	approvals, err := h.usecase.GetApprovalsByUserID(c.Request.Context(), userID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch approvals")
		return
	}

	c.JSON(http.StatusOK, approvals)
}

// SignApproval godoc
// @Summary Подписать одобрение
// @Description Подписание одобрения указанным пользователем
// @Tags approval
// @Param approval_id path uint true "ID одобрения"
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Одобрение подписано успешно"
// @Failure 400 {object} domain.ErrorResponse "Неверный ID одобрения"
// @Failure 401 {object} domain.ErrorResponse "Пользователь не аутентифицирован"
// @Failure 403 {object} domain.ErrorResponse "Нет прав для подписания или требуется завершение"
// @Failure 404 {object} domain.ErrorResponse "Одобрение не найдено"
// @Failure 500 {object} domain.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/approval/{approval_id}/sign [post]
func (h *ApprovalHandler) SignApproval(c *gin.Context) {
	approvalIDStr := c.Param("approval_id")
	approvalID, err := strconv.ParseUint(approvalIDStr, 10, 64)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "INVALID_APPROVAL_ID", "Invalid approval ID")
		return
	}

	userID, err := utils.ExtractUserID(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	err = h.usecase.SignApproval(c.Request.Context(), uint(approvalID), userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrApprovalNotFound):
			utils.SendErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "Approval not found")
		case errors.Is(err, domain.ErrNoPermission):
			utils.SendErrorResponse(c, http.StatusForbidden, "ACCESS_DENIED", "User has no permission to sign approval or need to finalize it")
		default:
			utils.SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to sign approval")
		}

		return
	}
	c.Status(http.StatusNoContent)
}

// AnnotateApproval godoc
// @Summary Добавить аннотацию к одобрению
// @Description Добавляет сообщение (аннотацию) к одобрению
// @Tags approval
// @Param approval_id path uint true "ID одобрения"
// @Param message body string true "Текст аннотации"
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Аннотация добавлена успешно"
// @Failure 400 {object} domain.ErrorResponse "Неверный ID одобрения или тело запроса"
// @Failure 401 {object} domain.ErrorResponse "Пользователь не аутентифицирован"
// @Failure 403 {object} domain.ErrorResponse "Нет прав для добавления аннотации"
// @Failure 404 {object} domain.ErrorResponse "Одобрение не найдено"
// @Failure 500 {object} domain.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/approval/{approval_id}/annotate [post]
func (h *ApprovalHandler) AnnotateApproval(c *gin.Context) {
	approvalIDStr := c.Param("approval_id")
	approvalID, err := strconv.ParseUint(approvalIDStr, 10, 64)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "INVALID_APPROVAL_ID", "Invalid approval ID")
		return
	}

	var req struct {
		Message string `json:"message"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	userID, err := utils.ExtractUserID(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	err = h.usecase.AnnotateApproval(c.Request.Context(), uint(approvalID), userID, req.Message)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrApprovalNotFound):
			utils.SendErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "Approval not found")
		case errors.Is(err, domain.ErrNoPermission):
			utils.SendErrorResponse(c, http.StatusForbidden, "FORBIDDEN", "User has no permission to annotate this approval")
		default:
			utils.SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to annotate approval")
		}
		return
	}
	c.Status(http.StatusNoContent)
}

// FinalizeApproval godoc
// @Summary Завершить одобрение
// @Description Завершает процесс одобрения для указанного одобрения
// @Tags approval
// @Param approval_id path uint true "ID одобрения"
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Одобрение завершено успешно"
// @Failure 400 {object} domain.ErrorResponse "Неверный ID одобрения"
// @Failure 401 {object} domain.ErrorResponse "Пользователь не аутентифицирован"
// @Failure 403 {object} domain.ErrorResponse "Только последний пользователь может завершить"
// @Failure 404 {object} domain.ErrorResponse "Одобрение не найдено"
// @Failure 500 {object} domain.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/approval/{approval_id}/finalize [post]
func (h *ApprovalHandler) FinalizeApproval(c *gin.Context) {
	approvalIDStr := c.Param("approval_id")
	approvalID, err := strconv.ParseUint(approvalIDStr, 10, 64)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "INVALID_APPROVAL_ID", "Invalid approval ID")
		return
	}

	userID, err := utils.ExtractUserID(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	err = h.usecase.FinalizeApproval(c.Request.Context(), uint(approvalID), userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrApprovalNotFound):
			utils.SendErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "Approval not found")
		case errors.Is(err, domain.ErrNoPermission):
			utils.SendErrorResponse(c, http.StatusForbidden, "FORBIDDEN", "Only the last user in the workflow can finalize this approval")
		default:
			utils.SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to finalize approval")
		}
		return
	}
	c.Status(http.StatusNoContent)
}
