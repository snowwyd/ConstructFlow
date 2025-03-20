package http

import (
	"errors"
	"net/http"

	"backend/pkg/config"
	"backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Пропускаем OPTIONS-запросы
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// Извлечение токена из куки
		tokenString, err := c.Cookie("auth_token")
		if err != nil || tokenString == "" {
			utils.SendErrorResponse(c, http.StatusUnauthorized, "MISSING_TOKEN", "Authentication token is missing")
			c.Abort()
			return
		}

		// Проверка JWT
		claims, err := utils.ParseJWT(tokenString, cfg.AppSecret)
		if err != nil {
			switch {
			case errors.Is(err, jwt.ErrTokenExpired):
				utils.SendErrorResponse(c, http.StatusUnauthorized, "TOKEN_EXPIRED", "Token has expired")
			case errors.Is(err, jwt.ErrTokenMalformed):
				utils.SendErrorResponse(c, http.StatusUnauthorized, "MALFORMED_TOKEN", "Token is malformed")
			default:
				utils.SendErrorResponse(c, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid token")
			}
			c.Abort()
			return
		}

		// Извлечение userID из claims
		userID, ok := claims["uid"].(float64)
		if !ok {
			utils.SendErrorResponse(c, http.StatusUnauthorized, "INVALID_TOKEN_CLAIMS", "Missing or invalid user ID in token claims")
			c.Abort()
			return
		}

		// Сохранение userID в контексте
		c.Set("userID", uint(userID))
		c.Next()
	}
}
