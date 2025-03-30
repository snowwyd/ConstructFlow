package http

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"service-file/internal/domain"
	"service-file/internal/domain/interfaces"
	"service-file/pkg/utils"

	"github.com/gin-gonic/gin"
)

type TreeHandler struct {
	directoryUsecase interfaces.DirectoryUsecase
	fileUsecase      interfaces.FileUsecase
}

func NewTreeHandler(directoryUsecase interfaces.DirectoryUsecase, fileUsecase interfaces.FileUsecase) *TreeHandler {
	return &TreeHandler{
		directoryUsecase: directoryUsecase,
		fileUsecase:      fileUsecase,
	}
}

func (h *TreeHandler) GetTree(c *gin.Context) {
	var req struct {
		IsArchive bool `json:"is_archive"`
	}

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

func (h *TreeHandler) CreateDirectory(c *gin.Context) {
	userID, err := utils.ExtractUserID(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}
	var req struct {
		ParentPathID *uint  `json:"parent_path_id"`
		Name         string `json:"name"`
	}

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

func (h *TreeHandler) UploadFile(c *gin.Context) {
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

	directoryIDStr := form.Value["directory_id"][0]
	directoryID, _ := strconv.Atoi(directoryIDStr)
	name := form.Value["name"][0]

	err = h.fileUsecase.CreateFile(c.Request.Context(), uint(directoryID), name, fileData, userID)
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

func (h *TreeHandler) DeleteDirectory(c *gin.Context) {
	userID, err := utils.ExtractUserID(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	var req struct {
		DirectoryID uint `json:"directory_id"`
	}

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

func (h *TreeHandler) DeleteFile(c *gin.Context) {
	userID, err := utils.ExtractUserID(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	var req struct {
		FileID uint `json:"file_id"`
	}

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
