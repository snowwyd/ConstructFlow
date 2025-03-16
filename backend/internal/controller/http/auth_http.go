package http

import (
	"errors"
	"net/http"

	"backend/internal/domain"
	"backend/internal/domain/interfaces"
	"backend/pkg/config"
	"backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	usecase interfaces.AuthUsecase
	cfg     *config.Config
}

// конструктор
func NewAuthHandler(usecase interfaces.AuthUsecase, cfg *config.Config) *AuthHandler {
	return &AuthHandler{usecase: usecase, cfg: cfg}
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

// GetCurrentUser - обработчик эндпоинта /auth/me
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	// получение userID из контекста
	userID, exists := c.Get("userID")
	if !exists {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	userIDint := userID.(int)

	// вызов Usecase GetCurrentUser
	userResponse, err := h.usecase.GetCurrentUser(c.Request.Context(), uint(userIDint))
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

// RegisterUser - обработчик эндпоинта /auth/register
func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// вызов Usecase RegisterUser
	userID, err := h.usecase.RegisterUser(c.Request.Context(), req.Login, req.Password, req.Role)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserAlreadyExists):
			utils.SendErrorResponse(c, http.StatusConflict, "USER_ALREADY_EXISTS", "User with this login already exists")
		default:
			utils.SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to register user")
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user_id": userID})
}
