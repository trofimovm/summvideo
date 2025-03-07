package repositories

import (
	"context"

	"github.com/trofimovm/summvideo/database"
	"github.com/trofimovm/summvideo/models"
)

// UsageRepository предоставляет методы для работы с историей использования
type UsageRepository struct{}

// Create создает новую запись об использовании сервиса
func (r *UsageRepository) Create(userID int64, videoName, promptText string, summaryLength, processingTime int) (*models.UsageHistory, error) {
	var usage models.UsageHistory

	err := database.DB.QueryRow(
		context.Background(),
		`INSERT INTO usage_history (user_id, video_name, prompt_text, summary_length, processing_time) 
         VALUES ($1, $2, $3, $4, $5) 
         RETURNING id, user_id, video_name, prompt_text, summary_length, processing_time, created_at`,
		userID, videoName, promptText, summaryLength, processingTime,
	).Scan(
		&usage.ID, &usage.UserID, &usage.VideoName, &usage.PromptText,
		&usage.SummaryLength, &usage.ProcessingTime, &usage.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &usage, nil
}

// FindByUserID возвращает историю использования конкретным пользователем
func (r *UsageRepository) FindByUserID(userID int64, limit, offset int) ([]models.UsageHistory, error) {
	rows, err := database.DB.Query(
		context.Background(),
		`SELECT id, user_id, video_name, prompt_text, summary_length, processing_time, created_at 
         FROM usage_history 
         WHERE user_id = $1 
         ORDER BY created_at DESC 
         LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usageList []models.UsageHistory
	for rows.Next() {
		var usage models.UsageHistory
		if err := rows.Scan(
			&usage.ID, &usage.UserID, &usage.VideoName, &usage.PromptText,
			&usage.SummaryLength, &usage.ProcessingTime, &usage.CreatedAt,
		); err != nil {
			return nil, err
		}
		usageList = append(usageList, usage)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return usageList, nil
}

// GetTotalCount возвращает общее количество записей для пользователя
func (r *UsageRepository) GetTotalCount(userID int64) (int, error) {
	var count int
	err := database.DB.QueryRow(
		context.Background(),
		"SELECT COUNT(*) FROM usage_history WHERE user_id = $1",
		userID,
	).Scan(&count)
	
	return count, err
}

// GetRecentActivity возвращает недавнюю активность всех пользователей
func (r *UsageRepository) GetRecentActivity(limit int) ([]models.UsageHistory, error) {
	rows, err := database.DB.Query(
		context.Background(),
		`SELECT id, user_id, video_name, prompt_text, summary_length, processing_time, created_at 
         FROM usage_history 
         ORDER BY created_at DESC 
         LIMIT $1`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usageList []models.UsageHistory
	for rows.Next() {
		var usage models.UsageHistory
		if err := rows.Scan(
			&usage.ID, &usage.UserID, &usage.VideoName, &usage.PromptText,
			&usage.SummaryLength, &usage.ProcessingTime, &usage.CreatedAt,
		); err != nil {
			return nil, err
		}
		usageList = append(usageList, usage)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return usageList, nil
}