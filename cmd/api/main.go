package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"example.com/pz16/internal/db"
	"example.com/pz16/internal/httpapi"
	"example.com/pz16/internal/repo"
	"example.com/pz16/internal/service"
)

func main() {
	// Получаем DSN из переменных окружения или используем дефолт
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/pz16?sslmode=disable"
	}

	// Подключение к БД
	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to connect to db:", err)
	}
	defer sqlDB.Close()

	// Проверка соединения
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("cannot ping db:", err)
	}

	// Миграции
	db.MustApplyMigrations(sqlDB)

	// Репозиторий и сервис
	noteRepo := repo.NoteRepo{DB: sqlDB}
	svc := service.Service{Notes: noteRepo}

	// Gin роутер
	r := gin.Default()
	router := httpapi.Router{Svc: &svc}
	router.Register(r)

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Starting server on port", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
