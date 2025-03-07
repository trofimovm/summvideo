package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/trofimovm/summvideo/database"
	"github.com/trofimovm/summvideo/handlers"
	"github.com/trofimovm/summvideo/middleware"
	"github.com/trofimovm/summvideo/utils"
)

func main() {
	// Загрузка переменных окружения из .env файла
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, используем переменные окружения системы")
	}

	// Инициализация директории для логов
	logDir := utils.GetEnv("LOG_DIR", "/var/log/summvideo")
	os.MkdirAll(logDir, os.ModePerm)

	// Установка режима Gin в зависимости от окружения
	if os.Getenv("GIN_MODE") != "release" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Создаем роутер Gin
	router := gin.Default()

	// Настройка CORS
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	// Максимальный размер файла (500 MB по умолчанию)
	router.MaxMultipartMemory = 500 << 20 // 500 MB

	// Определяем пути к статическим файлам и шаблонам
	staticDir := utils.GetEnv("STATIC_DIR", "./static")
	templatesDir := utils.GetEnv("TEMPLATES_DIR", "./templates")
	migrationsDir := utils.GetEnv("MIGRATIONS_DIR", "./migrations")

	// Проверяем и логируем пути
	log.Printf("Статические файлы: %s", staticDir)
	log.Printf("Шаблоны: %s", templatesDir)

	// Инициализация базы данных
	if err := database.Initialize(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Применение миграций
	if err := database.RunMigrations(migrationsDir); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Статические файлы
	router.Static("/static", staticDir)

	// Настройка HTML templates
	templatesPattern := templatesDir + "/*"
	router.LoadHTMLGlob(templatesPattern)

	// Открытые маршруты (не требуют авторизации)
	router.GET("/", handlers.HomePage)
	router.GET("/index.html", handlers.RedirectToHome)
	router.GET("/auth/telegram", handlers.TelegramAuthHandler)
	router.POST("/auth/admin", handlers.AdminLoginHandler)
	
	// Страницы администратора
	router.GET("/admin", handlers.AdminLoginPage)
	router.GET("/admin/dashboard", middleware.AuthMiddleware(), handlers.AdminMiddleware, handlers.AdminDashboardPage)
	
	// Защищенные маршруты (требуют авторизации)
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/upload_video/", handlers.UploadVideo)
		protected.GET("/profile", handlers.GetUserProfile)
		protected.GET("/history", handlers.GetUserHistory)
	}

	// Маршруты API для администратора
	adminAPI := router.Group("/api/admin")
	adminAPI.Use(middleware.AuthMiddleware(), handlers.AdminMiddleware)
	{
		adminAPI.GET("/users", handlers.GetAllUsers)
		adminAPI.GET("/users/:id/usage", handlers.GetUserUsage)
		adminAPI.PUT("/users/limit", handlers.UpdateUserUsageLimit)
	}

	// Проверка OPENAI_API_KEY
	if os.Getenv("OPENAI_API_KEY") == "" {
		log.Println("ВНИМАНИЕ: OPENAI_API_KEY не найден. Установите переменную окружения для работы с OpenAI API.")
	}

	// Запуск сервера
	port := utils.GetEnv("PORT", "8000")
	log.Printf("Сервер запущен на порту %s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}