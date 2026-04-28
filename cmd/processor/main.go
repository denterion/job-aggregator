package main

import (
	"context"
	"encoding/json"
	"fmt"
	"job-aggregator/internal/models"
	"job-aggregator/internal/storage"
	"job-aggregator/internal/transport"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	fmt.Printf("DEBUG: Пытаюсь подключиться как %s с паролем %s\n", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"))
	if err != nil {
		log.Println("Предупреждение: .env Файл не найден, используем системные переменные")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName)

	fmt.Println("Сервис обработки вакансий запущен!")

	db, err := storage.NewPostgresStorage(connStr)
	if err != nil {
		log.Fatalf("Ошибка авторизации в БД: %v", err)
	}

	defer db.Conn.Close(context.Background())

	reader := transport.GetReader()
	defer reader.Close()

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Ошибка при чтении: %v", err)
			continue
		}

		var v models.Vacancy
		err = json.Unmarshal(m.Value, &v)
		if err != nil {
			log.Printf("Ошибка декодирования: %v", err)
			continue
		}

		titleLower := strings.ToLower(v.Title)
		if !strings.Contains(titleLower, "go") && !strings.Contains(titleLower, "golang") {
			fmt.Printf("Не нужные вакансии: %s", v.Title)
			continue
		}

		err = db.SaveVacancy(v)
		if err != nil {
			log.Printf("Ошибка сохранения в БД: %v", err)
		} else {
			fmt.Printf("Сохранено в базу: %v", v.Title)
		}
	}
}
