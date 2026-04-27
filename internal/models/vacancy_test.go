package models

import (
	"strings"
	"testing"
)

func TestVacancyValidation(t *testing.T) {
	v := Vacancy{
		Title: "Golang Developer",
	}

	if v.Title == "" {
		t.Error("Заголовок не может быть пустым")
	}

	expectedPrefix := "Golang"
	if !contains(v.Title, expectedPrefix) {
		t.Errorf("Ожидалось, что заголовок содержит %s, но получили %s", expectedPrefix, v.Title)
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
