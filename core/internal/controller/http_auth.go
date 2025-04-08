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