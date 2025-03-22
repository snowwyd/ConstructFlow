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

// @Router /api/v1/auth/login [post]
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

	c.JSON(http.StatusOK, gin.H{"message": "File sent for approval"})
}

// GetApprovalsByUser обработчик GET /approvals
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

func (h *ApprovalHandler) SignApproval(c *gin.Context) {
	approvalIDStr := c.Param("approval_id")
	approvalID, err := strconv.ParseUint(approvalIDStr, 10, 64)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "INVALID_FILE_ID", "Invalid file ID")
		return
	}

	userID, err := utils.ExtractUserID(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	err = h.usecase.SignApproval(c.Request.Context(), uint(approvalID), userID)
	if err != nil {
		if errors.Is(err, domain.ErrApprovalNotFound) {
			utils.SendErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "Approval not found")
			return
		}
		utils.SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to sign approval")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Approval signed successfully"})
}
