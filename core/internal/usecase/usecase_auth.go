package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"service-core/internal/domain"
	"service-core/internal/domain/interfaces"
	"service-core/pkg/config"
	"service-core/pkg/logger/slogger"
	"service-core/pkg/utils"
)

type AuthUsecase struct {
	userRepo interfaces.UserRepository
	roleRepo interfaces.RoleRepository
	cfg      *config.Config
	log      *slog.Logger
}

func NewAuthUsecase(userRepo interfaces.UserRepository, roleRepo interfaces.RoleRepository, cfg *config.Config, log *slog.Logger) *AuthUsecase {
	return &AuthUsecase{
		userRepo: userRepo,
		roleRepo: roleRepo,
		cfg:      cfg,
		log:      log,
	}
}

// Login - проверяет, существует ли пользователь, если есть - то логинит
func (u *AuthUsecase) Login(ctx context.Context, login, password string) (string, error) {
	const op = "usecase.auth.Login"

	log := u.log.With(slog.String("op", op), slog.String("login", login))
	log.Info("logging user in")

	log.Debug("getting user by login")
	user, err := u.userRepo.GetUserByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			u.log.Warn("user not found", slogger.Err(err))
			return "", fmt.Errorf("%s: %w", op, domain.ErrUserNotFound)
		}

		u.log.Error("failed to get user", slogger.Err(err))
		return "", fmt.Errorf("%s: %w", op, domain.ErrInvalidCredentials)
	}

	log.Debug("comparing pass and hash")
	if err := utils.CheckPassword(user.PassHash, password); err != nil {
		u.log.Error("invalid credentials", slogger.Err(err))
		return "", fmt.Errorf("%s: %w", op, domain.ErrInvalidCredentials)
	}

	log.Debug("generating JWT")
	token, err := utils.GenerateJWT(user, u.cfg.SecretKeys.KeyJwt, u.cfg.SecretKeys.TokenTTL)
	if err != nil {
		u.log.Error("failed to generate jwt", slogger.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged in successfully")
	return token, nil
}

// GetCurrentUser возвращает информацию о пользователе по ID
func (u *AuthUsecase) GetCurrentUser(ctx context.Context, userID uint) (domain.GetCurrentUserResponse, error) {
	const op = "usecase.auth.GetCurrentUser"

	log := u.log.With(slog.String("op", op), slog.Any("userID", userID))
	log.Info("getting current user info")

	log.Debug("getting user by ID")
	user, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			u.log.Warn("user not found", slogger.Err(err))
			return domain.GetCurrentUserResponse{}, fmt.Errorf("%s: %w", op, domain.ErrUserNotFound)
		}
		return domain.GetCurrentUserResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("getting role by ID")
	roleName, err := u.roleRepo.GetRoleByID(ctx, user.RoleID)
	if err != nil {
		if errors.Is(err, domain.ErrRoleNotFound) {
			u.log.Warn("role not found", slogger.Err(err))
			return domain.GetCurrentUserResponse{}, fmt.Errorf("%s: %w", op, domain.ErrRoleNotFound)
		}
		return domain.GetCurrentUserResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user info got successfully")
	return domain.GetCurrentUserResponse{
		ID:    user.ID,
		Login: user.Login,
		Role:  roleName,
	}, nil
}
