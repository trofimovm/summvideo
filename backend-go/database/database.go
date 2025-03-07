package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
)

// DB глобальный пул соединений с базой данных
var DB *pgxpool.Pool

// Initialize инициализирует соединение с базой данных
func Initialize() error {
	// Получаем URL подключения
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return errors.New("DATABASE_URL не указан")
	}

	// Создаем конфигурацию для пула соединений
	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return fmt.Errorf("ошибка при разборе DATABASE_URL: %v", err)
	}

	// Создаем пул соединений
	DB, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		return fmt.Errorf("ошибка при подключении к базе данных: %v", err)
	}

	// Проверяем соединение
	err = DB.Ping(context.Background())
	if err != nil {
		return fmt.Errorf("ошибка при проверке соединения с базой данных: %v", err)
	}

	log.Println("Подключение к базе данных установлено")
	return nil
}

// RunMigrations запускает миграции базы данных
func RunMigrations(migrationsDir string) error {
	// Получаем URL подключения
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return errors.New("DATABASE_URL не указан")
	}

	// Формируем путь к миграциям
	migrationPath := fmt.Sprintf("file://%s", migrationsDir)
	
	// Создаем экземпляр мигратора
	m, err := migrate.New(migrationPath, dbURL)
	if err != nil {
		return fmt.Errorf("ошибка при инициализации миграций: %v", err)
	}
	
	// Запускаем миграции
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("ошибка при выполнении миграций: %v", err)
	}
	
	log.Println("Миграции успешно применены")
	
	return nil
}

// Close закрывает соединение с базой данных
func Close() {
	if DB != nil {
		DB.Close()
		log.Println("Соединение с базой данных закрыто")
	}
}