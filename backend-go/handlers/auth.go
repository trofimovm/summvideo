package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/trofimovm/summvideo/models"
	"github.com/trofimovm/summvideo/repositories"
	"github.com/trofimovm/summvideo/utils"
)

// TelegramAuthHandler обрабатывает авторизацию через Telegram
func TelegramAuthHandler(c *gin.Context) {
	var authData models.TelegramAuthData

	// Получаем параметры из запроса
	authData.ID, _ = strconv.ParseInt(c.Query("id"), 10, 64)
	authData.FirstName = c.Query("first_name")
	authData.LastName = c.Query("last_name")
	authData.Username = c.Query("username")
	authData.PhotoURL = c.Query("photo_url")
	authData.AuthDate, _ = strconv.ParseInt(c.Query("auth_date"), 10, 64)
	authData.Hash = c.Query("hash")

	// Валидация данных
	if authData.ID == 0 || authData.Hash == "" || authData.AuthDate == 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid auth data",
		})
		return
	}

	// Проверяем подпись данных
	if err := utils.ValidateTelegramAuth(&authData); err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "Invalid auth signature: " + err.Error(),
		})
		return
	}

	// Создаем или обновляем пользователя
	userRepo := repositories.UserRepository{}
	user, err := userRepo.CreateOrUpdate(&authData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to create/update user: " + err.Error(),
		})
		return
	}

	// Генерируем JWT токен
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to generate token: " + err.Error(),
		})
		return
	}

	// Возвращаем токен и информацию о пользователе
	c.JSON(http.StatusOK, models.AuthResponse{
		Token: token,
		User:  *user,
	})
}

// AdminLoginHandler обрабатывает вход администратора
func AdminLoginHandler(c *gin.Context) {
	var credentials models.AdminCredentials

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Неверные данные: " + err.Error(),
		})
		return
	}

	// Проверяем учетные данные
	userRepo := repositories.UserRepository{}
	user, err := userRepo.CheckAdminCredentials(credentials.Username, credentials.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "Ошибка авторизации: " + err.Error(),
		})
		return
	}

	// Генерируем JWT токен
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Ошибка генерации токена: " + err.Error(),
		})
		return
	}

	// Возвращаем токен и информацию о пользователе
	c.JSON(http.StatusOK, models.AuthResponse{
		Token: token,
		User:  *user,
	})
}

// GetUserProfile возвращает профиль текущего пользователя
func GetUserProfile(c *gin.Context) {
	// Получаем ID пользователя из контекста (установлен в AuthMiddleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "Not authenticated",
		})
		return
	}

	// Получаем информацию о пользователе
	userRepo := repositories.UserRepository{}
	user, err := userRepo.FindByID(userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to get user profile: " + err.Error(),
		})
		return
	}

	// Получаем оставшееся время использования
	remainingSeconds, err := userRepo.GetRemainingUsageSeconds(userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to get remaining usage: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
		"remaining_usage_seconds": remainingSeconds,
	})
}

// GetUserHistory возвращает историю использования для текущего пользователя
func GetUserHistory(c *gin.Context) {
	// Получаем ID пользователя из контекста
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "Not authenticated",
		})
		return
	}

	// Получаем параметры пагинации
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	
	offset := (page - 1) * pageSize

	// Получаем историю использования
	usageRepo := repositories.UsageRepository{}
	history, err := usageRepo.FindByUserID(userID.(int64), pageSize, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to get usage history: " + err.Error(),
		})
		return
	}

	// Получаем общее количество записей для пагинации
	total, err := usageRepo.GetTotalCount(userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to get total count: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"history": history,
		"pagination": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
			"pages":     (total + pageSize - 1) / pageSize,
		},
	})
}