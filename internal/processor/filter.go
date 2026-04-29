package processor

import (
	"job-aggregator/internal/models"
	"strings"
)

func ShouldProcess(v models.Vacancy) bool {
	title := strings.ToLower(v.Title)
	description := strings.ToLower(v.Description)

	isGo := containsAny(title, "go", "golang", "gopher") || containsAny(description, "golang", "язык go")

	if !isGo {
		return false
	}

	isJunior := containsAny(title, "junior", "intern", "стажер", "младший")
	if isJunior {
		return true
	}

	return isGo
}

func containsAny(text string, keywords ...string) bool {
	for _, kw := range keywords {
		if strings.Contains(text, kw) {
			return true
		}

	}
	return false
}
