package http

import (
	"errors"
	"net/http"
	"strconv"

	"service-core/internal/domain"
	"service-core/internal/domain/interfaces"
	"service-core/pkg/utils"

	"github.com/gin-gonic/gin"
)

type FileApprovalsHandler struct {
	usecase interfaces.ApprovalUsecase
}

func NewFileApprovalsHandler(usecase interfaces.ApprovalUsecase) *FileApprovalsHandler {
	return &FileApprovalsHandler{usecase: usecase}
}

// GetApprovalsByUser godoc
// @Summary Получить список согласований пользователя
// @Description Возвращает все согласования, в которых участвует текущий пользователь
// @Tags approval
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {array} domain.Approval "Список согласований"
// @Failure 401 {object} domain.ErrorResponse "Отсутствует/недействителен API-ключ"
// @Failure 500 {object} domain.ErrorResponse "Ошибка при получении данных"
// @Router /file-approvals [get]
func (h *FileApprovalsHandler) GetApprovalsByUser(c *gin.Context) {
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
// @Summary Подписать согласование
// @Description Подтверждает согласование текущим пользователем. Пользователь должен иметь права на подписание.
// @Tags approval
// @Security ApiKeyAuth
// @Param approval_id path string true "ID согласования (числовой формат, например: 456)"
// @Accept json
// @Produce json
// @Success 204 {object} nil "Согласование успешно подписано"
// @Failure 400 {object} domain.ErrorResponse "Невалидный ID согласования"
// @Failure 401 {object} domain.ErrorResponse "Отсутствует/недействителен API-ключ"
// @Failure 403 {object} domain.ErrorResponse "Недостаточно прав или требуется завершение согласования"
// @Failure 404 {object} domain.ErrorResponse "Согласование не найдено"
// @Failure 500 {object} domain.ErrorResponse "Ошибка при обработке подписи"
// @Router /file-approvals/{approval_id}/sign [put]
func (h *FileApprovalsHandler) SignApproval(c *gin.Context) {
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

type annotateApprovalInput struct {
	Message string `json:"message"`
}

// AnnotateApproval godoc
// @Summary Добавить комментарий к согласованию
// @Description Добавляет примечание к согласованию. Пользователь должен участвовать в этом согласовании.
// @Tags approval
// @Security ApiKeyAuth
// @Param approval_id path string true "ID согласования (числовой формат)"
// @Param message body annotateApprovalInput true "Текст примечания"
// @Accept json
// @Produce json
// @Success 204 {object} nil "Примечание добавлено"
// @Failure 400 {object} domain.ErrorResponse "Невалидный ID или тело запроса"
// @Failure 401 {object} domain.ErrorResponse "Отсутствует/недействителен API-ключ"
// @Failure 403 {object} domain.ErrorResponse "Пользователь не участвует в согласовании"
// @Failure 404 {object} domain.ErrorResponse "Согласование не найдено"
// @Failure 500 {object} domain.ErrorResponse "Ошибка при добавлении примечания"
// @Router /file-approvals/{approval_id}/annotate [put]
func (h *FileApprovalsHandler) AnnotateApproval(c *gin.Context) {
	approvalIDStr := c.Param("approval_id")
	approvalID, err := strconv.ParseUint(approvalIDStr, 10, 64)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "INVALID_APPROVAL_ID", "Invalid approval ID")
		return
	}

	var req annotateApprovalInput
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
// @Summary Завершить согласование
// @Description Завершает процесс согласования. Доступно только последнему участнику в цепочке.
// @Tags approval
// @Security ApiKeyAuth
// @Param approval_id path string true "ID согласования (числовой формат)"
// @Accept json
// @Produce json
// @Success 204 {object} nil "Согласование завершено"
// @Failure 400 {object} domain.ErrorResponse "Невалидный ID согласования"
// @Failure 401 {object} domain.ErrorResponse "Отсутствует/недействителен API-ключ"
// @Failure 403 {object} domain.ErrorResponse "Только последний участник может завершить согласование"
// @Failure 404 {object} domain.ErrorResponse "Согласование не найдено"
// @Failure 500 {object} domain.ErrorResponse "Ошибка при завершении согласования"
// @Router /file-approvals/{approval_id}/finalize [put]
func (h *FileApprovalsHandler) FinalizeApproval(c *gin.Context) {
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
