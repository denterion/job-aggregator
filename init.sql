CREATE TABLE IF NOT EXISTS vacancies(
    id SERIAL PRIMARY KEY,
    external_id TEXT UNIQUE,
    title TEXT NOT NULL,
    company TEXT,
    location TEXT,
    salary TEXT,
    url TEXT UNIQUE,
    source TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
);

CREATE INDEX IF NOT EXISTS idx_vacancies_title ON vacancies(title);