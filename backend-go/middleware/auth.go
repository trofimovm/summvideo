package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/trofimovm/summvideo/models"
	"github.com/trofimovm/summvideo/repositories"
	"github.com/trofimovm/summvideo/utils"
)

// AuthMiddleware защищает маршруты, требующие аутентификации
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Проверяем режим разработки
		if os.Getenv("DEV_MODE") == "true" {
			log.Println("Используется DEV_MODE - аутентификация пропущена")
			// В режиме разработки используем тестового пользователя
			c.Set("userID", int64(1))
			c.Next()
			return
		}

		// Получаем токен из заголовка Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Authorization header is required",
			})
			return
		}

		// Проверяем формат (Bearer token)
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Authorization header format must be Bearer {token}",
			})
			return
		}

		tokenString := parts[1]

		// Валидируем JWT токен
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Invalid or expired token",
			})
			return
		}

		// Проверяем существование пользователя
		userRepo := repositories.UserRepository{}
		user, err := userRepo.FindByID(claims.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "User not found",
			})
			return
		}

		// Проверяем активность пользователя
		if !user.IsActive {
			c.AbortWithStatusJSON(http.StatusForbidden, models.ErrorResponse{
				Error: "User is deactivated",
			})
			return
		}

		// Сохраняем ID пользователя в контексте для дальнейшего использования
		c.Set("userID", user.ID)

		c.Next()
	}
}