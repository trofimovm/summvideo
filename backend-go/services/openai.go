package services

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

// TranscribeAudio транскрибирует аудио файл через OpenAI API
func TranscribeAudio(apiKey, audioFile string) (string, error) {
	// Создаем клиента OpenAI
	client := openai.NewClient(apiKey)

	// Открываем аудио файл
	audioData, err := os.Open(audioFile)
	if err != nil {
		return "", fmt.Errorf("ошибка открытия аудио файла: %v", err)
	}
	defer audioData.Close()

	// Устанавливаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Отправляем запрос на транскрипцию
	resp, err := client.CreateTranscription(
		ctx,
		openai.AudioRequest{
			Model:    openai.Whisper1,
			FilePath: audioFile,
			Language: "ru",
		},
	)
	if err != nil {
		return "", fmt.Errorf("ошибка транскрипции аудио: %v", err)
	}

	return resp.Text, nil
}

// GenerateSummary генерирует саммари на основе транскрипции и промта
func GenerateSummary(apiKey, transcript, prompt string) (string, error) {
	// Создаем клиента OpenAI
	client := openai.NewClient(apiKey)

	// Устанавливаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	// Отправляем запрос на генерацию текста
	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: "gpt-4o-mini",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: prompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: transcript,
				},
			},
			Temperature: 0,
		},
	)
	if err != nil {
		return "", fmt.Errorf("ошибка генерации саммари: %v", err)
	}

	// Проверяем наличие ответа
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("пустой ответ от OpenAI API")
	}

	// Собираем результат
	summary := strings.TrimSpace(resp.Choices[0].Message.Content)
	return summary, nil
}