package scraper

import (
	"fmt"
	"job-aggregator/internal/models"
	"strings"
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
			colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
		),
	}
}

func (s *HabrScrapper) Parse(query string) ([]models.Vacancy, error) {
	var vacancies []models.Vacancy
	s.Collector.OnHTML(".vacancy-card__title-link", func(e *colly.HTMLElement) {
		vacancyURL := e.Request.AbsoluteURL(e.Attr("href"))
		e.Request.Visit(vacancyURL)
	})

	s.Collector.OnHTML(".page-container", func(e *colly.HTMLElement) {
		title := strings.TrimSpace(e.ChildText(".page-title__title"))
		if title == "" {
			return
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

		fmt.Printf("Нашел вакансию: %s\n", v.Title)
		vacancies = append(vacancies, v)
	})

	searchURL := fmt.Sprintf("https://career.habr.com/vacancies?q=%s&type=all", query)
	err := s.Collector.Visit(searchURL)

	return vacancies, err
}
