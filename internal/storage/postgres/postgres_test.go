package postgres

import (
	"database/sql"
	"database/sql/driver"
	"errors"
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
	s := Storage{db}

	now := time.Now()
	alias := domain.Alias{
		Value:  "sqlmock",
		URL:    "https://github.com/DATA-DOG/go-sqlmock",
		Expire: now,
	}

	mock.ExpectExec("INSERT INTO urls").
		WithArgs("sqlmock", "https://github.com/DATA-DOG/go-sqlmock", AnyTime{}).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = s.SaveURL(&alias)
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
	require.ErrorIs(t, err, storage.ErrAliasExists)
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

	resultError := errors.New("full disk error")

	mock.ExpectExec("INSERT INTO urls").
		WithArgs("sqlmock", "https://github.com/DATA-DOG/go-sqlmock", AnyTime{}).
		WillReturnError(resultError)

	err = s.SaveURL(&alias)
	require.ErrorIs(t, err, resultError)
}

func TestGetURL(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub connection: %s", err)
	}
	defer db.Close()
	s := Storage{db}

	row := sqlmock.NewRows([]string{"url"}).
		AddRow("https://github.com/DATA-DOG/go-sqlmock")
	mock.ExpectQuery("SELECT url FROM urls WHERE alias =").
		WillReturnRows(row)

	urlExpected := "https://github.com/DATA-DOG/go-sqlmock"
	urlResult, err := s.GetURL("sqlmock")
	require.NoError(t, err)
	require.Equal(t, urlExpected, urlResult)
}

func TestGetURLErrorURLNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub connection: %s", err)
	}
	defer db.Close()
	s := Storage{db}

	mock.ExpectQuery("SELECT url FROM urls WHERE alias =").
		WillReturnError(sql.ErrNoRows)
	urlResult, err := s.GetURL("ui")
	require.Empty(t, urlResult)
	require.ErrorIs(t, err, storage.ErrURLNotFound)
}

func TestGetURLErrorInAssignment(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub connection: %s", err)
	}
	defer db.Close()
	s := Storage{db}

	resultError := errors.New("error in assignment")

	mock.ExpectQuery("SELECT url FROM urls WHERE alias =").
		WillReturnError(resultError)
	urlResult, err := s.GetURL("ui")
	require.Empty(t, urlResult)
	require.ErrorIs(t, err, resultError)
}

func TestDeleteURL(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub connection: %s", err)
	}
	defer db.Close()
	s := Storage{db}

	mock.ExpectExec("DELETE FROM urls WHERE alias = ").
		WithArgs("adsfaf").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = s.DeleteURL("adsfaf")
	require.NoError(t, err)
}

func TestDeleteURLErrorWarning(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open stub connection: %s", err)
	}
	defer db.Close()
	s := Storage{db}

	resultError := errors.New("error warning")

	mock.ExpectExec("DELETE FROM urls WHERE alias = ").
		WillReturnError(resultError)
	err = s.DeleteURL("ui")
	require.ErrorIs(t, err, resultError)
}
