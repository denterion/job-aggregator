package processor

import (
	"job-aggregator/internal/models"
	"testing"
)

func TestShouldProcess(t *testing.T) {

	tests := []struct {
		name     string
		title    string
		expected bool
	}{
		{"Valid Go", "Middle Golang Developer", true},
		{"Valid lowercase", "junior go dev", true},
		{"Invalid Python", "Python Senior Architect", false},
		{"Invalid PHP", "Laravel backend", false},
		{"Valid with noise", "Fullstack {Go/React}", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := models.Vacancy{Title: tt.title}
			result := ShouldProcess(v)

			if result != tt.expected {
				t.Errorf("ShouldProcess() для %s, ожидалось %v, получили %v", tt.title, tt.expected, result)
			}
		})
	}
}
