package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"service-core/internal/domain"
	"service-core/internal/domain/interfaces"
	"service-core/pkg/logger/slogger"
)

type RoleUsecase struct {
	roleRepo interfaces.RoleRepository
	userRepo interfaces.UserRepository
	log      *slog.Logger
}

func NewRoleUsecase(roleRepo interfaces.RoleRepository, userRepo interfaces.UserRepository, log *slog.Logger) *RoleUsecase {
	return &RoleUsecase{
		roleRepo: roleRepo,
		userRepo: userRepo,
		log:      log,
	}
}

func (roleUsecase *RoleUsecase) RegisterRole(ctx context.Context, roleName string, userID uint) error {
	const op = "usecase.role.RegisterRole"

	log := roleUsecase.log.With(slog.String("op", op), slog.Any("role", roleName))
	log.Info("registering new role")

	log.Debug("checking if user is admin")
	if err := roleUsecase.checkAdmin(ctx, userID); err != nil {
		log.Error("failed admin check", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("inserting role into DB")
	err := roleUsecase.roleRepo.CreateRole(ctx, roleName)
	if err != nil {
		if errors.Is(err, domain.ErrRoleAlreadyExists) {
			log.Error("role already exists", slogger.Err(domain.ErrRoleAlreadyExists))
			return fmt.Errorf("%s: %w", op, domain.ErrRoleAlreadyExists)
		}
		log.Error("failed to save role", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("role registered successfully")
	return nil
}

func (roleUsecase *RoleUsecase) GetRoles(ctx context.Context, userID uint) ([]domain.RoleResponse, error) {
	const op = "usecase.role.GetRoles"

	log := roleUsecase.log.With(slog.String("op", op), slog.Any("user_id", userID))
	log.Info("getting roles")

	log.Debug("checking if user is admin")
	if err := roleUsecase.checkAdmin(ctx, userID); err != nil {
		log.Error("failed admin check", slogger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("getting roles from database")
	roles, err := roleUsecase.roleRepo.GetRoles(ctx)
	if err != nil {
		// TODO: custom errors?
		log.Error("failed to get roles", slogger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("roles got succesfully")
	return roles, nil
}

func (roleUsecase *RoleUsecase) GetRoleByID(ctx context.Context, roleID uint, userID uint) (string, error) {
	const op = "usecase.role.GetRoleByID"

	log := roleUsecase.log.With(slog.String("op", op), slog.Any("user_id", userID))
	log.Info("getting roles")

	log.Debug("checking if user is admin")
	if err := roleUsecase.checkAdmin(ctx, userID); err != nil {
		log.Error("failed admin check", slogger.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("getting role from database")
	role, err := roleUsecase.roleRepo.GetRoleByID(ctx, roleID)
	if err != nil {
		// TODO: custom errors?
		log.Error("failed to get role", slogger.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("role got succesfully")
	return role, nil
}

func (roleUsecase *RoleUsecase) UpdateRole(ctx context.Context, roleID uint, roleName string, userID uint) error {
	const op = "usecase.role.UpdateRole"

	log := roleUsecase.log.With(slog.String("op", op))
	log.Info("updating role")

	log.Debug("checking if user is admin")
	if err := roleUsecase.checkAdmin(ctx, userID); err != nil {
		log.Error("failed admin check", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("checking role access")
	if err := roleUsecase.checkRole(ctx, roleID); err != nil {
		log.Error("failed role access check", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("checking is role name exists")
	if err := roleUsecase.checkRoleByName(ctx, roleName); err != nil {
		log.Error("role already exists", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("updating role")
	if err := roleUsecase.roleRepo.UpdateRole(ctx, roleID, roleName); err != nil {
		// TODO: custom errors?
		log.Error("failed to update role", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("role updated successfully")
	return nil
}

func (roleUsecase *RoleUsecase) DeleteRole(ctx context.Context, roleID uint, userID uint) error {
	const op = "usecase.role.DeleteRole"

	log := roleUsecase.log.With(slog.String("op", op))
	log.Info("deleting role")

	log.Debug("checking if user is admin")
	if err := roleUsecase.checkAdmin(ctx, userID); err != nil {
		log.Error("failed admin check", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("checking role access")
	if err := roleUsecase.checkRole(ctx, roleID); err != nil {
		log.Error("failed role access check", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("checking if role in use")
	hasUsers, err := roleUsecase.userRepo.CheckUsersWithRole(ctx, roleID)
	if err != nil {
		log.Error("failed to check users with role", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	if hasUsers {
		log.Warn("role is in use by users")
		return fmt.Errorf("%s: %w", op, domain.ErrRoleInUse)
	}

	log.Debug("deleting role")
	if err := roleUsecase.roleRepo.DeleteRole(ctx, roleID); err != nil {
		// TODO: custom errors
		log.Error("failed to delete role", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (roleUsecase *RoleUsecase) checkAdmin(ctx context.Context, userID uint) error {
	role, err := roleUsecase.userRepo.GetUserRole(ctx, userID)
	if err != nil {
		return err
	}
	if role != "admin" {
		return domain.ErrAccessDenied
	}
	return nil
}

func (roleUsecase *RoleUsecase) checkRole(ctx context.Context, roleID uint) error {
	roleName, err := roleUsecase.roleRepo.GetRoleByID(ctx, roleID)
	switch {
	case err != nil:
		return err
	case roleName == "":
		return domain.ErrRoleNotFound
	case roleName == "admin":
		return domain.ErrAccessDenied
	}

	return nil
}

func (roleUsecase *RoleUsecase) checkRoleByName(ctx context.Context, roleName string) error {
	exists, err := roleUsecase.roleRepo.CheckRoleByName(ctx, roleName)
	if err != nil {
		return err
	}
	if exists {
		return domain.ErrRoleAlreadyExists
	}
	return nil
}
