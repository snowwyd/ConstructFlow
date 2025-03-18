package http

import (
	"net/http"
	"strconv"

	"backend/internal/domain/interfaces"
	"backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

type TreeHandler struct {
	usecase interfaces.FileTreeUsecase
}

func NewTreeHandler(usecase interfaces.FileTreeUsecase) *TreeHandler {
	return &TreeHandler{usecase: usecase}
}

func (h *TreeHandler) GetTree(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	userIDuint := userID.(uint)

	isArchiveStr := c.Query("is_archive")
	isArchive, err := strconv.ParseBool(isArchiveStr)
	if err != nil {
		isArchive = false
	}

	response, err := h.usecase.GetFileTree(c.Request.Context(), isArchive, userIDuint)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get file tree")
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *TreeHandler) UploadDirectory(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	userIDuint := userID.(uint)
	var req struct {
		ParentPathID *uint  `json:"parent_path_id"`
		Name         string `json:"name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	dirID, err := h.usecase.UploadDirectory(c.Request.Context(), req.ParentPathID, req.Name, userIDuint)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create directory")
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": dirID})
}

func (h *TreeHandler) UploadFile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	userIDuint := userID.(uint)
	var req struct {
		DirectoryID uint   `json:"directory_id"`
		Name        string `json:"name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	fileID, err := h.usecase.UploadFile(c.Request.Context(), req.DirectoryID, req.Name, userIDuint)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create file")
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": fileID})
}
