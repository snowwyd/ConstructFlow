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
	token, err := utils.GenerateJWT(user, u.cfg.AppSecret, u.cfg.TokenTTL)
	if err != nil {
		u.log.Error("failed to generate jwt", slogger.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged in successfully")
	return token, nil
}

// RegisterUser проверяет, есть ли такой пользователь, если нет - регистрирует
func (u *AuthUsecase) RegisterUser(ctx context.Context, login, password string, roleID uint) error {
	const op = "usecase.auth.RegisterUser"

	log := u.log.With(slog.String("op", op), slog.String("login", login))
	log.Info("registering user")

	log.Debug("hashing password")
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		u.log.Error("failed to hash password", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("getting role by ID")
	_, err = u.roleRepo.GetRoleByID(ctx, roleID)
	if err != nil {
		if errors.Is(err, domain.ErrRoleNotFound) {
			u.log.Warn("role not found", slogger.Err(err))
			return fmt.Errorf("%s: %w", op, domain.ErrRoleNotFound)
		}
		u.log.Error("failed to get role", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("saving user")
	err = u.userRepo.SaveUser(ctx, login, hashedPassword, roleID)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			u.log.Error("user already exists", slogger.Err(domain.ErrUserAlreadyExists))
			return fmt.Errorf("%s: %w", op, domain.ErrUserAlreadyExists)
		}
		u.log.Error("failed to save user", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user registered successfully")
	return nil
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

// RegisterRole создает роль
func (u *AuthUsecase) RegisterRole(ctx context.Context, roleName string) error {
	const op = "usecase.auth.RegisterRole"

	log := u.log.With(slog.String("op", op), slog.Any("role", roleName))
	log.Info("registering new role")

	log.Debug("inserting role into DB")
	err := u.roleRepo.CreateRole(ctx, roleName)
	if err != nil {
		if errors.Is(err, domain.ErrRoleAlreadyExists) {
			u.log.Error("role already exists", slogger.Err(domain.ErrRoleAlreadyExists))
			return fmt.Errorf("%s: %w", op, domain.ErrRoleAlreadyExists)
		}
		u.log.Error("failed to save role", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("role registered successfully")
	return nil
}
