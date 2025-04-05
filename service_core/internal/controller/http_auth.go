package http

import (
	"errors"
	"net/http"

	"service-core/internal/domain"
	"service-core/internal/domain/interfaces"
	"service-core/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	usecase interfaces.AuthUsecase
}

// конструктор
func NewAuthHandler(usecase interfaces.AuthUsecase) *AuthHandler {
	return &AuthHandler{usecase: usecase}
}

type loginInput struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Login godoc
// @Summary Аутентификация пользователя
// @Description Возвращает JWT токен при успешной аутентификации
// @Tags auth
// @Accept json
// @Produce json
// @Param input body loginInput true "Данные для входа"
// @Success 200 {object} map[string]string "Токен доступа"
// @Failure 400 {object} domain.ErrorResponse "Неверный запрос"
// @Failure 401 {object} domain.ErrorResponse "Неверные учетные данные"
// @Failure 404 {object} domain.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} domain.ErrorResponse "Внутренняя ошибка сервера"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginInput

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

type registerInput struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	RoleID   uint   `json:"role_id"`
}

// RegisterUser godoc
// @Summary Регистрация нового пользователя
// @Description Регистрирует нового пользователя на основе предоставленных данных (логин, пароль, ID роли) и возвращает HTTP статус 201 при успешной регистрации.
// @Tags auth
// @Accept json
// @Produce json
// @Param input body registerInput true "Данные для регистрации пользователя"
// @Success 201 {object} nil "Пользователь успешно зарегистрирован. Тело ответа пустое."
// @Failure 400 {object} domain.ErrorResponse "Неверный запрос: отсутствуют обязательные поля или некорректный формат данных."
// @Failure 409 {object} domain.ErrorResponse "Конфликт: пользователь с таким логином уже существует."
// @Failure 404 {object} domain.ErrorResponse "Роль с указанным ID не найдена."
// @Failure 500 {object} domain.ErrorResponse "Внутренняя ошибка сервера: не удалось зарегистрировать пользователя."
// @Router /auth/register [post]
func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var req registerInput

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
// @Description Возвращает информацию о пользователе на основе JWT токена, извлеченного из заголовка Authorization.
// @Tags auth
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} domain.GetCurrentUserResponse "Информация о пользователе"
// @Failure 401 {object} domain.ErrorResponse "Не авторизован: отсутствует или недействителен JWT токен."
// @Failure 404 {object} domain.ErrorResponse "Пользователь не найден: пользователь с указанным ID в токене не существует."
// @Failure 500 {object} domain.ErrorResponse "Внутренняя ошибка сервера: не удалось получить информацию о пользователе."
// @Router /auth/me [get]
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

type registerRoleInput struct {
	RoleName string `json:"role_name"`
}

// RegisterRole godoc
// @Summary Регистрация новой роли
// @Description Регистрирует новую роль на основе предоставленного названия и возвращает HTTP статус 201 при успешной регистрации.
// @Tags auth
// @Accept json
// @Produce json
// @Param input body registerRoleInput true "Название роли"
// @Success 201 {object} map[string]uint "ID созданной роли"
// @Failure 400 {object} domain.ErrorResponse "Неверный запрос: отсутствует название роли или некорректный формат данных."
// @Failure 409 {object} domain.ErrorResponse "Конфликт: роль с таким названием уже существует."
// @Failure 500 {object} domain.ErrorResponse "Внутренняя ошибка сервера: не удалось зарегистрировать роль."
// @Router /auth/role [post]
func (h *AuthHandler) RegisterRole(c *gin.Context) {
	var req registerRoleInput

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
