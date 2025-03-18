package http

import (
	"errors"
	"net/http"

	"backend/internal/domain"
	"backend/internal/domain/interfaces"
	"backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	usecase interfaces.AuthUsecase
}

// конструктор
func NewAuthHandler(usecase interfaces.AuthUsecase) *AuthHandler {
	return &AuthHandler{usecase: usecase}
}

// Login - обработчик эндпоинта /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if req.Login == "" || req.Password == "" {
		utils.SendErrorResponse(c, http.StatusBadRequest, "MISSING_FIELDS", "Login and password are required")
		return
	}

	// вызов Usecase Login
	token, err := h.usecase.Login(c.Request.Context(), req.Login, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidCredentials):
			utils.SendErrorResponse(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid login or password")
		case errors.Is(err, domain.ErrUserNotFound):
			utils.SendErrorResponse(c, http.StatusNotFound, "USER_NOT_FOUND", "User not found")
		default:
			utils.SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// RegisterUser - обработчик эндпоинта /auth/register
func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
		RoleID   uint   `json:"role_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if req.Login == "" || req.Password == "" || req.RoleID == 0 {
		utils.SendErrorResponse(c, http.StatusBadRequest, "MISSING_FIELDS", "Login, password, and role_id are required")
		return
	}

	// вызов Usecase RegisterUser
	userID, err := h.usecase.RegisterUser(c.Request.Context(), req.Login, req.Password, req.RoleID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserAlreadyExists):
			utils.SendErrorResponse(c, http.StatusConflict, "USER_ALREADY_EXISTS", "User with this login already exists")
		case errors.Is(err, domain.ErrRoleNotFound):
			utils.SendErrorResponse(c, http.StatusNotFound, "ROLE_NOT_FOUND", "Role not found")
		default:
			utils.SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to register user")
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user_id": userID})
}

// GetCurrentUser - обработчик эндпоинта /auth/me
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	// получение userID из контекста
	userID, exists := c.Get("userID")
	if !exists {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	userIDuint := userID.(uint)

	// вызов Usecase GetCurrentUser
	userResponse, err := h.usecase.GetCurrentUser(c.Request.Context(), userIDuint)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			utils.SendErrorResponse(c, http.StatusNotFound, "USER_NOT_FOUND", "User not found")
		default:
			utils.SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		}
		return
	}

	c.JSON(http.StatusOK, userResponse)
}

// RegisterRole - обработчик эндпоинта /auth/role
func (h *AuthHandler) RegisterRole(c *gin.Context) {
	var req struct {
		RoleName string `json:"role_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if req.RoleName == "" {
		utils.SendErrorResponse(c, http.StatusBadRequest, "MISSING_FIELDS", "Role name is required")
		return
	}

	// вызов Usecase RegisterRole
	roleID, err := h.usecase.RegisterRole(c.Request.Context(), req.RoleName)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRoleAlreadyExists):
			utils.SendErrorResponse(c, http.StatusConflict, "ROLE_ALREADY_EXISTS", "Role with this name already exists")
		default:
			utils.SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to register role")
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"role_id": roleID})
}
