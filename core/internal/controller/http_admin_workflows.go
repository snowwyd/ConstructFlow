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

type WorkflowHandler struct {
	usecase interfaces.WorkflowUsecase
}

func NewWorkflowHandler(usecase interfaces.WorkflowUsecase) *WorkflowHandler {
	return &WorkflowHandler{usecase: usecase}
}

// TODO: swagger docs
func (workflowHandler *WorkflowHandler) GetWorkflows(c *gin.Context) {
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

func (workflowHandler *WorkflowHandler) GetWorkflowByID(c *gin.Context) {
	workflowIDStr := c.Param("workflow_id")
	workflowID, err := strconv.ParseUint(workflowIDStr, 10, 64)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "INVALID_APPROVAL_ID", "Invalid workflow ID")
		return
	}

	userID, err := utils.ExtractUserID(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	workflow, err := workflowHandler.usecase.GetWorkflowByID(c.Request.Context(), uint(workflowID), userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrAccessDenied):
			utils.SendErrorResponse(c, http.StatusForbidden, "FORBIDDEN", "User has no access")
		default:
			utils.SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get workflows")
		}
		return
	}

	c.JSON(http.StatusOK, workflow)
}

type workflowInput struct {
	WorkflowName string                 `json:"workflow_name"`
	Stages       []domain.WorkflowStage `json:"stages"`
}

// TODO: swagger docs
func (workflowHandler *WorkflowHandler) CreateWorkflow(c *gin.Context) {
	var req workflowInput
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
		case errors.Is(err, domain.ErrUserNotFound):
			utils.SendErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "Some users are not found")
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

// TODO: swagger docs
func (workflowHandler *WorkflowHandler) DeleteWorkflow(c *gin.Context) {

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
		case errors.Is(err, domain.ErrWorkflowNotFound):
			utils.SendErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "Workflow not found")
		case errors.Is(err, domain.ErrWorkflowInUse):
			utils.SendErrorResponse(c, http.StatusConflict, "CONFILCT", "Workflow is in use")
		default:
			utils.SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete workflow")
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// TODO: swagger docs
func (workflowHandler *WorkflowHandler) UpdateWorkflow(c *gin.Context) {
	workflowIDStr := c.Param("workflow_id")
	workflowID, err := strconv.ParseUint(workflowIDStr, 10, 64)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "INVALID_APPROVAL_ID", "Invalid workflow ID")
		return
	}

	var req workflowInput

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	userID, err := utils.ExtractUserID(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	err = workflowHandler.usecase.UpdateWorkflow(c.Request.Context(), uint(workflowID), req.WorkflowName, req.Stages, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrAccessDenied):
			utils.SendErrorResponse(c, http.StatusForbidden, "FORBIDDEN", "User has no access")
		case errors.Is(err, domain.ErrWorkflowNotFound):
			utils.SendErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "Workflow not found")
		case errors.Is(err, domain.ErrUserNotFound):
			utils.SendErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "Some users are not found")
		default:
			utils.SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update workflow")
		}
		return
	}

	c.Status(http.StatusNoContent)
}

type assignWorkflowInput struct {
	DirectoryIDs []uint `json:"directory_ids"`
}

// TODO: swagger docs
func (workflowHandler *WorkflowHandler) AssignWorkflow(c *gin.Context) {
	workflowIDStr := c.Param("workflow_id")
	workflowID, err := strconv.ParseUint(workflowIDStr, 10, 64)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "INVALID_APPROVAL_ID", "Invalid workflow ID")
		return
	}

	var req assignWorkflowInput

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	userID, err := utils.ExtractUserID(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	err = workflowHandler.usecase.AssignWorkflow(c.Request.Context(), uint(workflowID), req.DirectoryIDs, userID)
	if err != nil {
		switch {
		// TODO: custom errors
		case errors.Is(err, domain.ErrAccessDenied):
			utils.SendErrorResponse(c, http.StatusForbidden, "FORBIDDEN", "User has no access")
		default:
			utils.SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update workflow")
		}
		return
	}

	c.Status(http.StatusNoContent)
}
