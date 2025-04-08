package postgresrepo

import (
	"context"
	"errors"
	"fmt"
	"service-core/internal/domain"

	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *Database) *RoleRepository {
	return &RoleRepository{db: db.db}
}

func (roleRepo *RoleRepository) GetRoles(ctx context.Context) ([]domain.RoleResponse, error) {
	const op = "infrastructure.postgresrepo.role.GetRoles"

	var roles []domain.RoleResponse
	err := roleRepo.db.WithContext(ctx).
		Table("roles").
		Select("id AS role_id, role_name").
		Where("deleted_at IS NULL").
		Order("role_name ASC").
		Scan(&roles).Error

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return roles, nil
}

func (roleRepo *RoleRepository) GetRoleByID(ctx context.Context, roleID uint) (string, error) {
	const op = "infrastructure.postgresrepo.role.GetRoleByID"

	var Role domain.Role
	result := roleRepo.db.WithContext(ctx).First(&Role, roleID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("%s: %w", op, domain.ErrRoleNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, result.Error)
	}

	return Role.RoleName, nil
}

func (roleRepo *RoleRepository) CreateRole(ctx context.Context, roleName string) error {
	const op = "infrastructure.postgresrepo.role.CreateRole"

	var existingRole domain.Role
	result := roleRepo.db.WithContext(ctx).Where("role_name = ?", roleName).First(&existingRole)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%s: %w", op, domain.ErrRoleNotFound)
		}
	} else {
		return fmt.Errorf("%s: %w", op, domain.ErrRoleAlreadyExists)
	}

	newRole := domain.Role{
		RoleName: roleName,
	}

	if err := roleRepo.db.WithContext(ctx).Create(&newRole).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return fmt.Errorf("%s: %w", op, err)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (roleRepo *RoleRepository) UpdateRole(ctx context.Context, roleID uint, roleName string) error {
	const op = "infrastructure.postgresrepo.role.UpdateRole"

	tx := roleRepo.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Обновление имени роли
	if err := tx.Model(&domain.Role{}).
		Where("id = ?", roleID).
		Update("role_name", roleName).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (roleRepo *RoleRepository) DeleteRole(ctx context.Context, roleID uint) error {
	const op = "infrastructure.postgresrepo.role.DeleteRole"

	tx := roleRepo.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	result := tx.Where("id = ?", roleID).Delete(&domain.Role{})
	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, result.Error)
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, domain.ErrRoleNotFound)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (roleRepo *RoleRepository) CheckRole(ctx context.Context, roleID uint) (bool, error) {
	const op = "infrastructure.postgresrepo.role.CheckRole"

	var count int64
	err := roleRepo.db.WithContext(ctx).
		Model(&domain.Role{}).
		Where("id = ? AND deleted_at IS NULL", roleID).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return count > 0, nil
}

func (roleRepo *RoleRepository) CheckRoleByName(ctx context.Context, roleName string) (bool, error) {
	const op = "infrastructure.postgresrepo.role.CheckRoleByName"

	var count int64
	err := roleRepo.db.WithContext(ctx).
		Model(&domain.Role{}).
		Where("role_name = ? AND deleted_at IS NULL", roleName).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return count > 0, nil
}
