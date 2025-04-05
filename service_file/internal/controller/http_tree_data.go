package http

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"service-file/internal/domain"
	"service-file/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UploadFile godoc
//
//	@Summary		Загрузить файл
//	@Description	Загружает файл в указанную директорию. Данные передаются в формате multipart/form-data: файл (file), ID директории (directory_id) и имя файла (name).
//	@Tags			file
//	@Security		ApiKeyAuth
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			file			formData	file					true	"Файл для загрузки"
//	@Param			directory_id	formData	int						true	"ID директории, куда загружается файл"
//	@Param			name			formData	string					true	"Имя файла"
//	@Success		201				{string}	string					"Файл успешно загружен"
//	@Failure		400				{object}	domain.ErrorResponse	"Нет файла или переданы некорректные данные"
//	@Failure		401				{object}	domain.ErrorResponse	"Пользователь не авторизован"
//	@Failure		403				{object}	domain.ErrorResponse	"Доступ к директории запрещен"
//	@Failure		404				{object}	domain.ErrorResponse	"Директория не найдена"
//	@Failure		409				{object}	domain.ErrorResponse	"Файл с таким именем уже существует"
//	@Failure		500				{object}	domain.ErrorResponse	"Ошибка при загрузке файла"
//	@Router			/files/upload [post]
func (h *TreeHandler) UploadFile(c *gin.Context) {
	userID, err := utils.ExtractUserID(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	form, _ := c.MultipartForm()
	files := form.File["file"]
	if len(files) == 0 {
		utils.SendErrorResponse(c, http.StatusBadRequest, "MISSING_FILE", "No file uploaded")
		return
	}

	fileHeader := files[0]
	file, err := fileHeader.Open()
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "FILE_READ_ERROR", err.Error())
		return
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "FILE_READ_ERROR", err.Error())
		return
	}

	directoryID, _ := strconv.Atoi(form.Value["directory_id"][0])
	name := form.Value["name"][0]

	// Определяем MIME-тип
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream" // Значение по умолчанию
	}

	err = h.fileUsecase.CreateFile(c.Request.Context(), uint(directoryID), name, fileData, contentType, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrDirectoryNotFound):
			utils.SendErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "Directory not found")
		case errors.Is(err, domain.ErrAccessDenied):
			utils.SendErrorResponse(c, http.StatusForbidden, "ACCESS_DENIED", "User has no access to this directory")
		case errors.Is(err, domain.ErrFileAlreadyExists):
			utils.SendErrorResponse(c, http.StatusConflict, "CONFLICT", "File already exists")
		default:
			utils.SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to upload file")
		}
		return
	}

	c.Status(http.StatusCreated)
}

// DownloadFileDirect godoc
//
//	@Summary		Скачать файл напрямую
//	@Description	Отдает файл для скачивания в виде потока. Если файл не найден – возвращает 404, а при отсутствии доступа – 403.
//	@Tags			file
//	@Security		ApiKeyAuth
//	@Param			file_id	path	int	true	"ID файла (например, 123)"
//	@Accept			json
//	@Produce		octet-stream
//	@Success		200	{file}		string					"Файл для скачивания"
//	@Failure		400	{object}	domain.ErrorResponse	"Невалидный ID файла"
//	@Failure		401	{object}	domain.ErrorResponse	"Пользователь не авторизован"
//	@Failure		403	{object}	domain.ErrorResponse	"Доступ к файлу запрещен"
//	@Failure		404	{object}	domain.ErrorResponse	"Файл не найден"
//	@Failure		500	{object}	domain.ErrorResponse	"Ошибка при получении файла"
//	@Router			/files/{file_id}/download-direct [get]
func (h *TreeHandler) DownloadFileDirect(c *gin.Context) {
	fileID, err := strconv.ParseUint(c.Param("file_id"), 10, 64)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid file ID")
		return
	}

	userID, err := utils.ExtractUserID(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	fileMeta, fileObject, err := h.fileUsecase.DownloadFileDirect(c.Request.Context(), uint(fileID), userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrFileNotFound):
			utils.SendErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "File not found")
		case errors.Is(err, domain.ErrAccessDenied):
			utils.SendErrorResponse(c, http.StatusForbidden, "ACCESS_DENIED", "No access to file")
		default:
			utils.SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve file")
		}
		return
	}
	defer fileObject.Close()

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileMeta.Name))
	c.Header("Content-Type", fileMeta.ContentType)
	c.Header("Content-Length", strconv.FormatInt(fileMeta.Size, 10))

	// Потоковая передача файла
	buf := make([]byte, 1024*1024) // 1MB буфер
	c.Stream(func(w io.Writer) bool {
		n, err := fileObject.Read(buf)
		if n > 0 {
			if _, writeErr := w.Write(buf[:n]); writeErr != nil {
				c.Error(writeErr)
				return false
			}
		}
		if err != nil {
			if err == io.EOF {
				return false
			}
			c.Error(err)
			return false
		}
		return true
	})
}

// UpdateFile godoc
//
//	@Summary		Обновить файл
//	@Description	Обновляет содержимое файла, переданное через multipart/form-data. Если файл не найден или нет доступа – возвращает соответствующую ошибку.
//	@Tags			file
//	@Security		ApiKeyAuth
//	@Param			file_id	path		int		true	"ID файла (например, 123)"
//	@Param			file	formData	file	true	"Новый файл для загрузки"
//	@Accept			multipart/form-data
//	@Produce		json
//	@Success		201	{string}	string					"Файл успешно обновлен"
//	@Failure		400	{object}	domain.ErrorResponse	"Невалидный ID файла или данные формы"
//	@Failure		401	{object}	domain.ErrorResponse	"Пользователь не авторизован"
//	@Failure		403	{object}	domain.ErrorResponse	"Доступ к файлу запрещен"
//	@Failure		404	{object}	domain.ErrorResponse	"Файл не найден"
//	@Failure		500	{object}	domain.ErrorResponse	"Ошибка при обновлении файла"
//	@Router			/files/{file_id} [put]
func (h *TreeHandler) UpdateFile(c *gin.Context) {
	fileIDStr := c.Param("file_id")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 64)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid file ID")
		return
	}

	userID, err := utils.ExtractUserID(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid multipart form")
		return
	}

	files := form.File["file"]
	if len(files) == 0 {
		utils.SendErrorResponse(c, http.StatusBadRequest, "MISSING_FILE", "No file uploaded")
		return
	}
	file, err := files[0].Open()
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "FILE_READ_ERROR", err.Error())
		return
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "FILE_READ_ERROR", err.Error())
		return
	}

	err = h.fileUsecase.UpdateFile(c.Request.Context(), uint(fileID), fileData, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrFileNotFound):
			utils.SendErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "File not found")
		case errors.Is(err, domain.ErrAccessDenied):
			utils.SendErrorResponse(c, http.StatusForbidden, "ACCESS_DENIED", "No access to file")
		default:
			utils.SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update file")
		}
		return
	}

	c.Status(http.StatusCreated)
}

// ConvertSTPToGLTF godoc
//
//	@Summary		Конвертировать STP в GLTF
//	@Description	Конвертирует файл STP в формат GLTF и возвращает GLB-файл. Если конвертация завершается ошибкой – возвращает сообщение об ошибке.
//	@Tags			file
//	@Security		ApiKeyAuth
//	@Param			file_id	path	int	true	"ID исходного STP файла (например, 123)"
//	@Accept			json
//	@Produce		octet-stream
//	@Success		200	{file}		string					"GLB-файл для просмотра"
//	@Failure		400	{object}	domain.ErrorResponse	"Невалидный ID файла"
//	@Failure		401	{object}	domain.ErrorResponse	"Пользователь не авторизован"
//	@Failure		500	{object}	domain.ErrorResponse	"Ошибка при конвертации файла"
//	@Router			/files/{file_id}/convert/gltf [get]
func (h *TreeHandler) ConvertSTPToGLTF(c *gin.Context) {
	fileID, err := strconv.ParseUint(c.Param("file_id"), 10, 64)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid file ID")
		return
	}

	userID, err := utils.ExtractUserID(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	outputPath, err := h.fileUsecase.ConvertSTPToGLTF(c.Request.Context(), uint(fileID), userID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "CONVERSION_ERROR", err.Error())
		return
	}

	// Устанавливаем заголовки для GLB-файла
	c.Header("Content-Type", "model/gltf-binary")
	c.Header("Content-Disposition", "inline; filename=\"model.glb\"")

	// Отдаём файл
	c.File(outputPath)
}
