package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// ExtractAudio извлекает аудио из видео файла
func ExtractAudio(videoPath string) (string, error) {
	// Создаем временный файл для аудио
	audioFile := filepath.Join(os.TempDir(), fmt.Sprintf("%s.wav", filepath.Base(videoPath)))

	// Используем ffmpeg для извлечения аудио
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-vn", "-acodec", "pcm_s16le", "-ar", "44100", "-ac", "2", audioFile)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ошибка извлечения аудио: %v", err)
	}

	return audioFile, nil
}

// ConvertToMP3 конвертирует аудио файл в MP3 формат
func ConvertToMP3(audioFile string) (string, error) {
	// Создаем имя для MP3 файла
	mp3File := filepath.Join(os.TempDir(), fmt.Sprintf("%s.mp3", filepath.Base(audioFile)))

	// Используем ffmpeg для конвертации в MP3
	cmd := exec.Command("ffmpeg", "-i", audioFile, "-codec:a", "libmp3lame", "-qscale:a", "2", "-b:a", "32k", mp3File)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ошибка конвертации аудио в MP3: %v", err)
	}

	return mp3File, nil
}