package processor

import (
	"job-aggregator/internal/models"
	"testing"
)

func TestShouldProcess(t *testing.T) {

	tests := []struct {
		name     string
		title    string
		desc     string
		expected bool
	}{
		{"Junior Go", "Junior Developer", "We use Golang for our mircoservices", true},
		{"Intern Russian", "Стажер-разработчик", "Будем писать на языке Go", true},
		{"Not Go", "Junior Java Developer", "Spring boot is cool", false},
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
