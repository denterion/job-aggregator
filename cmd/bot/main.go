package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"job-aggregator/internal/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	db, err := storage.NewPostgresStorage(connStr)
	if err != nil {
		slog.Error("Не удалось подключиться к бд", "err", err)
		os.Exit(1)
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		slog.Error("Ошибка инициализации бота", "err", err)
		os.Exit(1)
	}

	slog.Info("Бот авторизован", "user", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("Бот останавливается...")
			return
		case update := <-updates:
			if update.Message == nil {
				continue
			}

			if update.Message.IsCommand() {
				handleCommand(bot, update.Message, db)
			}
		}
	}

}

func handleCommand(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, db *storage.PostgresStorage) {
	switch msg.Command() {
	case "start":
		reply := "Привет! Я твой личный агрегатор Go-вакансий. Используй /latest, чтобы увидеть свежие предложения."
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, reply))

	case "latest":
		vacs, err := db.GetLatestVacancies(10)
		if err != nil {
			slog.Error("Ошибка получения вакансий", "err", err)
			bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Упс, что-то пошло не так, при обращении к базе."))
			return
		}

		if len(vacs) == 0 {
			bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Пока вакансий нет. Запустите скапер"))
			return
		}

		for _, v := range vacs {
			text := fmt.Sprintf(
				"<b>%s</b>\n Компания: %s\n Зарплата: %s\n Стек: <i>%s</i>\n\n <a href=\"%s\">Открыть вакансию</a>",
				v.Title, v.Company, v.Salary, v.Description, v.URL,
			)

			newMsg := tgbotapi.NewMessage(msg.Chat.ID, text)
			newMsg.ParseMode = "HTML"
			bot.Send(newMsg)
		}

	}

}
