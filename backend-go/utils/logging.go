package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	logMutex sync.Mutex
)

// LogData записывает данные обработки видео в лог файл
func LogData(videoFilename, prompt, transcription, summary string) error {
	logMutex.Lock()
	defer logMutex.Unlock()

	logDir := GetEnv("LOG_DIR", "/var/log/summvideo")
	logFile := filepath.Join(logDir, "log.txt")

	// Создаем директорию для логов, если её нет
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		return fmt.Errorf("не удалось создать директорию для логов: %v", err)
	}

	// Открываем файл для записи (append)
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("не удалось открыть лог файл: %v", err)
	}
	defer f.Close()

	// Форматируем и записываем данные
	logData := fmt.Sprintf(
		"Дата и время: %s\nНазвание видео: %s\nПромт: %s\nТранскрибация: %s\nСаммари: %s\n%s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		videoFilename,
		prompt,
		transcription,
		summary,
		"--------------------------------------------------",
	)

	if _, err := f.WriteString(logData); err != nil {
		return fmt.Errorf("ошибка записи в лог: %v", err)
	}

	return nil
}