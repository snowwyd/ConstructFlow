package http

import (
	"errors"
	"net/http"
	"service-core/internal/domain"
	"service-core/internal/domain/interfaces"
	"service-core/pkg/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	usecase interfaces.UserUsecase
}

func NewUserHandler(usecase interfaces.UserUsecase) *UserHandler {
	return &UserHandler{usecase: usecase}
}

func (userHandler *UserHandler) GetUsers(c *gin.Context) {
	userID, err := utils.ExtractUserID(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	users, err := userHandler.usecase.GetUsersGrouped(c.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrAccessDenied):
			utils.SendErrorResponse(c, http.StatusForbidden, "FORBIDDEN", "User has no access")
		default:
			utils.SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get users")
		}
		return
	}

	c.JSON(http.StatusOK, users)
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
// @Router /admin/users/register [post]
func (userHandler *UserHandler) RegisterUser(c *gin.Context) {
	var req registerInput

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if req.Login == "" || req.Password == "" || req.RoleID == 0 {
		utils.SendErrorResponse(c, http.StatusBadRequest, "MISSING_FIELDS", "Login, password, and role_id are required")
		return
	}

	userID, err := utils.ExtractUserID(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	// вызов Usecase RegisterUser
	err = userHandler.usecase.RegisterUser(c.Request.Context(), req.Login, req.Password, req.RoleID, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrAccessDenied):
			utils.SendErrorResponse(c, http.StatusForbidden, "FORBIDDEN", "User has no access")
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
