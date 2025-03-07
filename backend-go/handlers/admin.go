package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/trofimovm/summvideo/models"
	"github.com/trofimovm/summvideo/repositories"
)

// AdminMiddleware проверяет, что пользователь является администратором
func AdminMiddleware(c *gin.Context) {
	// Получаем ID пользователя из контекста (установлен в AuthMiddleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "Требуется авторизация",
		})
		c.Abort()
		return
	}

	// Получаем информацию о пользователе
	userRepo := repositories.UserRepository{}
	user, err := userRepo.FindByID(userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Ошибка получения информации о пользователе: " + err.Error(),
		})
		c.Abort()
		return
	}

	// Проверяем, что пользователь - админ
	if !user.IsAdmin {
		c.JSON(http.StatusForbidden, models.ErrorResponse{
			Error: "Недостаточно прав для доступа",
		})
		c.Abort()
		return
	}

	c.Next()
}

// AdminLoginPage отображает страницу входа для администратора
func AdminLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_login.html", gin.H{})
}

// AdminDashboardPage отображает панель администратора
func AdminDashboardPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_dashboard.html", gin.H{})
}

// GetAllUsers возвращает список всех пользователей (для админа)
func GetAllUsers(c *gin.Context) {
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

	// Получаем список пользователей
	userRepo := repositories.UserRepository{}
	users, err := userRepo.GetAllUsers(pageSize, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Ошибка получения списка пользователей: " + err.Error(),
		})
		return
	}

	// Получаем общее количество пользователей для пагинации
	total, err := userRepo.GetTotalUsersCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Ошибка получения общего количества пользователей: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"pagination": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
			"pages":     (total + pageSize - 1) / pageSize,
		},
	})
}

// UpdateUserUsageLimit обновляет лимит использования для пользователя
func UpdateUserUsageLimit(c *gin.Context) {
	var update models.UsageLimitUpdate

	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Неверные данные: " + err.Error(),
		})
		return
	}

	// Проверяем, что пользователь существует
	userRepo := repositories.UserRepository{}
	user, err := userRepo.FindByID(update.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Пользователь не найден",
		})
		return
	}

	// Обновляем лимит использования
	err = userRepo.UpdateUsageLimit(update.UserID, update.LimitSeconds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Ошибка обновления лимита: " + err.Error(),
		})
		return
	}

	// Получаем обновленную информацию о пользователе
	updatedUser, err := userRepo.FindByID(update.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Ошибка получения информации о пользователе: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Лимит использования обновлен",
		"user":    updatedUser,
	})
}

// GetUserUsage возвращает информацию об использовании сервиса конкретным пользователем
func GetUserUsage(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Неверный ID пользователя",
		})
		return
	}

	// Проверяем, что пользователь существует
	userRepo := repositories.UserRepository{}
	user, err := userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Пользователь не найден",
		})
		return
	}

	// Получаем оставшееся время использования
	remainingSeconds, err := userRepo.GetRemainingUsageSeconds(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Ошибка получения информации об использовании: " + err.Error(),
		})
		return
	}

	// Получаем историю использования (последние 10 записей)
	usageRepo := repositories.UsageRepository{}
	history, err := usageRepo.FindByUserID(userID, 10, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Ошибка получения истории использования: " + err.Error(),
		})
		return
	}

	// Получаем общее количество записей истории
	totalCount, err := usageRepo.GetTotalCount(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Ошибка получения общего количества записей: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":                   user,
		"remaining_usage_seconds": remainingSeconds,
		"recent_activity":         history,
		"total_activity_count":    totalCount,
	})
}