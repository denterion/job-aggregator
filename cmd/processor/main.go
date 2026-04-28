package main

import (
	"context"
	"encoding/json"
	"fmt"
	"job-aggregator/internal/models"
	"job-aggregator/internal/transport"
	"log"
	"strings"
)

func main() {
	fmt.Println("Сервис обработки вакансий запущен!")

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
		fmt.Printf("Вывод в консоль вакансии: %s от %s\n", v.Title, v.Company)
	}
}
