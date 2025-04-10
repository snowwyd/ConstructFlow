package http

import (
	"errors"
	"net/http"
	"strconv"

	"service-file/internal/domain"
	"service-file/pkg/utils"

	"github.com/gin-gonic/gin"
)

func (h *TreeHandler) GetUserTree(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid file ID")
		return
	}

	actorID, err := utils.ExtractUserID(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	response, err := h.adminUsecase.GetUserTree(c.Request.Context(), uint(userID), actorID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrAccessDenied):
			utils.SendErrorResponse(c, http.StatusForbidden, "ACCESS_DENIED", "You do not have access to this file")
		default:
			utils.SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get file info")
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *TreeHandler) GetWorkflowTree(c *gin.Context) {
	workflowIDStr := c.Param("workflow_id")
	workflowID, err := strconv.ParseUint(workflowIDStr, 10, 64)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid file ID")
		return
	}

	actorID, err := utils.ExtractUserID(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	response, err := h.adminUsecase.GetWorkflowTree(c.Request.Context(), uint(workflowID), actorID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrAccessDenied):
			utils.SendErrorResponse(c, http.StatusForbidden, "ACCESS_DENIED", "You do not have access to this file")
		default:
			utils.SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get file info")
		}
		return
	}

	c.JSON(http.StatusOK, response)
}
