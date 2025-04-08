package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"service-core/internal/domain"
	"service-core/internal/domain/interfaces"
	"service-core/pkg/logger/slogger"
	"service-core/pkg/utils"
)

type UserUsecase struct {
	userRepo interfaces.UserRepository
	roleRepo interfaces.RoleRepository
	log      *slog.Logger
}

func NewUserUsecase(userRepo interfaces.UserRepository, roleRepo interfaces.RoleRepository, log *slog.Logger) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
		roleRepo: roleRepo,
		log:      log,
	}
}

func (userUsecase *UserUsecase) GetUsersGrouped(ctx context.Context, userID uint) ([]domain.RoleData, error) {
	const op = "usecase.user.GetUsersGrouped"

	log := userUsecase.log.With(slog.String("op", op), slog.Any("user_id", userID))
	log.Info("getting users")

	log.Debug("checking if user is admin")
	if err := userUsecase.checkAdmin(ctx, userID); err != nil {
		log.Error("failed admin check", slogger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("getting users from database")
	roleData, err := userUsecase.userRepo.GetUsersGroupedByRoles(ctx)
	if err != nil {
		// TODO: custom errors?
		log.Error("failed to get users", slogger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("users got succesfully")
	return roleData, nil
}

// RegisterUser проверяет, есть ли такой пользователь, если нет - регистрирует
func (userUsecase *UserUsecase) RegisterUser(ctx context.Context, login, password string, roleID uint, userID uint) error {
	const op = "usecase.auth.RegisterUser"

	log := userUsecase.log.With(slog.String("op", op), slog.String("login", login))
	log.Info("registering user")

	log.Debug("checking if user is admin")
	if err := userUsecase.checkAdmin(ctx, userID); err != nil {
		log.Error("failed admin check", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("hashing password")
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Error("failed to hash password", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("getting role by ID")
	_, err = userUsecase.roleRepo.GetRoleByID(ctx, roleID)
	if err != nil {
		if errors.Is(err, domain.ErrRoleNotFound) {
			log.Warn("role not found", slogger.Err(err))
			return fmt.Errorf("%s: %w", op, domain.ErrRoleNotFound)
		}
		log.Error("failed to get role", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("saving user")
	err = userUsecase.userRepo.SaveUser(ctx, login, hashedPassword, roleID)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			log.Error("user already exists", slogger.Err(domain.ErrUserAlreadyExists))
			return fmt.Errorf("%s: %w", op, domain.ErrUserAlreadyExists)
		}
		log.Error("failed to save user", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user registered successfully")
	return nil
}

func (userUsecase *UserUsecase) checkAdmin(ctx context.Context, userID uint) error {
	role, err := userUsecase.userRepo.GetUserRole(ctx, userID)
	if err != nil {
		return err
	}
	if role != "admin" {
		return domain.ErrAccessDenied
	}
	return nil
}
