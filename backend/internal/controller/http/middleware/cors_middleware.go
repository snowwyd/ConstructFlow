package http

import "github.com/gin-gonic/gin"

// CORSMiddleware возвращает middleware для обработки CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Разрешаем все источники
		c.Header("Access-Control-Allow-Origin", "*")
		// Разрешенные методы
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		// Продолжаем обработку запроса
		c.Next()
	}
}
