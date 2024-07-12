package postgres

import (
	"database/sql/driver"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MaKYaro/url-shortener/internal/domain"
	"github.com/MaKYaro/url-shortener/internal/storage"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

type AnyTime struct{}

func (s AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestSaveURL(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub connection: %s", err)
	}
	defer db.Close()
	storage := Storage{db}

	now := time.Now()
	alias := domain.Alias{
		Value:  "sqlmock",
		URL:    "https://github.com/DATA-DOG/go-sqlmock",
		Expire: now,
	}

	mock.ExpectExec("INSERT INTO urls").
		WithArgs("sqlmock", "https://github.com/DATA-DOG/go-sqlmock", AnyTime{}).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = storage.SaveURL(&alias)
	require.NoError(t, err)
}

func TestSaveURLUniqueError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub connection: %s", err)
	}
	defer db.Close()
	s := Storage{db}

	now := time.Now()
	alias := domain.Alias{
		Value:  "sqlmock",
		URL:    "https://github.com/DATA-DOG/go-sqlmock",
		Expire: now,
	}

	mock.ExpectExec("INSERT INTO urls").
		WithArgs("sqlmock", "https://github.com/DATA-DOG/go-sqlmock", AnyTime{}).
		WillReturnError(&pq.Error{Code: pq.ErrorCode("23505")})

	err = s.SaveURL(&alias)
	require.EqualError(
		t,
		fmt.Errorf("storage.postgres.SaveURL: %w", storage.ErrURLExists),
		err.Error(),
	)
}

func TestSaveURLFullDiskError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub connection: %s", err)
	}
	defer db.Close()
	s := Storage{db}

	now := time.Now()
	alias := domain.Alias{
		Value:  "sqlmock",
		URL:    "https://github.com/DATA-DOG/go-sqlmock",
		Expire: now,
	}

	mock.ExpectExec("INSERT INTO urls").
		WithArgs("sqlmock", "https://github.com/DATA-DOG/go-sqlmock", AnyTime{}).
		WillReturnError(&pq.Error{Code: pq.ErrorCode("53100")})

	err = s.SaveURL(&alias)
	require.EqualError(
		t,
		fmt.Errorf("storage.postgres.SaveURL: %w", &pq.Error{Code: pq.ErrorCode("53100")}),
		err.Error(),
	)
}

func TestGetURL(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub connection: %s", err)
	}
	defer db.Close()
	s := Storage{db}

	exp, _ := time.Parse("2016-06-22 19:10:25-07", "2016-06-22 19:10:25-07")
	row := sqlmock.NewRows([]string{"alias", "url", "expire"}).
		AddRow("sqlmock", "https://github.com/DATA-DOG/go-sqlmock", exp)
	mock.ExpectQuery("SELECT").WillReturnRows(row)

	aliasExpected := domain.Alias{
		Value:  "sqlmock",
		URL:    "https://github.com/DATA-DOG/go-sqlmock",
		Expire: exp,
	}
	aliasResult, err := s.GetURL("sqlmock")
	require.NoError(t, err)
	require.Equal(t, aliasExpected, *aliasResult)
}
