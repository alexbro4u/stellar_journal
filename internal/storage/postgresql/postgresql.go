package postgresql

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"stellar_journal/internal/models/nasa_api_models"
	"stellar_journal/internal/models/stellar_journal_models"
	"stellar_journal/internal/storage"
)

type PostgresDriver struct{}

func (d *PostgresDriver) Open(db *sql.DB) (database.Driver, error) {
	return postgres.WithInstance(db, &postgres.Config{})
}

type Storage struct {
	DB     *sql.DB
	Driver *PostgresDriver
}

func NewStorage(user, pass, name, host string) (*Storage, error) {
	const op = "internal/storage/postgresql.NewStorage"

	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable", user, pass, name, host))
	if err != nil {
		return nil, fmt.Errorf("%s: failed to connect to the database: %w", op, err)
	}

	return &Storage{DB: db}, nil
}

func (s *Storage) SaveAPOD(apod *nasa_api_models.APODResp) error {
	const op = "internal/storage/postgresql.SaveAPOD"

	stmt, err := s.DB.Prepare(`
		INSERT INTO nasa_apod (copyright, apod_date, explanation, hdurl, media_type, service_version, title, url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`)
	if err != nil {
		return fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}

	_, err = stmt.Exec(apod.Copyright, apod.Date, apod.Explanation, apod.Hdurl, apod.MediaType, apod.ServiceVersion, apod.Title, apod.Url)
	if err != nil {
		if postgresErr, ok := err.(*pq.Error); ok && postgresErr.Code == "23505" {
			return fmt.Errorf("%s: failed to insert data: %w", op, storage.ErrAPODExists)
		}
		return fmt.Errorf("%s: failed to insert data: %w", op, err)
	}

	return nil
}

func (s *Storage) GetAPOD(date string) (*stellar_journal_models.APOD, error) {
	const op = "internal/storage/postgresql.GetAPOD"

	stmt, err := s.DB.Prepare(`
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
	const op = "internal/storage/postgresql.GetJournal"

	stmt, err := s.DB.Prepare(`
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
	const op = "internal/storage/postgresql.Close"

	err := s.DB.Close()
	if err != nil {
		return fmt.Errorf("%s: failed to close db: %w", op, err)
	}

	return nil
}
