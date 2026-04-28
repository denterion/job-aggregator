package processor

import (
	"job-aggregator/internal/models"
	"strings"
)

func ShouldProcess(v models.Vacancy) bool {
	title := strings.ToLower(v.Title)

	keywords := []string{"go", "golang", "gopher"}

	for _, kw := range keywords {
		if strings.Contains(title, kw) {
			return true
		}
	}
	return false
}
