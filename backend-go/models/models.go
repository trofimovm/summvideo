package models

import "time"

// VideoResponse представляет ответ API на запрос обработки видео
type VideoResponse struct {
	Summary       string `json:"summary"`
	Transcription string `json:"transcription"`
}

// ErrorResponse представляет ответ API при ошибке
type ErrorResponse struct {
	Error string `json:"error"`
}

// Config содержит основные настройки приложения
type Config struct {
	MaxFileSize int64  // Максимальный размер файла в байтах
	OpenAIKey   string // API ключ OpenAI
	LogDir      string // Директория для логов
}

// User представляет пользователя Telegram
type User struct {
	ID              int64     `json:"id"`
	TelegramID      int64     `json:"telegram_id"`
	Username        string    `json:"username"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	PhotoURL        string    `json:"photo_url"`
	AuthDate        int64     `json:"auth_date"`
	Hash            string    `json:"hash"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	LastLogin       time.Time `json:"last_login"`
	IsActive        bool      `json:"is_active"`
	IsAdmin         bool      `json:"is_admin"`         // Флаг админа
	UsageLimitSecs  int       `json:"usage_limit_secs"` // Лимит использования в секундах
	UsageTotalSecs  int       `json:"usage_total_secs"` // Общее использованное время в секундах
	PasswordHash    string    `json:"-"`                // Хеш пароля (для админов)
}

// AdminCredentials представляет данные для входа администратора
type AdminCredentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UsageLimitUpdate представляет обновление лимита использования для пользователя
type UsageLimitUpdate struct {
	UserID       int64 `json:"user_id" binding:"required"`
	LimitSeconds int   `json:"limit_seconds" binding:"required"`
}

// TelegramAuthData представляет данные авторизации из Telegram
type TelegramAuthData struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	PhotoURL  string `json:"photo_url"`
	AuthDate  int64  `json:"auth_date"`
	Hash      string `json:"hash"`
}

// AuthResponse представляет ответ на успешную авторизацию
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// UsageHistory представляет запись об использовании сервиса
type UsageHistory struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id"`
	VideoName      string    `json:"video_name"`
	PromptText     string    `json:"prompt_text"`
	SummaryLength  int       `json:"summary_length"`
	ProcessingTime int       `json:"processing_time"` // в секундах
	CreatedAt      time.Time `json:"created_at"`
}