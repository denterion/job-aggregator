package main

import (
	"context"
	"encoding/json"
	"fmt"
	"job-aggregator/internal/models"
	"job-aggregator/internal/processor"
	"job-aggregator/internal/storage"
	"job-aggregator/internal/transport"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	slog.Info("Сервис запущен", "service", "processor", "pid", os.Getpid())

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
		slog.Error("Ошибка авторизации в БД:", "details", err)
	}

	defer db.Conn.Close(context.Background())

	reader := transport.GetReader()
	defer reader.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	slog.Info("Сервис обработки вакансий запущен", "status", "waiting_for_messages")

	for {
		select {
		case <-ctx.Done():
			slog.Info("Завершение работы по сигналу")
			return
		default:
			m, err := reader.ReadMessage(ctx)
			if ctx.Err() != nil {
				slog.Error("Ошибка при чтении:", "details", err)
				continue
			}

			var v models.Vacancy
			err = json.Unmarshal(m.Value, &v)
			if err != nil {
				slog.Error("Ошибка декодирования:", "details", err)
				continue
			}

			if !processor.ShouldProcess(v) {
				slog.Debug("Вакансия пропущена фильтром", "title", v.Title)
				continue
			}

			err = db.SaveVacancy(v)
			if err != nil {
				slog.Error("Ошибка сохранения в БД:", "details", err)
			} else {
				fmt.Printf("Сохранено в базу: %v", v.Title)
			}

		}
	}
}
