package utils

import (
	"github.com/gin-gonic/gin"
)

// SendErrorResponse отправляет структурированный ответ об ошибке
func SendErrorResponse(c *gin.Context, statusCode int, code string, message string) {
	c.JSON(statusCode, gin.H{
		"error": gin.H{
			"code":    code,
			"message": message,
		},
	})
}
