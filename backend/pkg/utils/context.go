package utils

import (
	"context"
	"errors"
)

// GetUserIDFromContext получвает userID(string) из контекста
func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return "", errors.New("failed to get user_id from context")
	}
	return userID, nil
}
