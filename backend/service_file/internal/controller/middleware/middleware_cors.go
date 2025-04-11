package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware возвращает middleware для обработки CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Разрешаем конкретный источник (замените на ваш frontend URL)
		origin := c.Request.Header.Get("Origin")
		allowedOrigins := []string{"http://localhost:5173"} // Список разрешенных origin

		// Проверяем, разрешен ли origin
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Header("Access-Control-Allow-Credentials", "true") // Разрешаем credentials
		}

		// Обработка предварительного запроса (OPTIONS)
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		// Продолжаем обработку запроса
		c.Next()
	}
}
