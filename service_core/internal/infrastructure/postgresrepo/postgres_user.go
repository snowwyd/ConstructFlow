package postgresrepo

import (
	"context"
	"errors"
	"fmt"
	"service-core/internal/domain"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *Database) *UserRepository {
	return &UserRepository{db: db.db}
}

// SaveUser добавляет пользователя в БД
func (r *UserRepository) SaveUser(ctx context.Context, login string, passHash []byte, roleID uint) error {
	const op = "infrastructure.postgresrepo.user.SaveUser"

	var existingUser domain.User
	result := r.db.WithContext(ctx).Where("login = ?", login).First(&existingUser)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%s: %w", op, domain.ErrUserNotFound)
		}
	} else {
		return fmt.Errorf("%s: %w", op, domain.ErrUserAlreadyExists)
	}

	newUser := domain.User{
		Login:    login,
		PassHash: passHash,
		RoleID:   roleID,
	}

	if err := r.db.WithContext(ctx).Create(&newUser).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return fmt.Errorf("%s: %w", op, err)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// GetUserByID возвращает пользователя по ID
func (r *UserRepository) GetUserByID(ctx context.Context, userID uint) (domain.User, error) {
	const op = "infrastructure.postgresrepo.user.GetUserByID"

	var user domain.User
	result := r.db.WithContext(ctx).First(&user, userID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.User{}, fmt.Errorf("%s: %w", op, domain.ErrUserNotFound)
		}
		return domain.User{}, fmt.Errorf("%s: %w", op, result.Error)
	}

	return user, nil
}

// GetUserByLogin возвращает пользователя по логину
func (r *UserRepository) GetUserByLogin(ctx context.Context, login string) (domain.User, error) {
	const op = "infrastructure.postgresrepo.user.GetUserByLogin"

	var user domain.User
	result := r.db.WithContext(ctx).Where("login = ?", login).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.User{}, fmt.Errorf("%s: %w", op, domain.ErrUserNotFound)
		}
		return domain.User{}, fmt.Errorf("%s: %w", op, result.Error)
	}

	return user, nil
}

func (r *UserRepository) GetUserRole(ctx context.Context, userID uint) (string, error) {
	const op = "infrastructure.postgresrepo.user.GetUserRole"

	type Result struct {
		ID   uint
		Role string `gorm:"column:role_name"`
	}

	var result Result
	err := r.db.WithContext(ctx).
		Table("users").
		Select("users.id, roles.role_name").
		Joins("JOIN roles ON users.role_id = roles.id").
		Where("users.id = ?", userID).
		Scan(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return result.Role, nil
}

func (r *UserRepository) CheckUsersExist(ctx context.Context, userIDs []uint) (bool, error) {
	const op = "infrastructure.postgresrepo.user.CheckUsersExist"

	if len(userIDs) == 0 {
		return false, fmt.Errorf("%s: empty user list", op)
	}

	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("id IN (?)", userIDs).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return count == int64(len(userIDs)), nil
}

func (userRepo *UserRepository) CheckUsersWithRole(ctx context.Context, roleID uint) (bool, error) {
	const op = "infrastructure.postgresrepo.user.CheckUsersWithRole"

	var count int64
	err := userRepo.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("role_id = ?", roleID).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return count > 0, nil
}

func (userRepo *UserRepository) GetUsersGroupedByRoles(ctx context.Context) ([]domain.RoleData, error) {
	const op = "infrastructure.postgresrepo.user.GetUsersGroupedByRoles"

	type UserWithRole struct {
		UserID   uint   `json:"user_id"`
		Login    string `json:"login"`
		RoleName string `json:"role_name"`
	}

	var usersWithRoles []UserWithRole

	err := userRepo.db.WithContext(ctx).
		Table("users").
		Select("users.id AS user_id, users.login, roles.role_name").
		Joins("JOIN roles ON users.role_id = roles.id").
		Where("users.deleted_at IS NULL AND roles.deleted_at IS NULL").
		Scan(&usersWithRoles).Error

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	groupedData := make(map[string][]domain.UserData)
	for _, userWithRole := range usersWithRoles {
		groupedData[userWithRole.RoleName] = append(groupedData[userWithRole.RoleName], domain.UserData{
			UserID: userWithRole.UserID,
			Login:  userWithRole.Login,
		})
	}

	var roleData []domain.RoleData
	for roleName, users := range groupedData {
		roleData = append(roleData, domain.RoleData{
			RoleName: roleName,
			Users:    users,
		})
	}

	return roleData, nil
}
