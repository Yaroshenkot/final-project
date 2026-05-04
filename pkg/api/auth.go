package api

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// secretKey возвращает ключ для подписи JWT на основе пароля
func secretKey() []byte {
	pass := os.Getenv("TODO_PASSWORD")
	hash := sha256.Sum256([]byte(pass))
	return []byte(fmt.Sprintf("%x", hash))
}

// GenerateToken создаёт JWT-токен
func GenerateToken() (string, error) {
	claims := jwt.MapClaims{
		"hash": fmt.Sprintf("%x", sha256.Sum256([]byte(os.Getenv("TODO_PASSWORD")))),
		"exp":  time.Now().Add(8 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey())
}

// ValidateToken проверяет JWT-токен
func ValidateToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}
		return secretKey(), nil
	})

	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expectedHash := fmt.Sprintf("%x", sha256.Sum256([]byte(os.Getenv("TODO_PASSWORD"))))
		if claims["hash"] == expectedHash {
			return true, nil
		}
	}

	return false, nil
}

// auth — middleware для проверки аутентификации
func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var jwtToken string
			cookie, err := r.Cookie("token")
			if err == nil {
				jwtToken = cookie.Value
			}

			valid, _ := ValidateToken(jwtToken)
			if !valid {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
