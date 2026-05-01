package storage

import (
	"context"
	"fmt"
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
		INSERT INTO vacancies (external_id, title, description, company, location, salary, url, source)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (url) DO NOTHING
	`

	_, err := s.Conn.Exec(context.Background(), query, v.ID, v.Title, v.Description, v.Company, v.Location, v.Salary, v.URL, v.Source)
	return err
}

func (s *PostgresStorage) GetLatestVacancies(limit int) ([]models.Vacancy, error) {
	query := `
		SELECT title, company, url, description, salary
		FROM vacancies
		ORDER BY created_at DESC
		LIMIT $1
	`
	rows, err := s.Conn.Query(context.Background(), query, limit)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к бд: %v", err)
	}

	defer rows.Close()

	var vacancies []models.Vacancy
	for rows.Next() {
		var v models.Vacancy
		err := rows.Scan(&v.Title, &v.Company, &v.URL, &v.Description, &v.Salary)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %v", err)
		}
		vacancies = append(vacancies, v)
	}
	return vacancies, nil
}
