package http

import (
	"errors"
	"net/http"
	"strconv"

	"service-file/internal/domain"
	"service-file/pkg/utils"

	"github.com/gin-gonic/gin"
)

type getTreeInput struct {
	IsArchive bool `json:"is_archive"`
}

// GetTree godoc
//
//	@Summary		Получить файловое дерево
//	@Description	Отдает дерево файлов, доступных конкретному пользователю. Если флаг isArchive = true, отдает полностью дерево с статусом "archive".
//	@Tags			tree
//	@Security		ApiKeyAuth
//	@Param			input	body	getTreeInput	true	"Параметры для получения дерева (например, флаг isArchive)"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	domain.FileResponse		"Дерево файлов успешно отображено"
//	@Failure		400	{object}	domain.ErrorResponse	"Невалидный запрос"
//	@Failure		401	{object}	domain.ErrorResponse	"Пользователь не авторизован"
//	@Failure		403	{object}	domain.ErrorResponse	"Доступ к репозиторию запрещен"
//	@Failure		500	{object}	domain.ErrorResponse	"Ошибка при получении дерева файлов"
//	@Router			/directories [post]
func (h *TreeHandler) GetTree(c *gin.Context) {
	var req getTreeInput

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	userID, err := utils.ExtractUserID(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	response, err := h.directoryUsecase.GetFileTree(c.Request.Context(), req.IsArchive, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrAccessDenied):
			utils.SendErrorResponse(c, http.StatusForbidden, "ACCESS_DENIED", "user has no access to this repository")
		default:
			utils.SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get file tree")
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetFileInfo godoc
//
//	@Summary		Получить информацию о файле
//	@Description	Отдает подробную информацию о файле по его ID. Если файл не найден – возвращает 404, если пользователь не имеет доступа – 403.
//	@Tags			file
//	@Security		ApiKeyAuth
//	@Param			file_id	path	string	true	"ID файла (например, 123)"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	domain.FileResponse		"Информация о файле"
//	@Failure		400	{object}	domain.ErrorResponse	"Невалидный ID файла"
//	@Failure		401	{object}	domain.ErrorResponse	"Пользователь не авторизован"
//	@Failure		403	{object}	domain.ErrorResponse	"Доступ к файлу запрещен"
//	@Failure		404	{object}	domain.ErrorResponse	"Файл не найден"
//	@Failure		500	{object}	domain.ErrorResponse	"Ошибка при получении информации о файле"
//	@Router			/files/{file_id} [get]
func (h *TreeHandler) GetFileInfo(c *gin.Context) {
	userID, err := utils.ExtractUserID(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	fileIDStr := c.Param("file_id")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 64)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid file ID")
		return
	}

	fileInfo, err := h.fileUsecase.GetFileInfo(c.Request.Context(), uint(fileID), userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrAccessDenied):
			utils.SendErrorResponse(c, http.StatusForbidden, "ACCESS_DENIED", "You do not have access to this file")
		case errors.Is(err, domain.ErrFileNotFound):
			utils.SendErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "File not found")
		default:
			utils.SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get file info")
		}
		return
	}

	c.JSON(http.StatusOK, fileInfo)
}

type createDirectoryInput struct {
	ParentPathID *uint  `json:"parent_path_id"`
	Name         string `json:"name"`
}

// CreateDirectory godoc
//
//	@Summary		Создать директорию
//	@Description	Создает новую директорию в файловой системе для пользователя. Обязательны поля parent_path_id и name.
//	@Tags			directory
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			input	body		createDirectoryInput	true	"Параметры для создания директории"
//	@Success		201		{string}	string					"Директория успешно создана"
//	@Failure		400		{object}	domain.ErrorResponse	"Невалидный запрос или отсутствуют обязательные поля"
//	@Failure		401		{object}	domain.ErrorResponse	"Пользователь не авторизован"
//	@Failure		403		{object}	domain.ErrorResponse	"Доступ к директории запрещен"
//	@Failure		404		{object}	domain.ErrorResponse	"Родительская директория не найдена"
//	@Failure		409		{object}	domain.ErrorResponse	"Директория уже существует"
//	@Failure		500		{object}	domain.ErrorResponse	"Ошибка при создании директории"
//	@Router			/directories/create [post]
func (h *TreeHandler) CreateDirectory(c *gin.Context) {
	userID, err := utils.ExtractUserID(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}
	var req createDirectoryInput

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if req.Name == "" || req.ParentPathID == nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "MISSING_FIELDS", "Name and parent_path_id are required")
		return
	}

	err = h.directoryUsecase.CreateDirectory(c.Request.Context(), req.ParentPathID, req.Name, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrDirectoryNotFound):
			utils.SendErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "Directory not found")
		case errors.Is(err, domain.ErrAccessDenied):
			utils.SendErrorResponse(c, http.StatusForbidden, "ACCESS_DENIED", "User has no access to this directory")
		case errors.Is(err, domain.ErrDirectoryAlreadyExists):
			utils.SendErrorResponse(c, http.StatusConflict, "CONFLICT", "Directory already exists")
		default:
			utils.SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to upload directory")
		}
		return
	}

	c.Status(http.StatusCreated)
}

type deleteDirectoryInput struct {
	DirectoryID uint `json:"directory_id"`
}

// DeleteDirectory godoc
//
//	@Summary		Удалить директорию
//	@Description	Удаляет директорию по её ID, переданному в теле запроса. Если директория не найдена, содержит недопустимые файлы или доступ запрещен – возвращает ошибку.
//	@Tags			directory
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			input	body		deleteDirectoryInput	true	"Параметры для удаления директории"
//	@Success		204		{string}	string					"Директория успешно удалена"
//	@Failure		400		{object}	domain.ErrorResponse	"Невалидный запрос или отсутствуют обязательные поля"
//	@Failure		401		{object}	domain.ErrorResponse	"Пользователь не авторизован"
//	@Failure		403		{object}	domain.ErrorResponse	"Доступ к директории запрещен"
//	@Failure		404		{object}	domain.ErrorResponse	"Директория не найдена"
//	@Failure		409		{object}	domain.ErrorResponse	"Директория содержит файлы, не соответствующие черновикам"
//	@Failure		500		{object}	domain.ErrorResponse	"Ошибка при удалении директории"
//	@Router			/directories [delete]
func (h *TreeHandler) DeleteDirectory(c *gin.Context) {
	userID, err := utils.ExtractUserID(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	var req deleteDirectoryInput

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if req.DirectoryID == 0 {
		utils.SendErrorResponse(c, http.StatusBadRequest, "MISSING_FIELDS", "Directory_id is required")
		return
	}

	err = h.directoryUsecase.DeleteDirectory(c.Request.Context(), req.DirectoryID, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrDirectoryNotFound):
			utils.SendErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "Directory not found")
		case errors.Is(err, domain.ErrAccessDenied):
			utils.SendErrorResponse(c, http.StatusForbidden, "ACCESS_DENIED", "User has no access to delete this directory")
		case errors.Is(err, domain.ErrDirectoryContainsNonDraftFiles):
			utils.SendErrorResponse(c, http.StatusConflict, "CONFLICT", "Directory contains non-draft files and cannot be deleted")
		default:
			utils.SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete directory")
		}
		return
	}

	c.Status(http.StatusNoContent)
}

type deleteFileInput struct {
	FileID uint `json:"file_id"`
}

// DeleteFile godoc
//
//	@Summary		Удалить файл
//	@Description	Удаляет файл по его ID, переданному в теле запроса. Если файл не найден или доступ запрещен – возвращает ошибку.
//	@Tags			file
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			input	body		deleteFileInput			true	"Параметры для удаления файла"
//	@Success		204		{string}	string					"Файл успешно удален"
//	@Failure		400		{object}	domain.ErrorResponse	"Невалидный запрос или отсутствуют обязательные поля"
//	@Failure		401		{object}	domain.ErrorResponse	"Пользователь не авторизован"
//	@Failure		403		{object}	domain.ErrorResponse	"Доступ к файлу запрещен"
//	@Failure		404		{object}	domain.ErrorResponse	"Файл не найден"
//	@Failure		409		{object}	domain.ErrorResponse	"Нельзя удалить файл, не являющийся черновиком"
//	@Failure		500		{object}	domain.ErrorResponse	"Ошибка при удалении файла"
//	@Router			/files [delete]
func (h *TreeHandler) DeleteFile(c *gin.Context) {
	userID, err := utils.ExtractUserID(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	var req deleteFileInput

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if req.FileID == 0 {
		utils.SendErrorResponse(c, http.StatusBadRequest, "MISSING_FIELDS", "File_id is required")
		return
	}

	err = h.fileUsecase.DeleteFile(c.Request.Context(), req.FileID, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrFileNotFound):
			utils.SendErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "File not found")
		case errors.Is(err, domain.ErrAccessDenied):
			utils.SendErrorResponse(c, http.StatusForbidden, "ACCESS_DENIED", "User has no access to delete this file")
		case errors.Is(err, domain.ErrCannotDeleteNonDraftFile):
			utils.SendErrorResponse(c, http.StatusConflict, "CONFLICT", "Cannot delete non-draft files")
		default:
			utils.SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete file")
		}
		return
	}

	c.Status(http.StatusNoContent)
}
