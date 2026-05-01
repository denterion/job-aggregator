package scraper

import (
	"fmt"
	"job-aggregator/internal/models"
	"time"

	"github.com/gocolly/colly/v2"
)

type HabrScrapper struct {
	Collector *colly.Collector
}

func NewHabrScraper() *HabrScrapper {
	return &HabrScrapper{
		Collector: colly.NewCollector(
			colly.AllowedDomains("career.habr.com"),
		),
	}
}

func (s *HabrScrapper) Parse(query string) ([]models.Vacancy, error) {
	var vacancies []models.Vacancy

	// 1. На странице списка находим все ссылки на вакансии и заходим в них
	s.Collector.OnHTML(".vacancy-card__title-link", func(e *colly.HTMLElement) {
		vacancyURL := e.Request.AbsoluteURL(e.Attr("href"))
		// Заставляем коллектор зайти внутрь каждой найденной вакансии
		e.Request.Visit(vacancyURL)
	})

	// 2. А здесь описываем, ЧТО собирать ВНУТРИ страницы вакансии
	s.Collector.OnHTML(".page-container", func(e *colly.HTMLElement) {
		// Проверяем, что мы действительно на странице вакансии (есть заголовок)
		title := e.ChildText(".vacancy-header__title")
		if title == "" {
			return // Это не страница вакансии, пропускаем
		}

		v := models.Vacancy{
			Title:       title,
			Company:     e.ChildText(".company_name"),
			Salary:      e.ChildText(".vacancy-header__salary"),
			Description: e.ChildText(".vacancy-description__text"), // Текст описания
			URL:         e.Request.URL.String(),
			Source:      "Habr",
			CreatedAt:   time.Now(),
		}

		fmt.Printf("✅ Нашел вакансию: %s\n", v.Title)
		vacancies = append(vacancies, v)
	})

	searchURL := fmt.Sprintf("https://career.habr.com/vacancies?q=%s&type=all", query)
	err := s.Collector.Visit(searchURL)

	return vacancies, err
}
