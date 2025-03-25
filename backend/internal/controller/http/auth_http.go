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

// Login godoc
// @Summary Аутентификация пользователя
// @Description Возвращает JWT токен при успешной аутентификации
// @Tags auth
// @Param login body string true "Логин для входа"
// @Param password body string true "Пароль для входа"
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Токен доступа"
// @Failure 400 {object} domain.ErrorResponse "Неверный запрос"
// @Failure 401 {object} domain.ErrorResponse "Неверные учетные данные"
// @Failure 404 {object} domain.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} domain.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/auth/login [post]
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

	// Установка HTTP-only куки
	c.SetCookie(
		"auth_token", // Имя куки
		token,        // Значение куки (токен)
		3600,         // Время жизни куки в секундах (1 час)
		"/",          // Путь, для которого куки доступен
		"",           // Домен (пустая строка = текущий домен)
		true,         // Secure: true (куки отправляется только по HTTPS)
		true,         // HttpOnly: true (куки недоступен через JavaScript)
	)

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// RegisterUser godoc
// @Summary Регистрация нового пользователя
// @Description Регистрирует нового пользователя и возвращает его ID
// @Tags auth
// @Accept json
// @Produce json
// @Param login body string true "Логин пользователя"
// @Param password body string true "Пароль пользователя"
// @Param role_id body uint true "ID роли, назначенной пользователю"
// @Success 201 {object} map[string]uint "ID созданного пользователя"
// @Failure 400 {object} domain.ErrorResponse "Неверный запрос"
// @Failure 409 {object} domain.ErrorResponse "Пользователь с таким логином уже существует"
// @Failure 404 {object} domain.ErrorResponse "Роль не найдена"
// @Failure 500 {object} domain.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/auth/register [post]
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
	err := h.usecase.RegisterUser(c.Request.Context(), req.Login, req.Password, req.RoleID)
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

	c.Status(http.StatusCreated)
}

// GetCurrentUser godoc
// @Summary Получение информации о текущем пользователе
// @Description Возвращает информацию о пользователе на основе JWT токена
// @Tags auth
// @Security JWT
// @Produce json
// @Success 200 {object} domain.GetCurrentUserResponse "Информация о пользователе"
// @Failure 401 {object} domain.ErrorResponse "Не авторизован"
// @Failure 404 {object} domain.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} domain.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, err := utils.ExtractUserID(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	// вызов Usecase GetCurrentUser
	userResponse, err := h.usecase.GetCurrentUser(c.Request.Context(), userID)
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

// RegisterRole godoc
// @Summary Регистрация новой роли
// @Description Регистрирует новую роль и возвращает её ID
// @Tags auth
// @Accept json
// @Produce json
// @Param role_name body string true "Название роли"
// @Success 201 {object} map[string]uint "ID созданной роли"
// @Failure 400 {object} domain.ErrorResponse "Неверный запрос"
// @Failure 409 {object} domain.ErrorResponse "Роль с таким названием уже существует"
// @Failure 500 {object} domain.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/auth/role [post]
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
	err := h.usecase.RegisterRole(c.Request.Context(), req.RoleName)
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
