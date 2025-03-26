package http

import (
	"net/http"

	"service-file/internal/domain/interfaces"
	"service-file/pkg/utils"

	"github.com/gin-gonic/gin"
)

type TreeHandler struct {
	usecase interfaces.FileTreeUsecase
}

func NewTreeHandler(usecase interfaces.FileTreeUsecase) *TreeHandler {
	return &TreeHandler{usecase: usecase}
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

	response, err := h.usecase.GetFileTree(c.Request.Context(), req.IsArchive, userID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get file tree")
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *TreeHandler) GetFileInfo(c *gin.Context) {
	c.JSON(http.StatusOK, "unimplemented")
}

func (h *TreeHandler) CreateDirectory(c *gin.Context) {
	c.JSON(http.StatusOK, "unimplemented")
}

func (h *TreeHandler) UploadFile(c *gin.Context) {
	c.JSON(http.StatusOK, "unimplemented")
}

func (h *TreeHandler) DeleteDirectory(c *gin.Context) {
	c.JSON(http.StatusOK, "unimplemented")
}

func (h *TreeHandler) DeleteFile(c *gin.Context) {
	c.JSON(http.StatusOK, "unimplemented")
}
