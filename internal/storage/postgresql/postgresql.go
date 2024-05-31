package postgresql

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"stellar_journal/internal/models/nasa_api_models"
	"stellar_journal/internal/models/stellar_journal_models"
	"stellar_journal/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(dbUri string) (*Storage, error) {
	const op = "internal/storage.postgresql.NewStorage"

	db, err := sql.Open("postgres", dbUri)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to connect to the database: %w", op, err)
	}

	stmt1, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS nasa_apod (
			id SERIAL PRIMARY KEY,
			copyright TEXT,
			apod_date DATE UNIQUE,
			explanation TEXT,
			hdurl TEXT,
			media_type TEXT,
			service_version TEXT,
			title TEXT,
			url TEXT
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create db: %w", op, err)
	}

	_, err = stmt1.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create db: %w", op, err)
	}

	stmt2, err := db.Prepare(`
    	CREATE INDEX IF NOT EXISTS nasa_apod_date_idx ON nasa_apod (date)
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create index: %w", op, err)
	}

	_, err = stmt2.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create index: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveAPOD(apod *nasa_api_models.APODResp) (int64, error) {
	const op = "internal/storage.postgresql.SaveAPOD"

	stmt, err := s.db.Prepare(`
		INSERT INTO nasa_apod (copyright, apod_date, explanation, hdurl, media_type, service_version, title, url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}

	res, err := stmt.Exec(apod.Copyright, apod.Date, apod.Explanation, apod.Hdurl, apod.MediaType, apod.ServiceVersion, apod.Title, apod.Url)
	if err != nil {
		if postgresErr, ok := err.(*pq.Error); ok && postgresErr.Code == "23505" {
			return 0, fmt.Errorf("%s: failed to insert data: %w", op, storage.ErrAPODExists)
		}
		return 0, fmt.Errorf("%s: failed to insert data: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetAPOD(date string) (*stellar_journal_models.APOD, error) {
	const op = "internal/storage.postgresql.GetAPOD"

	stmt, err := s.db.Prepare(`
		SELECT id, copyright, apod_date, explanation, hdurl, media_type, service_version, title, url
		FROM nasa_apod
		WHERE apod_date = $1
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}

	var apod stellar_journal_models.APOD
	err = stmt.QueryRow(date).Scan(&apod.Id, &apod.Copyright, &apod.Date, &apod.Explanation, &apod.Hdurl, &apod.MediaType, &apod.ServiceVersion, &apod.Title, &apod.Url)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: failed to get data: %w", op, storage.ErrAPODNotFound)
		}
		return nil, fmt.Errorf("%s: failed to get data: %w", op, err)
	}

	return &apod, nil
}

func (s *Storage) GetJournal() (*[]stellar_journal_models.APOD, error) {
	const op = "internal/storage.postgresql.GetJournal"

	stmt, err := s.db.Prepare(`
		SELECT id, copyright, apod_date, explanation, hdurl, media_type, service_version, title, url
		FROM nasa_apod
		ORDER BY apod_date DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get data: %w", op, err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Printf("%s: failed to close rows: %v\n", op, err)
		}
	}(rows)

	var apods []stellar_journal_models.APOD
	for rows.Next() {
		var apod stellar_journal_models.APOD
		err = rows.Scan(&apod.Id, &apod.Copyright, &apod.Date, &apod.Explanation, &apod.Hdurl, &apod.MediaType, &apod.ServiceVersion, &apod.Title, &apod.Url)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to scan data: %w", op, err)
		}
		apods = append(apods, apod)
	}

	return &apods, nil
}

func (s *Storage) Close() error {
	const op = "internal/storage.postgresql.Close"

	err := s.db.Close()
	if err != nil {
		return fmt.Errorf("%s: failed to close db: %w", op, err)
	}

	return nil
}
