package interfaces

import (
	"backend/internal/domain"
	"context"
)

type AuthUsecase interface {
	Login(ctx context.Context, login, password string) (token string, err error)
	GetCurrentUser(ctx context.Context, userID uint) (userInfo domain.GetCurrentUserResponse, err error)

	// для админа
	RegisterUser(ctx context.Context, login, password, role string) (userID uint, err error)
}
