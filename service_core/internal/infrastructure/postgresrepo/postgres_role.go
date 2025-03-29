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

// CreateRole добавляет роль в БД
func (r *RoleRepository) CreateRole(ctx context.Context, roleName string) error {
	const op = "infrastructure.postgresrepo.role.CreateRole"

	var existingRole domain.Role
	result := r.db.WithContext(ctx).Where("role_name = ?", roleName).First(&existingRole)

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

	if err := r.db.WithContext(ctx).Create(&newRole).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return fmt.Errorf("%s: %w", op, err)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// GetRoleByID возвращает название роли по ID
func (r *RoleRepository) GetRoleByID(ctx context.Context, roleID uint) (string, error) {
	const op = "infrastructure.postgresrepo.role.GetRoleByID"

	var Role domain.Role
	result := r.db.WithContext(ctx).First(&Role, roleID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("%s: %w", op, domain.ErrRoleNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, result.Error)
	}

	return Role.RoleName, nil
}
