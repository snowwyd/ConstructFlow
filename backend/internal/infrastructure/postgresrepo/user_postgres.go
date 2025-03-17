package postgresrepo

import (
	"backend/internal/domain"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository конструктор
func NewUserRepository(db *Database) *UserRepository {
	return &UserRepository{db: db.db}
}

// SaveUser добавляет пользователя в БД
func (r *UserRepository) SaveUser(ctx context.Context, login string, passHash []byte, roleID uint) (uint, error) {
	const op = "postgresrepo.user.SaveUser"

	var existingUser domain.User
	result := r.db.WithContext(ctx).Where("login = ?", login).First(&existingUser)

	// обработка ошибок и отсутствия пользователя
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("%s: user not found: %w", op, domain.ErrUserNotFound)
		}
	} else {
		return 0, fmt.Errorf("%s: %w", op, domain.ErrUserAlreadyExists)
	}

	newUser := domain.User{
		Login:    login,
		PassHash: passHash,
		RoleID:   roleID,
	}

	// создает пользователя и парсит в модель User
	if err := r.db.WithContext(ctx).Create(&newUser).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return 0, fmt.Errorf("%s: duplicate key error: %w", op, err)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return newUser.ID, nil
}

// GetUserByID возвращает пользователя по ID
func (r *UserRepository) GetUserByID(ctx context.Context, userID uint) (domain.User, error) {
	const op = "postgresrepo.user.GetUserByID"

	var user domain.User
	result := r.db.WithContext(ctx).First(&user, userID)

	// обработка ошибок и отсутствия пользователя
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.User{}, fmt.Errorf("%s: user not found: %w", op, domain.ErrUserNotFound)
		}
		return domain.User{}, fmt.Errorf("%s: %w", op, result.Error)
	}

	return user, nil
}

// GetUserByLogin возвращает пользователя по логину
func (r *UserRepository) GetUserByLogin(ctx context.Context, login string) (domain.User, error) {
	const op = "postgresrepo.user.GetUserByLogin"

	var user domain.User
	result := r.db.WithContext(ctx).Where("login = ?", login).First(&user)

	// обработка ошибок и отсутствия пользователя
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.User{}, fmt.Errorf("%s: user not found: %w", op, domain.ErrUserNotFound)
		}
		return domain.User{}, fmt.Errorf("%s: %w", op, result.Error)
	}

	return user, nil
}

// TODO: implement?
func (r *UserRepository) UpdateUserRole(ctx context.Context, userID uint, role string) (success bool, err error) {
	panic("unimplemented")
}
