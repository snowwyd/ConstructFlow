package utils

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
)

// GetUserIDFromContext получвает userID(string) из контекста
func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return "", errors.New("failed to get user_id from context")
	}
	return userID, nil
}

func ExtractUserID(c *gin.Context) (uint, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, errors.New("user not authenticated")
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		return 0, errors.New("failed to parse user ID")
	}

	return userIDUint, nil
}
