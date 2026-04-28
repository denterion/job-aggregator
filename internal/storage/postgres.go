package storage

import (
	"context"
	"job-aggregator/internal/models"

	"github.com/jackc/pgx/v5"
)

type PostgresStorage struct {
	Conn *pgx.Conn
}

func NewPostgresStorage(connStr string) (*PostgresStorage, error) {
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}
	return &PostgresStorage{Conn: conn}, nil
}

func (s *PostgresStorage) SaveVacancy(v models.Vacancy) error {

	query := `
		INSERT INTO vacancies (external_id, title, company, location, salary, url, source)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (url) DO NOTHING
	`

	_, err := s.Conn.Exec(context.Background(), query, v.ID, v.Title, v.Company, v.Location, v.Salary, v.URL, v.Source)
	return err
}
