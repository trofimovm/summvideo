package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/trofimovm/summvideo/models"
	"github.com/trofimovm/summvideo/repositories"
	"github.com/trofimovm/summvideo/services"
	"github.com/trofimovm/summvideo/utils"
)

// Константы приложения
const (
	MaxFileSize = 500 * 1024 * 1024 // 500 MB в байтах
)

// HomePage отображает главную страницу
func HomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "upload_video.html", gin.H{})
}

// RedirectToHome перенаправляет на главную страницу
func RedirectToHome(c *gin.Context) {
	c.Redirect(http.StatusFound, "/")
}

// UploadVideo обрабатывает загрузку и обработку видео
func UploadVideo(c *gin.Context) {
	// Проверка аутентификации
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "Требуется авторизация",
		})
		return
	}
	userID := userIDValue.(int64)

	// Проверяем, не превышен ли лимит использования
	userRepo := repositories.UserRepository{}
	remainingSeconds, err := userRepo.GetRemainingUsageSeconds(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Ошибка проверки лимита использования: " + err.Error(),
		})
		return
	}

	if remainingSeconds <= 0 {
		c.JSON(http.StatusForbidden, models.ErrorResponse{
			Error: "Превышен лимит бесплатного использования. Пожалуйста, свяжитесь с администратором для увеличения лимита.",
		})
		return
	}

	// Проверка OpenAI API ключа
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "API ключ OpenAI не найден",
		})
		return
	}

	// Получаем промт из формы
	prompt := c.PostForm("prompt")
	if prompt == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Промт не указан",
		})
		return
	}

	// Получаем файл видео
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Ошибка получения файла: " + err.Error(),
		})
		return
	}

	// Проверка размера файла
	if file.Size > MaxFileSize {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Размер файла превышает 500 МБ",
		})
		return
	}

	// Создаем временный файл для сохранения видео
	tempDir := os.TempDir()
	videoPath := filepath.Join(tempDir, file.Filename)

	// Сохраняем загруженный файл
	if err := c.SaveUploadedFile(file, videoPath); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Ошибка сохранения файла: " + err.Error(),
		})
		return
	}
	defer os.Remove(videoPath) // Удаляем временный файл после обработки

	// Фиксируем время начала обработки
	startTime := time.Now()

	// Обработка видео в отдельной горутине с использованием каналов для результатов
	resultCh := make(chan *models.VideoResponse)
	errorCh := make(chan error)

	go func() {
		// Извлечение аудио из видео
		audioFile, err := services.ExtractAudio(videoPath)
		if err != nil {
			errorCh <- err
			return
		}
		defer os.Remove(audioFile) // Удаляем временный аудиофайл

		// Конвертация аудио в mp3 для OpenAI API
		mp3File, err := services.ConvertToMP3(audioFile)
		if err != nil {
			errorCh <- err
			return
		}
		defer os.Remove(mp3File) // Удаляем временный mp3 файл

		// Транскрибация аудио через OpenAI API
		transcription, err := services.TranscribeAudio(apiKey, mp3File)
		if err != nil {
			errorCh <- err
			return
		}

		// Генерация саммари на основе транскрипции и промта
		summary, err := services.GenerateSummary(apiKey, transcription, prompt)
		if err != nil {
			errorCh <- err
			return
		}

		// Логирование результатов
		if err := utils.LogData(file.Filename, prompt, transcription, summary); err != nil {
			// Лог ошибки, но продолжаем обработку
			// Ошибка логирования не должна прервать обработку запроса
			// TODO: добавить логирование ошибки логирования
		}

		// Отправляем результат в канал
		resultCh <- &models.VideoResponse{
			Summary:       summary,
			Transcription: transcription,
		}
	}()

	// Ожидаем результат или ошибку
	select {
	case result := <-resultCh:
		// Вычисляем время обработки
		processingTime := int(time.Since(startTime).Seconds())
		summaryLength := len(result.Summary)

		// Обновляем использованное время пользователя
		if err := userRepo.UpdateUsage(userID, processingTime); err != nil {
			log.Printf("Ошибка обновления использованного времени: %v", err)
		}

		// Сохраняем запись об использовании
		usageRepo := repositories.UsageRepository{}
		_, err := usageRepo.Create(userID, file.Filename, prompt, summaryLength, processingTime)
		if err != nil {
			log.Printf("Ошибка сохранения записи об использовании: %v", err)
		}

		// Возвращаем успешный результат
		c.JSON(http.StatusOK, gin.H{
			"summary":        result.Summary,
			"transcription":  result.Transcription,
			"processing_time": processingTime,
		})

	case err := <-errorCh:
		// Возвращаем ошибку
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Ошибка обработки видео: " + err.Error(),
		})
	}
}