package http

import (
	"errors"
	"net/http"
	"service-core/internal/domain"
	"service-core/internal/domain/interfaces"
	"service-core/pkg/utils"

	"github.com/gin-gonic/gin"
)

type WorkflowlHandler struct {
	usecase interfaces.WorkflowUsecase
}

func NewWorkflowHandler(usecase interfaces.WorkflowUsecase) *WorkflowlHandler {
	return &WorkflowlHandler{usecase: usecase}
}

func (workflowHandler *WorkflowlHandler) GetWorkflows(c *gin.Context) {
	userID, err := utils.ExtractUserID(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	workflows, err := workflowHandler.usecase.GetWorkflows(c.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrAccessDenied):
			utils.SendErrorResponse(c, http.StatusForbidden, "FORBIDDEN", "User has no access")
		default:
			utils.SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get workflows")
		}
		return
	}

	c.JSON(http.StatusOK, workflows)
}

type createWorkflowInput struct {
	WorkflowName string                 `json:"workflow_name"`
	Stages       []domain.WorkflowStage `json:"stages"`
}

func (workflowHandler *WorkflowlHandler) CreateWorkflow(c *gin.Context) {

	var req createWorkflowInput
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	userID, err := utils.ExtractUserID(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	err = workflowHandler.usecase.CreateWorkflow(c.Request.Context(), req.WorkflowName, req.Stages, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrAccessDenied):
			utils.SendErrorResponse(c, http.StatusForbidden, "FORBIDDEN", "User has no access")
		default:
			utils.SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get workflows")
		}
		return
	}

	c.Status(http.StatusCreated)
}

type deleteWorkflowInput struct {
	WorkflowID uint `json:"workflow_id"`
}

func (workflowHandler *WorkflowlHandler) DeleteWorkflow(c *gin.Context) {

	var req deleteWorkflowInput
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	userID, err := utils.ExtractUserID(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	err = workflowHandler.usecase.DeleteWorkflow(c.Request.Context(), req.WorkflowID, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrAccessDenied):
			utils.SendErrorResponse(c, http.StatusForbidden, "FORBIDDEN", "User has no access")
		default:
			utils.SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get workflows")
		}
		return
	}

	c.Status(http.StatusCreated)
}
