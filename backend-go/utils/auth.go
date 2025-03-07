package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/trofimovm/summvideo/models"
)

// JWT claims структура
type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

// ValidateTelegramAuth проверяет данные аутентификации от Telegram
func ValidateTelegramAuth(authData *models.TelegramAuthData) error {
	// Проверяем, что auth_date не слишком старый (не более 24 часов)
	authTime := time.Unix(authData.AuthDate, 0)
	if time.Since(authTime) > 24*time.Hour {
		return errors.New("auth data is expired")
	}

	// Получаем ключ бота
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		return errors.New("TELEGRAM_BOT_TOKEN is not set")
	}

	// Создаем secret key из token
	secretKey := sha256.Sum256([]byte("WebAppData"))
	
	// Собираем все поля для проверки
	dataCheckString := createDataCheckString(authData)
	
	// Проверяем hash
	if !validateHash(dataCheckString, authData.Hash, secretKey[:]) {
		return errors.New("invalid hash")
	}

	return nil
}

// createDataCheckString создает строку для проверки hash
func createDataCheckString(authData *models.TelegramAuthData) string {
	// Создаем карту из полей
	fields := map[string]string{
		"auth_date": strconv.FormatInt(authData.AuthDate, 10),
		"first_name": authData.FirstName,
		"id": strconv.FormatInt(authData.ID, 10),
		"username": authData.Username,
	}
	
	if authData.LastName != "" {
		fields["last_name"] = authData.LastName
	}
	
	if authData.PhotoURL != "" {
		fields["photo_url"] = authData.PhotoURL
	}
	
	// Сортируем ключи
	var keys []string
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	
	// Собираем строку для проверки
	var dataCheckParts []string
	for _, k := range keys {
		if k != "hash" {
			dataCheckParts = append(dataCheckParts, fmt.Sprintf("%s=%s", k, fields[k]))
		}
	}
	
	return strings.Join(dataCheckParts, "\n")
}

// validateHash проверяет корректность hash
func validateHash(dataCheckString, hash string, secretKey []byte) bool {
	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(dataCheckString))
	expectedHash := hex.EncodeToString(h.Sum(nil))
	return expectedHash == hash
}

// GenerateJWT создает JWT токен для пользователя
func GenerateJWT(userID int64) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", errors.New("JWT_SECRET is not set")
	}

	// Устанавливаем срок действия токена (30 дней)
	expirationTime := time.Now().Add(30 * 24 * time.Hour)
	
	// Создаем claims
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "summvideo",
		},
	}
	
	// Создаем токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// Подписываем токен секретным ключом
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	
	return tokenString, nil
}

// ValidateJWT проверяет и возвращает claims из JWT токена
func ValidateJWT(tokenString string) (*Claims, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET is not set")
	}
	
	// Парсим токен
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	
	return claims, nil
}