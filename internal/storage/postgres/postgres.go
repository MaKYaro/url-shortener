package postgres

import (
	"database/sql"
	"fmt"

	"github.com/MaKYaro/url-shortener/internal/config"
	"github.com/MaKYaro/url-shortener/internal/domain"
	"github.com/MaKYaro/url-shortener/internal/storage"
	"github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(cfg config.DBConnConfig) (*Storage, error) {
	const op = "storage.postgres.New"

	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(alias *domain.Alias) error {
	const op = "storage.postgres.SaveURL"

	query := "INSERT INTO urls (alias, url, expire) VALUES ($1, $2, $3)"
	_, err := s.db.Exec(query, alias.Value, alias.URL, alias.Expire)

	fmt.Println(err)

	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
		return fmt.Errorf("%s: %w", op, storage.ErrURLExists)
	}
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetURL(alias string) (*domain.Alias, error) {
	const op = "storage.postgres.GetURL"

	query := "SELECT alias, url, expire FROM urls WHERE alias = $1"
	row := s.db.QueryRow(query, alias)

	var a domain.Alias
	err := row.Scan(&a.Value, &a.URL, &a.Expire)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &a, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.postgres.DeleteURL"

	query := "DELETE FROM urls WHERE alias = $1"
	_, err := s.db.Exec(query, alias)

	if err != nil {
		return fmt.Errorf("%s: can't delete url: %w", op, err)
	}

	return nil
}
