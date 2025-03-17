package postgresrepo

import (
	"backend/internal/domain"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

// NewRoleRepository конструктор
func NewRoleRepository(db *Database) *RoleRepository {
	return &RoleRepository{db: db.db}
}

// CreateRole добавляет роль в БД
func (r *RoleRepository) CreateRole(ctx context.Context, roleName string) (uint, error) {
	const op = "postgresrepo.role.CreateRole"

	var existingRole domain.Role
	result := r.db.WithContext(ctx).Where("role_name = ?", roleName).First(&existingRole)

	// обработка ошибок и отсутствия пользователя
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("%s: role not found: %w", op, domain.ErrRoleNotFound)
		}
	} else {
		return 0, fmt.Errorf("%s: %w", op, domain.ErrRoleAlreadyExists)
	}

	newRole := domain.Role{
		RoleName: roleName,
	}

	// создает пользователя и парсит в модель Role
	if err := r.db.WithContext(ctx).Create(&newRole).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return 0, fmt.Errorf("%s: duplicate key error: %w", op, err)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return newRole.ID, nil
}

// GetRoleByID возвращает название роли по ID
func (r *RoleRepository) GetRoleByID(ctx context.Context, roleID uint) (string, error) {
	const op = "postgresrepo.role.GetRoleByID"

	var Role domain.Role
	result := r.db.WithContext(ctx).First(&Role, roleID)

	// обработка ошибок и отсутствия пользователя
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("%s: Role not found: %w", op, domain.ErrRoleNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, result.Error)
	}

	return Role.RoleName, nil
}
