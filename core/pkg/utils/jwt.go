package utils

import (
	"service-core/internal/domain"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT генерирует JWT для конкретного пользователя, принимает AppSecret и время жизни токена, возвращает строку
func GenerateJWT(user domain.User, appSecret string, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["login"] = user.Login
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(appSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseJWT проверяет JWT-токен и возвращает его claims
func ParseJWT(tokenString string, appSecret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем, что используется алгоритм HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(appSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// Проверяем, что токен валиден
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Извлекаем claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	// Проверяем срок действия токена
	expirationTime, ok := claims["exp"].(float64)
	if !ok {
		return nil, errors.New("missing expiration time")
	}

	if time.Now().Unix() > int64(expirationTime) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

// VerifyJWT принимает токен и возвращает uid из него
func VerifyJWT(tokenString string, appSecret string) (int, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(appSecret), nil
	})
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["uid"].(float64)
		if !ok {
			return 0, jwt.ErrInvalidKey
		}
		return int(userID), nil
	}

	return 0, jwt.ErrTokenInvalidClaims
}

// ExtractBearerToken извлекает JWT-токен из заголовка Authorization
func ExtractBearerToken(header string) string {
	if header == "" {
		return ""
	}

	// Разделяем строку по пробелу
	parts := strings.Split(header, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}
