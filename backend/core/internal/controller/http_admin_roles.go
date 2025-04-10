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

type RoleHandler struct {
	usecase interfaces.RoleUsecase
}

func NewRoleHandler(usecase interfaces.RoleUsecase) *RoleHandler {
	return &RoleHandler{usecase: usecase}
}

type registerRoleInput struct {
	RoleName string `json:"role_name"`
}

// TODO: swagger docs
func (roleHandler *RoleHandler) RegisterRole(c *gin.Context) {
	var req registerRoleInput

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	userID, err := utils.ExtractUserID(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	err = roleHandler.usecase.RegisterRole(c.Request.Context(), req.RoleName, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRoleAlreadyExists):
			utils.SendErrorResponse(c, http.StatusConflict, "ROLE_ALREADY_EXISTS", "Role with this name already exists")
		default:
			utils.SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to register role")
		}
		return
	}

	c.Status(http.StatusCreated)
}

// TODO: swagger docs
func (roleHandler *RoleHandler) GetRoles(c *gin.Context) {
	userID, err := utils.ExtractUserID(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	roles, err := roleHandler.usecase.GetRoles(c.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrAccessDenied):
			utils.SendErrorResponse(c, http.StatusForbidden, "FORBIDDEN", "User has no access")
		default:
			utils.SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get roles")
		}
		return
	}

	c.JSON(http.StatusOK, roles)
}

// TODO: swagger docs
func (roleHandler *RoleHandler) GetRole(c *gin.Context) {
	roleIDStr := c.Param("role_id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 64)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "INVALID_APPROVAL_ID", "Invalid approval ID")
		return
	}
	userID, err := utils.ExtractUserID(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	role, err := roleHandler.usecase.GetRoleByID(c.Request.Context(), uint(roleID), userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrAccessDenied):
			utils.SendErrorResponse(c, http.StatusForbidden, "FORBIDDEN", "User has no access")
		default:
			utils.SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get role by id")
		}
		return
	}

	c.JSON(http.StatusOK, role)
}

type deleteRoleInput struct {
	RoleID uint `json:"role_id"`
}

// TODO: swagger docs
func (roleHandler *RoleHandler) DeleteRole(c *gin.Context) {

	var req deleteRoleInput
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	userID, err := utils.ExtractUserID(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	err = roleHandler.usecase.DeleteRole(c.Request.Context(), req.RoleID, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrAccessDenied):
			utils.SendErrorResponse(c, http.StatusForbidden, "FORBIDDEN", "User has no access")
		case errors.Is(err, domain.ErrRoleNotFound):
			utils.SendErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "Role not found")
		case errors.Is(err, domain.ErrRoleInUse):
			utils.SendErrorResponse(c, http.StatusConflict, "CONFILCT", "Role is in use")
		default:
			utils.SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete role")
		}
		return
	}

	c.Status(http.StatusNoContent)
}

type updateRoleInput struct {
	RoleName string `json:"role_name"`
}

// TODO: swagger docs
func (roleHandler *RoleHandler) UpdateRole(c *gin.Context) {
	roleIDStr := c.Param("role_id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 64)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "INVALID_APPROVAL_ID", "Invalid approval ID")
		return
	}

	var req updateRoleInput

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	userID, err := utils.ExtractUserID(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	if req.RoleName == "" {
		utils.SendErrorResponse(c, http.StatusBadRequest, "MISSING_FIELDS", "Role name is required")
		return
	}

	err = roleHandler.usecase.UpdateRole(c.Request.Context(), uint(roleID), req.RoleName, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRoleAlreadyExists):
			utils.SendErrorResponse(c, http.StatusConflict, "ROLE_ALREADY_EXISTS", "Role with this name already exists")
		default:
			utils.SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to register role")
		}
		return
	}

	c.Status(http.StatusCreated)
}
