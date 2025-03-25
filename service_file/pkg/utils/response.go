package utils

import (
	"github.com/gin-gonic/gin"
)

// SendErrorResponse отправляет структурированный ответ об ошибке
func SendErrorResponse(c *gin.Context, statusCode int, code string, message string) {
	// Получаем заголовок Origin из запроса
	origin := c.Request.Header.Get("Origin")
	allowedOrigins := []string{"http://localhost:5173"} // Список разрешенных источников

	// Проверяем, разрешён ли origin
	allowed := false
	for _, allowedOrigin := range allowedOrigins {
		if origin == allowedOrigin {
			allowed = true
			break
		}
	}

	// Если источник разрешён, устанавливаем CORS-заголовки
	if allowed {
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")
	}

	c.JSON(statusCode, gin.H{
		"error": gin.H{
			"code":    code,
			"message": message,
		},
	})
}
