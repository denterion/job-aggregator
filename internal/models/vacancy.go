package models

import "time"

type Vacancy struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Company     string    `json:"company"`
	Location    string    `json:"Location"`
	Salary      string    `json:"salary"`
	URL         string    `json:"url"`
	Source      string    `json:"source"`
	CreatedAt   time.Time `json:"created_at"`
}
