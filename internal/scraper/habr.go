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
		),
	}
}

func (s *HabrScrapper) Parse(query string) ([]models.Vacancy, error) {
	var vacancies []models.Vacancy

	s.Collector.OnHTML(".vacancy-card", func(e *colly.HTMLElement) {
		var skills []string
		e.ForEach(".vacancy-card__skills-item", func(_ int, el *colly.HTMLElement) {
			skills = append(skills, el.Text)
		})
		fullDesc := strings.Join(skills, " • ")
		v := models.Vacancy{
			ID:          e.Attr("id"),
			Title:       e.ChildText(".vacancy-card__title"),
			Description: fullDesc,
			Company:     e.ChildText(".vacancy-card__company-title a"),
			Location:    e.ChildText(".vacancy-card__meta"),
			URL:         "https://career.habr.com" + e.ChildAttr(".vacancy-card__title-link", "href"),
			Source:      "Habr",
			CreatedAt:   time.Now(),
		}
		vacancies = append(vacancies, v)

		if v.Title == "" || v.URL == "https://career.habr.com" {
			fmt.Println("Скрапер нашел пустую карточку, проверь селекторы!")
			return
		}

		fmt.Printf("Нашел: %s [%s]\n", v.Title, v.Company) // Увидишь это в консоли скрапера
		vacancies = append(vacancies, v)

	})

	searchURL := fmt.Sprintf("https://career.habr.com/vacancies?q=%s&type=all", query)
	err := s.Collector.Visit(searchURL)

	return vacancies, err
}
