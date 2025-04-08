package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword хэширует пароль
func HashPassword(password string) ([]byte, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hashedBytes, nil
}

// CheckPassword проверяет пароль
func CheckPassword(hashedPassword []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
}
