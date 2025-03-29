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
