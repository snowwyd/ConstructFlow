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
