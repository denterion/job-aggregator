package main

import (
	"fmt"
	"job-aggregator/internal/models"
	"job-aggregator/internal/scraper"
	"job-aggregator/internal/transport"
	"log"
	"sync"
)

func main() {
	fmt.Println("Запускаем парсер!")

	habr := scraper.NewHabrScraper()

	vacancies, err := habr.Parse("golang")
	if err != nil {
		log.Fatalf("Ошибка при парсинге: %v", err)
	}

	fmt.Printf("Найдено вакансий: %d\n", len(vacancies))

	var wg sync.WaitGroup

	for _, v := range vacancies {
		wg.Add(1)

		go func(v models.Vacancy) {
			defer wg.Done()

			err := transport.SendVacancy(v)
			if err != nil {
				fmt.Printf("Ошибка отправки вакансии [%s]: %v\n", v.Title, err)
				return
			}
			fmt.Printf("Вакансия отправлена: %s (%s)\n", v.Title, v.Company)
		}(v)
	}
	wg.Wait()
	fmt.Println("Парсер завершил работу!")
}
