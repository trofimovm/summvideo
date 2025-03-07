package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/trofimovm/summvideo/database"
	"github.com/trofimovm/summvideo/models"
	"golang.org/x/crypto/bcrypt"
)

// UserRepository предоставляет методы для работы с пользователями в БД
type UserRepository struct{}

// FindByTelegramID ищет пользователя по ID Telegram
func (r *UserRepository) FindByTelegramID(telegramID int64) (*models.User, error) {
	var user models.User

	err := database.DB.QueryRow(
		context.Background(),
		`SELECT id, telegram_id, username, first_name, last_name, photo_url, 
         auth_date, hash, created_at, updated_at, last_login, is_active,
         is_admin, usage_limit_secs, usage_total_secs, password_hash
         FROM users WHERE telegram_id = $1`,
		telegramID,
	).Scan(
		&user.ID, &user.TelegramID, &user.Username, &user.FirstName, &user.LastName,
		&user.PhotoURL, &user.AuthDate, &user.Hash, &user.CreatedAt, &user.UpdatedAt, 
		&user.LastLogin, &user.IsActive, &user.IsAdmin, &user.UsageLimitSecs, 
		&user.UsageTotalSecs, &user.PasswordHash,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByUsername ищет пользователя по имени пользователя
func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User

	err := database.DB.QueryRow(
		context.Background(),
		`SELECT id, telegram_id, username, first_name, last_name, photo_url, 
         auth_date, hash, created_at, updated_at, last_login, is_active,
         is_admin, usage_limit_secs, usage_total_secs, password_hash
         FROM users WHERE username = $1`,
		username,
	).Scan(
		&user.ID, &user.TelegramID, &user.Username, &user.FirstName, &user.LastName,
		&user.PhotoURL, &user.AuthDate, &user.Hash, &user.CreatedAt, &user.UpdatedAt, 
		&user.LastLogin, &user.IsActive, &user.IsAdmin, &user.UsageLimitSecs, 
		&user.UsageTotalSecs, &user.PasswordHash,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Create создает нового пользователя в БД
func (r *UserRepository) Create(user *models.TelegramAuthData) (*models.User, error) {
	var newUser models.User
	defaultUsageLimit := 3600 // 1 час в секундах

	err := database.DB.QueryRow(
		context.Background(),
		`INSERT INTO users (telegram_id, username, first_name, last_name, photo_url, auth_date, hash, usage_limit_secs) 
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
         RETURNING id, telegram_id, username, first_name, last_name, photo_url, 
         auth_date, hash, created_at, updated_at, last_login, is_active,
         is_admin, usage_limit_secs, usage_total_secs, password_hash`,
		user.ID, user.Username, user.FirstName, user.LastName, user.PhotoURL, user.AuthDate, user.Hash, defaultUsageLimit,
	).Scan(
		&newUser.ID, &newUser.TelegramID, &newUser.Username, &newUser.FirstName, &newUser.LastName,
		&newUser.PhotoURL, &newUser.AuthDate, &newUser.Hash, &newUser.CreatedAt, &newUser.UpdatedAt, 
		&newUser.LastLogin, &newUser.IsActive, &newUser.IsAdmin, &newUser.UsageLimitSecs, 
		&newUser.UsageTotalSecs, &newUser.PasswordHash,
	)

	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

// UpdateLoginTime обновляет время последнего входа пользователя
func (r *UserRepository) UpdateLoginTime(telegramID int64) error {
	now := time.Now()

	_, err := database.DB.Exec(
		context.Background(),
		`UPDATE users SET last_login = $1, updated_at = $1 WHERE telegram_id = $2`,
		now, telegramID,
	)

	return err
}

// FindByID ищет пользователя по ID в базе данных
func (r *UserRepository) FindByID(id int64) (*models.User, error) {
	var user models.User

	err := database.DB.QueryRow(
		context.Background(),
		`SELECT id, telegram_id, username, first_name, last_name, photo_url, 
         auth_date, hash, created_at, updated_at, last_login, is_active,
         is_admin, usage_limit_secs, usage_total_secs, password_hash
         FROM users WHERE id = $1`,
		id,
	).Scan(
		&user.ID, &user.TelegramID, &user.Username, &user.FirstName, &user.LastName,
		&user.PhotoURL, &user.AuthDate, &user.Hash, &user.CreatedAt, &user.UpdatedAt, 
		&user.LastLogin, &user.IsActive, &user.IsAdmin, &user.UsageLimitSecs, 
		&user.UsageTotalSecs, &user.PasswordHash,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Deactivate деактивирует пользователя
func (r *UserRepository) Deactivate(telegramID int64) error {
	_, err := database.DB.Exec(
		context.Background(),
		`UPDATE users SET is_active = false, updated_at = $1 WHERE telegram_id = $2`,
		time.Now(), telegramID,
	)

	return err
}

// CreateOrUpdate создает нового пользователя или обновляет существующего
func (r *UserRepository) CreateOrUpdate(authData *models.TelegramAuthData) (*models.User, error) {
	// Проверяем существует ли пользователь
	existingUser, err := r.FindByTelegramID(authData.ID)
	
	// Если пользователь найден, обновляем информацию
	if err == nil && existingUser != nil {
		_, err := database.DB.Exec(
			context.Background(),
			`UPDATE users 
             SET username = $1, first_name = $2, last_name = $3, photo_url = $4, 
             auth_date = $5, hash = $6, updated_at = $7, last_login = $7 
             WHERE telegram_id = $8`,
			authData.Username, authData.FirstName, authData.LastName, authData.PhotoURL,
			authData.AuthDate, authData.Hash, time.Now(), authData.ID,
		)
		if err != nil {
			return nil, err
		}

		// Получаем обновленного пользователя
		return r.FindByTelegramID(authData.ID)
	}

	// Если пользователь не найден, создаем нового
	return r.Create(authData)
}

// UpdateUsage обновляет использованное время пользователя
func (r *UserRepository) UpdateUsage(userID int64, additionalSeconds int) error {
	_, err := database.DB.Exec(
		context.Background(),
		`UPDATE users SET usage_total_secs = usage_total_secs + $1, updated_at = $2 WHERE id = $3`,
		additionalSeconds, time.Now(), userID,
	)

	return err
}

// UpdateUsageLimit обновляет лимит использования для пользователя
func (r *UserRepository) UpdateUsageLimit(userID int64, limitSeconds int) error {
	_, err := database.DB.Exec(
		context.Background(),
		`UPDATE users SET usage_limit_secs = $1, updated_at = $2 WHERE id = $3`,
		limitSeconds, time.Now(), userID,
	)

	return err
}

// GetAllUsers возвращает список всех пользователей (для админа)
func (r *UserRepository) GetAllUsers(limit, offset int) ([]models.User, error) {
	rows, err := database.DB.Query(
		context.Background(),
		`SELECT id, telegram_id, username, first_name, last_name, photo_url, 
         auth_date, hash, created_at, updated_at, last_login, is_active,
         is_admin, usage_limit_secs, usage_total_secs, password_hash
         FROM users 
         ORDER BY created_at DESC 
         LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID, &user.TelegramID, &user.Username, &user.FirstName, &user.LastName,
			&user.PhotoURL, &user.AuthDate, &user.Hash, &user.CreatedAt, &user.UpdatedAt,
			&user.LastLogin, &user.IsActive, &user.IsAdmin, &user.UsageLimitSecs,
			&user.UsageTotalSecs, &user.PasswordHash,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// GetRemainingUsageSeconds возвращает оставшееся время использования для пользователя в секундах
func (r *UserRepository) GetRemainingUsageSeconds(userID int64) (int, error) {
	var limit, used int
	
	err := database.DB.QueryRow(
		context.Background(),
		`SELECT usage_limit_secs, usage_total_secs FROM users WHERE id = $1`,
		userID,
	).Scan(&limit, &used)
	
	if err != nil {
		return 0, err
	}
	
	remaining := limit - used
	if remaining < 0 {
		remaining = 0
	}
	
	return remaining, nil
}

// CheckAdminCredentials проверяет учетные данные администратора
func (r *UserRepository) CheckAdminCredentials(username, password string) (*models.User, error) {
	user, err := r.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	
	if !user.IsAdmin {
		return nil, errors.New("пользователь не является администратором")
	}
	
	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("неверный пароль")
	}
	
	return user, nil
}

// GetTotalUsersCount возвращает общее количество пользователей
func (r *UserRepository) GetTotalUsersCount() (int, error) {
	var count int
	
	err := database.DB.QueryRow(
		context.Background(),
		`SELECT COUNT(*) FROM users`,
	).Scan(&count)
	
	return count, err
}