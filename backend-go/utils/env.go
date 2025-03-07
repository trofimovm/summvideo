package utils

import "os"

// GetEnv возвращает значение переменной окружения или значение по умолчанию,
// если переменная не найдена
func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}